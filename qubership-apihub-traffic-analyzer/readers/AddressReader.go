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

package readers

import (
	"bufio"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/Netcracker/qubership-apihub-traffic-analyzer/qubership-apihub-traffic-analyzer/db"
	"github.com/Netcracker/qubership-apihub-traffic-analyzer/qubership-apihub-traffic-analyzer/entities"
	"github.com/Netcracker/qubership-apihub-traffic-analyzer/qubership-apihub-traffic-analyzer/repository"
	"github.com/Netcracker/qubership-apihub-traffic-analyzer/qubership-apihub-traffic-analyzer/view"
	log "github.com/sirupsen/logrus"
)

type HostsReader interface {
	Read(hostsFile string) error
	GetServiceByIp(ip string) (*entities.ServiceAddress, error)
	Close() error
}

type hostsReaderImpl struct {
	workDir    string
	captureId  string
	db         db.ConnectionProvider
	hostsCache repository.ServiceAddressRepository
}

func NewHostsReader(workDir, captureId string, db db.ConnectionProvider) (HostsReader, error) {
	hc := repository.NewPeersCache(db)
	if hc == nil {
		return nil, errors.New("unable to create cache for hosts")
	}
	return &hostsReaderImpl{
		workDir:    workDir,
		captureId:  captureId,
		db:         db,
		hostsCache: hc,
	}, nil
}

func (hr *hostsReaderImpl) Read(hostsFile string) error {
	fh, err := os.Open(hostsFile)
	if err != nil {
		return err
	}
	defer func(fh *os.File) {
		err := fh.Close()
		if err != nil {
			log.Errorf("Error closing file %s: %s\n", hostsFile, err)
		}
	}(fh)
	if strings.HasSuffix(hostsFile, view.CompressedSuffix) {
		zr, err := gzip.NewReader(fh)
		if err != nil {
			return err
		}
		defer func(zr *gzip.Reader) {
			err := zr.Close()
			if err != nil {
				log.Errorf("Error closing compressed file %s: %s\n", hostsFile, err)
			}
		}(zr)
		return hr.readCaptureServiceMap(zr)
	}
	return hr.readCaptureServiceMap(fh)
}

func (hr *hostsReaderImpl) GetServiceByIp(ip string) (*entities.ServiceAddress, error) {
	ise, err := hr.hostsCache.GetServiceAddressByIp(ip)
	if err == nil {
		return ise, nil // found in cache
	}
	// not in cache - ask db for an IP address
	svcAddr, err := hr.hostsCache.GetServiceAddress(ip, view.EmptyString, view.EmptyString, hr.captureId)
	if err != nil {
		return nil, err // unable to find in DB
	}
	return &svcAddr, nil // found in DB
}

func (hr *hostsReaderImpl) readCaptureServiceMap(fh io.Reader) error {
	scanner := bufio.NewScanner(fh)
	if scanner == nil {
		return fmt.Errorf("unable to create file scanner")
	}
	var reqRe = regexp.MustCompile(`^([a-f\d\.:]+)\s+({.+})\s*$`)
	for scanner.Scan() {
		ms := reqRe.FindSubmatch(scanner.Bytes())
		if ms != nil {
			svc, err := view.UnmarshalServiceView(ms[2])
			if err != nil {
				return err
			}
			_, err = hr.hostsCache.GetServiceAddress(string(ms[1]), svc.Name, svc.Version, hr.captureId)
			if err != nil {
				log.Debugf("unable to get service by ip %s: %s", string(ms[1]), err)
			}
		}
	}
	return nil
}

func (hr *hostsReaderImpl) Close() error {
	//return hr.hostsCache.
	return nil
}
