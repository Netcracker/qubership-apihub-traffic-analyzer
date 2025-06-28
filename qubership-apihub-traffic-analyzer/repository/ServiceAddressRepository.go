// Copyright 2024-2025 NetCracker Technology Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package repository

import (
	"errors"
	"fmt"

	"github.com/Netcracker/qubership-apihub-traffic-analyzer/qubership-apihub-traffic-analyzer/db"
	"github.com/Netcracker/qubership-apihub-traffic-analyzer/qubership-apihub-traffic-analyzer/entities"
	"github.com/Netcracker/qubership-apihub-traffic-analyzer/qubership-apihub-traffic-analyzer/view"
	"github.com/go-pg/pg/v10"
	_ "github.com/shaj13/libcache/lru"
	log "github.com/sirupsen/logrus"
)

type ServiceAddressRepository interface {
	GetServiceAddress(address, name, version, captureId string) (entities.ServiceAddress, error)
	GetServiceAddressByIp(address string) (*entities.ServiceAddress, error)
	Close()
}
type serviceAddressRepository struct {
	//instance libcache.Cache
	db      db.ConnectionProvider
	addrMap map[string]entities.ServiceAddress
}

func NewPeersCache(db db.ConnectionProvider) ServiceAddressRepository {
	//nc := serviceAddressRepository{instance: libcache.LRU.New(MinCacheSize), db: db}
	//nc.instance.SetTTL(CachedRecAge)
	//return &nc
	return &serviceAddressRepository{db: db, addrMap: make(map[string]entities.ServiceAddress)}
}

func (sar *serviceAddressRepository) GetServiceAddressByIp(address string) (*entities.ServiceAddress, error) {
	//serviceAddr, exists := sar.instance.Load(address)
	serviceAddr, exists := sar.addrMap[address]
	if exists {
		//if val, converted := serviceAddr.(entities.ServiceAddress); converted {
		//	return &val, nil
		//} else {
		//	return nil, errors.New("unable to cast cache value for address: " + address)
		//}
		return &serviceAddr, nil
	}
	return nil, fmt.Errorf("service address not found for address %s", address)
}

func (sar *serviceAddressRepository) addCacheValue(serviceAddr entities.ServiceAddress) {
	//if _, exists := sar.instance.Load(serviceAddr.Address); exists {
	//	sar.instance.Delete(serviceAddr.Address)
	//}
	//sar.instance.Store(serviceAddr.Address, serviceAddr)
	sar.addrMap[serviceAddr.Address] = serviceAddr
}

func (sar *serviceAddressRepository) Close() {
	//sar.instance.Purge()
	sar.addrMap = make(map[string]entities.ServiceAddress)
}

func (sar *serviceAddressRepository) GetServiceAddress(address, name, version, captureId string) (entities.ServiceAddress, error) {
	result := entities.ServiceAddress{
		Address:   address,
		Name:      name,
		Version:   version,
		CaptureId: captureId,
	}
	var (
		err error
	)
	if name == view.EmptyString {
		err = sar.db.GetConnection().Model(&result).Where("ip_address=? and capture_id=?", address, captureId).First()
	} else {
		err = sar.db.GetConnection().Model(&result).Where("ip_address=? and service_name=? and capture_id=?", address, name, captureId).First()
	}
	if err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			result = entities.ServiceAddress{
				Address:   address,
				Name:      name,
				Version:   version,
				CaptureId: captureId,
			}
			log.Debugf("insertServiceAddress - Address=%s, Name=%s, Version=%s, captureId: %s", address, name, version, captureId)
			err = sar.insertServiceAddress(&result)
		} else {
			log.Debugf("select service address Id: %s", err.Error())
		}
		if err == nil {
			sar.addCacheValue(result)
		}
	} else {
		changed := false
		if name != view.EmptyString && result.Name != name {
			result.Name = name
			changed = true
		}
		if version != view.EmptyString && result.Version != version {
			result.Version = version
			changed = true
		}
		if changed {
			_, err := sar.db.GetConnection().Model(&result).Update()
			if err == nil {
				sar.addCacheValue(result)
			}
		}
	}
	return result, err
}

func (sar *serviceAddressRepository) insertServiceAddress(svcAddress *entities.ServiceAddress) error {
	result := new(entities.ServiceAddress)
	_, err := sar.db.GetConnection().Model(svcAddress).Returning("address_id").Insert(result)
	if err == nil {
		svcAddress.Id = result.Id
	}
	return err

}
