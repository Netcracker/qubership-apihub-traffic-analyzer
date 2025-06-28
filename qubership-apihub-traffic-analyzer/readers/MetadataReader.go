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
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/Netcracker/qubership-apihub-traffic-analyzer/qubership-apihub-traffic-analyzer/db"
	"github.com/Netcracker/qubership-apihub-traffic-analyzer/qubership-apihub-traffic-analyzer/entities"
	"github.com/Netcracker/qubership-apihub-traffic-analyzer/qubership-apihub-traffic-analyzer/view"
	log "github.com/sirupsen/logrus"
)

type MetadataReader interface {
	ReadBytes(metadata []byte) error
	ReadFile(metadataFile string) error
}

type metadataReaderImpl struct {
	db        db.ConnectionProvider
	captureId string
}

func NewMetadataReader(conn db.ConnectionProvider, captureId string) MetadataReader {
	return &metadataReaderImpl{
		db:        conn,
		captureId: captureId,
	}
}

func (md *metadataReaderImpl) ReadBytes(metadata []byte) error {
	p := entities.CaptureMetadata{
		CaptureId: md.captureId,
		Metadata:  string(metadata),
	}
	_, err := md.db.GetConnection().Model(&p).Insert(&p)
	if err != nil {
		_, err = md.db.GetConnection().Model(&p).Where("capture_id = ?", p.CaptureId).Update()
	}
	return err
}

func (md *metadataReaderImpl) ReadFile(metadataFile string) error {
	if strings.HasSuffix(metadataFile, view.CompressedSuffix) {
		fh, err := os.Open(metadataFile)
		if err != nil {
			return err
		}
		defer func(fh *os.File) {
			err := fh.Close()
			if err != nil {
				log.Errorf("Error closing metadata file %s: %s", metadataFile, err)
			}
		}(fh)
		zr, err := gzip.NewReader(fh)
		if err != nil {
			return err
		}
		defer func(zr *gzip.Reader) {
			err := zr.Close()
			if err != nil {
				log.Errorf("Error closing compressed metadata file %s: %s", metadataFile, err)
			}
		}(zr)
		return md.readFileToBytes(zr)
	}
	metadata, err := os.ReadFile(metadataFile)
	if metadata != nil && err == nil {
		return md.ReadBytes(metadata)
	} else {
		if metadata == nil {
			err = fmt.Errorf("no data read from %s", metadataFile)
		}
	}
	return err

}

func (md *metadataReaderImpl) readFileToBytes(input io.Reader) error {
	buf := make([]byte, 32767)
	nr, err := input.Read(buf)
	if err != nil {
		if nr < 1 {
			return err
		}
	}
	if nr > 0 {
		err = md.ReadBytes(buf[0:nr])
	} else {
		err = io.ErrShortBuffer
	}
	return err
}
