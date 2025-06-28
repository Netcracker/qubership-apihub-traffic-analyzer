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

package service

import (
	"context"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/Netcracker/qubership-apihub-traffic-analyzer/qubership-apihub-traffic-analyzer/readers"
	"github.com/Netcracker/qubership-apihub-traffic-analyzer/qubership-apihub-traffic-analyzer/view"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	log "github.com/sirupsen/logrus"
)

const TableName = "PacketCaptures"

type CloudStorage interface {
	ProcessCaptureFiles(captureId string, rdr readers.CaptureReader) (int, error)
	DeleteCaptureFiles(captureId string) (int, error)
	CleanupCaptureFiles() (int, error)
}

type cloudStorage struct {
	config      view.MinioStorageCreds
	lock        sync.Mutex
	minioClient *minioClient
	bannerTime  time.Time
}

type minioClient struct {
	client *minio.Client
	error  error
}

var ErrorMinioClient = errors.New("S3/Minio client was no initialized properly")

// public interface functions

// NewCloudStorage
// initialises a new interface instance
func NewCloudStorage(config view.MinioStorageCreds) (CloudStorage, error) {
	var err error = nil
	s3Client := createMinioClient(&config)
	if s3Client == nil {
		err = ErrorMinioClient
	}
	return &cloudStorage{
		config:      config,
		lock:        sync.Mutex{},
		minioClient: s3Client,
		bannerTime:  time.Now(),
	}, err
}

// ProcessCaptureFiles
// downloads capture files from S3/Minio to local
func (s3 *cloudStorage) ProcessCaptureFiles(captureId string, rdr readers.CaptureReader) (int, error) {
	receivedCount := 0
	ctx := context.Background()
	var addressLists []minio.ObjectInfo
	var captureFiles []minio.ObjectInfo
	var captureMetadata minio.ObjectInfo
	// iterate and sort objects
	if s3.minioClient == nil || s3.minioClient.client == nil {
		log.Errorln(ErrorMinioClient)
		return receivedCount, ErrorMinioClient
	}
	opts := minio.ListObjectsOptions{
		WithVersions: false,
		WithMetadata: true,
		Prefix:       TableName + "/" + captureId,
		Recursive:    true,
		MaxKeys:      1000,
	}
	if time.Now().After(s3.bannerTime) {
		log.Printf("loading capture '%s'", captureId)
		log.Debugf("bucket    :%s", s3.config.BucketName)
		log.Debugf("active    :%v", s3.config.IsActive)
		log.Debugf("endpoint  :%s", s3.config.Endpoint)
		log.Debugf("cert      :%s...", s3.config.Crt[:30])
		log.Debugf("access key:%s", s3.config.AccessKeyId)
		log.Debugf("secret    :%s", s3.config.SecretAccessKey)
		log.Debugf("production:%v", s3.config.ProductionMode)
		log.Debugf("work dir  :%s", s3.config.WorkDir)
		s3.bannerTime = time.Now().Add(time.Minute)
	}
	// receive a channel to objects
	s3Objects := s3.minioClient.client.ListObjects(context.Background(), s3.config.BucketName, opts)
	for objectInfo := range s3Objects {
		if objectInfo.Err != nil {
			log.Errorf("unable to list file %s from S3/minio: %v", objectInfo.Key, objectInfo.Err)
			continue
		}
		name := objectInfo.Key
		if strings.HasSuffix(objectInfo.Key, view.CompressedSuffix) {
			name = strings.TrimSuffix(objectInfo.Key, view.CompressedSuffix)
		}
		if strings.HasSuffix(name, view.AddressListSuffix) {
			addressLists = append(addressLists, objectInfo)
			continue
		}
		if strings.HasSuffix(name, view.CaptureSuffix) {
			captureFiles = append(captureFiles, objectInfo)
			continue
		}
		if strings.HasSuffix(name, view.MetadataSuffix) {
			captureMetadata = objectInfo
			continue
		}
	}
	if captureMetadata.Key == view.EmptyString || len(addressLists) == 0 || len(captureFiles) == 0 {
		//return receivedCount, fmt.Errorf("capture data not found at S3/Minio")
		log.Warnf("no metadata found for capture id %s", captureId)
	} else {
		// loading capture metadata
		md := rdr.GetMetadataReader(captureId)
		if md == nil {
			return receivedCount, fmt.Errorf("metadata reader is nil")
		}
		if strings.HasSuffix(captureMetadata.Key, view.CompressedSuffix) {
			localFileName, err := s3.getObject(ctx, &captureMetadata)
			if err != nil {
				return receivedCount, fmt.Errorf("unable to load capture metadata from S3/minio: %v", err)
			}
			// put loaded metadata into DB
			err = md.ReadFile(localFileName)
			if err != nil {
				return receivedCount, err
			}
		} else {
			// put loaded metadata into DB
			err := s3.getMetadataDirect(ctx, &captureMetadata, md)
			if err != nil {
				return receivedCount, err
			}
		}
		receivedCount++
	}
	// loading capture address cache
	for ai, addressList := range addressLists {
		localFileName, err := s3.getObject(ctx, &addressList)
		if err != nil {
			return receivedCount, fmt.Errorf("unable to load address list file %d.%s from S3/minio: %v", ai, addressList.Key, err)
		}
		// loads address resolving cache item
		tStart := time.Now()
		err = rdr.ReadHostsFile2(localFileName, captureId)
		tDiff := time.Now().Sub(tStart)
		log.Debugf("Reading hosts file %d in %v", ai+1, tDiff)
		if err != nil {
			return receivedCount, fmt.Errorf("unable to parse address list file %d.%s from S3/minio: %v", ai, addressList.Key, err)
		}
		receivedCount++
	}
	// loading capture data
	packetCount := 0
	for ci, captureFile := range captureFiles {
		localFileName, err := s3.getObject(ctx, &captureFile)
		if err != nil {
			return receivedCount, fmt.Errorf("unable to load capture file %d.%s from S3/minio: %v", ci, captureFile.Key, err)
		}
		// loads capture data from local file
		tStart := time.Now()
		count, err := rdr.ReadCaptureFile(captureId, localFileName)
		tDiff := time.Now().Sub(tStart)
		log.Debugf("Reading capture file %d in %v", ci+1, tDiff)
		receivedCount++
		if err != nil {
			return receivedCount, fmt.Errorf("unable to process capture file %d.%s from S3/minio: %v", ci, captureFile.Key, err)
		}
		packetCount += count
	}
	log.Printf("files read: %d, packets processed: %d", receivedCount, packetCount)
	return receivedCount, nil
}

