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
	"fmt"
	"time"

	"github.com/Netcracker/qubership-apihub-traffic-analyzer/qubership-apihub-traffic-analyzer/db"
	"github.com/Netcracker/qubership-apihub-traffic-analyzer/qubership-apihub-traffic-analyzer/entities"
	"github.com/Netcracker/qubership-apihub-traffic-analyzer/qubership-apihub-traffic-analyzer/view"
	"github.com/shaj13/libcache"
	_ "github.com/shaj13/libcache/lru"
)

const (
	// MinCacheSize initial cache size in records
	MinCacheSize = 300
	// CachedRecAge duration when record will be stored in cache
	CachedRecAge = 15 * time.Minute
)

type HttpHeadersCache interface {
	GetPrimaryKeyValue(key, value string) (string, error)
	Close()
}
type httpHeadersCache struct {
	instance libcache.Cache
	db       db.ConnectionProvider
}

func NewHttpHeadersCache(db db.ConnectionProvider) HttpHeadersCache {
	nc := httpHeadersCache{instance: libcache.LRU.New(MinCacheSize), db: db}
	nc.instance.SetTTL(CachedRecAge)
	// register expiration callback
	return &nc
}

func (hc *httpHeadersCache) GetPrimaryKeyValue(key, value string) (string, error) {
	h := entities.NewHttpHeader(key, value)
	_, exists := hc.instance.Load(h.Id)
	if exists {
		return h.Id, nil
	}
	err := hc.headerRecordExists(h)
	if err == nil {
		hc.instance.Store(h.Id, h)
		return h.Id, nil
	} else {
		err = hc.insertNewHeaderRecord(h)
		if err == nil {
			hc.instance.Store(h.Id, h)
			return h.Id, nil
		}
	}
	return view.EmptyString, fmt.Errorf("unable to get header id for: %s : %v", key, err)
}

func (hc *httpHeadersCache) Close() {
	hc.instance.Purge()
}

func (hc *httpHeadersCache) headerRecordExists(httpHeader entities.HttpHeaderItem) error {
	result := new(entities.HttpHeaderItem)
	return hc.db.GetConnection().Model(result).Where(
		"header_id=?", httpHeader.Id).Select()
}

func (hc *httpHeadersCache) insertNewHeaderRecord(httpHeader entities.HttpHeaderItem) error {
	result := new(entities.HttpHeaderItem)
	_, err := hc.db.GetConnection().Model(&httpHeader).Insert(result)
	return err
}