func (s3 *cloudStorage) DeleteCaptureFiles(captureId string) (int, error) {
	deletedCount := 0
	// iterate and sort objects
	if s3.minioClient == nil || s3.minioClient.client == nil {
		log.Errorln(ErrorMinioClient)
		return deletedCount, ErrorMinioClient
	}
	if time.Now().After(s3.bannerTime) {
		log.Debugf("Deleting capture '%s'", captureId)
		log.Debugf("bucket    :%s", s3.config.BucketName)
		log.Debugf("active    :%v", s3.config.IsActive)
		log.Debugf("production:%v", s3.config.ProductionMode)
		s3.bannerTime = time.Now().Add(time.Minute)
	}
	// List all objects from a bucket-name with a matching prefix.
	opts := minio.ListObjectsOptions{
		WithVersions: false,
		WithMetadata: true,
		Prefix:       TableName + "/" + captureId,
		Recursive:    true,
		MaxKeys:      1000,
	}
	for objectInfo := range s3.minioClient.client.ListObjects(context.Background(), s3.config.BucketName, opts) {
		if objectInfo.Err != nil {
			log.Debugf("unable to list object:%v", objectInfo.Err)
		}
		log.Debugf("deleting object:%s", objectInfo.Key)
		if strings.Contains(objectInfo.Key, captureId) {
			// delete listed object
			err := s3.minioClient.client.RemoveObject(context.Background(), s3.config.BucketName, objectInfo.Key, minio.RemoveObjectOptions{})
			if err != nil {
				log.Debugf("unable to delete object '%s':%v", objectInfo.Key, err)
			} else {
				deletedCount++
			}
		}
	}
	return deletedCount, nil
}

const (
	MinFileCount int32 = 3
	HasNothing   int32 = 0
	HasMetadata  int32 = 1
	HasPackets         = HasMetadata * 2
	HasAddresses       = HasPackets * 2
)

type CaptureStatus struct {
	FileCount int32
	Flags     int32
}
type captureMapKey [36]byte

func CanBeDeleted(status *CaptureStatus) bool {
	return status.FileCount < MinFileCount &&
		(status.Flags&HasMetadata != HasMetadata ||
			status.Flags&HasPackets != HasPackets ||
			status.Flags&HasAddresses != HasAddresses)
}

func (s3 *cloudStorage) CleanupCaptureFiles() (int, error) {
	deletedCount := 0
	// iterate and sort objects
	if s3.minioClient == nil || s3.minioClient.client == nil {
		log.Errorln(ErrorMinioClient)
		return deletedCount, ErrorMinioClient
	}
	if time.Now().After(s3.bannerTime) {
		log.Debugf("Clean up S3/Minio bucket '%s'", s3.config.BucketName)
		log.Debugf("active    :%v", s3.config.IsActive)
		log.Debugf("production:%v", s3.config.ProductionMode)
		s3.bannerTime = time.Now().Add(time.Minute)
	}
	// List all objects from a bucket-name with a matching prefix.
	opts := minio.ListObjectsOptions{
		WithVersions: false,
		WithMetadata: true,
		Prefix:       TableName + "/",
		Recursive:    true,
	}
	captures := make(map[captureMapKey]CaptureStatus)
	// expecting around 200000 records (36+8 bytes) here
	log.Debugf("collecting captures...")
	for objectInfo := range s3.minioClient.client.ListObjects(context.Background(), s3.config.BucketName, opts) {
		if objectInfo.Err != nil {
			log.Debugf("unable to list object:%v", objectInfo.Err)
		}
		name := objectInfo.Key
		capIdx := strings.Index(name, "_")
		if capIdx < 1 {
			continue // skip files with improper names
		}
		if strings.HasSuffix(name, view.CompressedSuffix) {
			name = name[:len(name)-len(view.CompressedSuffix)]
		}
		captureId := captureMapKey([]byte(name[:capIdx]))
		capStat, found := captures[captureId]
		if !found {
			capStat = CaptureStatus{
				FileCount: 1,
				Flags:     HasNothing,
			}
		} else {
			capStat.FileCount++
		}
		if strings.HasSuffix(name, view.MetadataSuffix) {
			capStat.Flags |= HasMetadata
		} else if strings.HasSuffix(name, view.CaptureSuffix) {
			capStat.Flags |= HasPackets
		} else if strings.HasSuffix(name, view.AddressListSuffix) {
			capStat.Flags |= HasAddresses
		}
		captures[captureId] = capStat
	}
	log.Debugf("captures collected: %d. deleting ...", len(captures))
	for captureId, capStat := range captures {
		if CanBeDeleted(&capStat) {
			capId := string(captureId[:])
			dc, err := s3.DeleteCaptureFiles(capId)
			if err != nil {
				log.Debugf("unable to delete capture '%s': %v", capId, err)
			} else {
				deletedCount += dc
			}
		}
	}
	return deletedCount, nil
}

// internal functions

// mustGetSystemCertPool
// acquires certification pool
func mustGetSystemCertPool() *x509.CertPool {
	pool, err := x509.SystemCertPool()
	if err != nil {
		return x509.NewCertPool()
	}
	return pool
}

// createMinioClient
// creates minio instance
func createMinioClient(minioCredentials *view.MinioStorageCreds) *minioClient {
	if !minioCredentials.IsActive {
		log.Warnf("S3/Minio cloud storage is not active")
		return nil // inactive storage does not require full fledged client
	}
	client := new(minioClient)
	tr, err := minio.DefaultTransport(true)
	if err != nil {
		log.Warnf("error creating the minio connection: error creating the default transport layer: %v", err)
		client.error = err
		return client
	}
	crt, err := os.CreateTemp("", "minio.cert")
	if err != nil {
		log.Warn("unable to create temporary certificate storage:%v", err)
		client.error = err
		return client
	}
	decodeSamlCert, err := base64.StdEncoding.DecodeString(minioCredentials.Crt)
	if err != nil {
		log.Warnf("unable to decode cert: %v", err)
		client.error = err
		return client
	}

	_, err = crt.WriteString(string(decodeSamlCert))
	if err != nil {
		log.Warn(err.Error())
		client.error = err
	}
	rootCAs := mustGetSystemCertPool()
	data, err := os.ReadFile(crt.Name())
	if err == nil {
		rootCAs.AppendCertsFromPEM(data)
	}
	tr.TLSClientConfig.RootCAs = rootCAs

	minioClient, err := minio.New(minioCredentials.Endpoint, &minio.Options{
		Creds:     credentials.NewStaticV4(minioCredentials.AccessKeyId, minioCredentials.SecretAccessKey, ""),
		Secure:    true,
		Transport: tr,
	})
	if err != nil {
		log.Warn(err.Error())
		client.error = err
		return client
	}
	log.Infof("MINIO instance initialized")
	client.client = minioClient
	return client
}

// bucketExists
// check whether S3/Minio bucket exists or not
func bucketExists(ctx context.Context, minioClient *minio.Client, bucketName string) (bool, error) {
	exists, err := minioClient.BucketExists(ctx, bucketName)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// createBucketIfNotExists
// creates bucket if it does not exist
func (s3 *cloudStorage) createBucketIfNotExists(ctx context.Context) error {
	exists, err := bucketExists(ctx, s3.minioClient.client, s3.config.BucketName)
	if err != nil {
		return err
	}
	if exists {
		log.Infof(fmt.Sprintf("Minio bucket - %s exists", s3.config.BucketName))
	} else {
		err = s3.minioClient.client.MakeBucket(ctx, s3.config.BucketName, minio.MakeBucketOptions{})
		if err != nil {
			return err
		}
		exists, err = bucketExists(ctx, s3.minioClient.client, s3.config.BucketName)
		if err != nil {
			return err
		}
		if exists {
			log.Infof(fmt.Sprintf("Minio bucket - %s was created", s3.config.BucketName))
		}
	}
	return nil
}

// getObject
// get object (file) from bucket to a local file system
func (s3 *cloudStorage) getObject(ctx context.Context, csObjectInfo *minio.ObjectInfo) (string, error) {
	csObject, err := s3.minioClient.client.GetObject(ctx, s3.config.BucketName, csObjectInfo.Key, minio.GetObjectOptions{})
	if err != nil {
		return view.EmptyString, err
	}
	fullPath := path.Join(s3.config.WorkDir, path.Base(csObjectInfo.Key))
	ofh, err := os.Create(fullPath)
	if err != nil {
		return view.EmptyString, err
	}
	defer func(ofh *os.File) {
		err := ofh.Close()
		if err != nil {
			log.Errorf("unable to close file %s: %v", csObjectInfo.Key, err)
		}
	}(ofh)
	fileSize := csObjectInfo.Size
	tryCount := 3
	for fileSize > 0 {
		buf := make([]byte, 8192)
		nr := 0
		nr, err = csObject.Read(buf)
		if err != nil && nr < 1 {
			err := os.Remove(fullPath)
			if err != nil {
				log.Errorf("unable to remove file %s: %v", csObjectInfo.Key, err)
			}
			break
		}
		if nr > 0 {
			nw := 0
			nw, err = ofh.Write(buf[:nr])
			if err != nil {
				err := os.Remove(fullPath)
				if err != nil {
					log.Errorf("unable to remove file %s: %v", csObjectInfo.Key, err)
				}
				break
			}
			if nw != nr {
				err := os.Remove(fullPath)
				if err != nil {
					log.Errorf("unable to remove file %s: %v", csObjectInfo.Key, err)
				}
				err = io.ErrShortWrite
				break
			}
		} else {
			tryCount--
			if tryCount <= 0 {
				err = io.ErrShortBuffer
				break
			}
		}
		fileSize -= int64(nr)
	}
	return fullPath, err
}

func (s3 *cloudStorage) getMetadataDirect(ctx context.Context, csMdInfo *minio.ObjectInfo, rdr readers.MetadataReader) error {
	csObject, err := s3.minioClient.client.GetObject(ctx, s3.config.BucketName, csMdInfo.Key, minio.GetObjectOptions{})
	if err != nil {
		return err
	}
	nSize := int(csMdInfo.Size)
	if nSize > 0 {
		buf := make([]byte, nSize)
		nr, err := csObject.Read(buf)
		if err != nil && nr < 1 {
			return err
		}
		if nr != nSize {
			return io.ErrShortBuffer
		}
		err = rdr.ReadBytes(buf)
	}
	return err
}
