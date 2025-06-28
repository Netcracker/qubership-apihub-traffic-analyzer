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
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/Netcracker/qubership-apihub-traffic-analyzer/qubership-apihub-traffic-analyzer/db"
	"github.com/Netcracker/qubership-apihub-traffic-analyzer/qubership-apihub-traffic-analyzer/decoders"
	"github.com/Netcracker/qubership-apihub-traffic-analyzer/qubership-apihub-traffic-analyzer/entities"
	"github.com/Netcracker/qubership-apihub-traffic-analyzer/qubership-apihub-traffic-analyzer/repository"
	"github.com/Netcracker/qubership-apihub-traffic-analyzer/qubership-apihub-traffic-analyzer/view"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	log "github.com/sirupsen/logrus"
)

type CaptureReader interface {
	ReadCaptureDir(captureId string, workDir string) error
	ReadCaptureFile(captureId, inputFile string) (int, error)
	GetMetadataReader(captureId string) MetadataReader
	ReadHostsFile2(fileName, captureId string) error
	ReadHostsFile(fileName string) error
	Close() error
}

func NewCaptureReader(
	headers repository.HttpHeadersCache,
	packets repository.PacketCache,
	peers repository.ServiceAddressRepository,
	pdb db.ConnectionProvider,
	workDir string) CaptureReader {
	return &captureReaderImpl{
		headers:   headers,
		packets:   packets,
		peers:     peers,
		db:        pdb,
		hosts:     nil,
		workDir:   workDir,
		captureId: view.EmptyString,
	}
}

type captureReaderImpl struct {
	headers   repository.HttpHeadersCache
	packets   repository.PacketCache
	peers     repository.ServiceAddressRepository
	hosts     HostsReader
	db        db.ConnectionProvider
	workDir   string
	captureId string
}

func (cr *captureReaderImpl) ReadCaptureFile(captureId, fileName string) (int, error) {
	cr.captureId = captureId
	if strings.HasSuffix(fileName, view.CompressedSuffix) {
		cs, err := os.Open(fileName)
		if err != nil {
			return 0, err
		}
		defer func(cs *os.File) {
			err := cs.Close()
			if err != nil {
				log.Errorf("unable to close compressed file %s. Error: %v", path.Base(fileName), err)
			}
		}(cs)
		zr, err := gzip.NewReader(cs)
		if zr == nil || err != nil {
			if err != nil {
				log.Errorf("unable to uncompress file %s. Error: %v", path.Base(fileName), err)
			} else {
				return 0, fmt.Errorf("unable to uncompress file %s", path.Base(fileName))
			}
		}
		defer func(zr *gzip.Reader) {
			err := zr.Close()
			if err != nil {
				log.Errorf("unable to close decompressor for file %s. Error: %v", path.Base(fileName), err)
			}
		}(zr)
		tmpFileName := fileName + ".tmp"
		fth, err := os.Create(tmpFileName)
		if err != nil {
			log.Errorf("unable to create intermediate file '%s'. Error: %v", tmpFileName, err)
		} else {
			defer func(fth *os.File) {
				err := fth.Close()
				if err != nil {
					log.Errorf("unable to close intermedate file '%s'. Error: %v", tmpFileName, err)
				}
			}(fth)
			_, err = io.Copy(fth, zr)
			if err != nil {
				log.Errorf("unable to write uncompressed data. error: %v", err)
				return 0, err
			}
			return cr.readPackets(captureId, tmpFileName)
		}
	}
	return cr.readPackets(captureId, fileName)
}

func (cr *captureReaderImpl) ReadCaptureDir(captureId string, workDir string) error {
	fileInfo, err := os.Lstat(workDir)
	if err != nil {
		log.Errorf("unable to get info for work dir %s : %v", workDir, err.Error())
		return err
	}
	if fileInfo == nil {
		log.Errorf("unable to get info for work dir %s : file info is nil", workDir)
		return nil
	}
	if fileInfo.Mode()&os.ModeSymlink == os.ModeSymlink {
		realDir, err := os.Readlink(workDir)
		if err != nil {
			log.Errorf("unable to read link for work dir %s : %v", workDir, err.Error())
			return err
		}
		fileInfo, err := os.Lstat(realDir)
		if err != nil {
			log.Errorf("unable to get info for work dir %s : %v", workDir, err.Error())
			return err
		}
		if fileInfo.Mode()&os.ModeDir != os.ModeDir {
			log.Errorf("unable to work with not dir %s", realDir)
			return err
		}
		workDir = realDir
	}
	if cr.hosts == nil {
		log.Debugf("ReadCaptureDir: %s", cr.captureId)
		cr.hosts, err = NewHostsReader(workDir, cr.captureId, cr.db)
		if err != nil {
			log.Errorf("unable to establish hosts cache: %v", err)
			return err
		}
	}
	items, err := os.ReadDir(workDir)
	if err != nil {
		log.Fatalf("unable to read dir %s : %v", workDir, err.Error())
	}
	// fill service names
	for _, item := range items {
		if item.IsDir() {
			continue
		}
		name := item.Name()
		if !strings.HasPrefix(name, captureId) {
			continue
		}
		inputFilename := path.Join(workDir, name)
		if strings.HasSuffix(name, view.CompressedSuffix) {
			name = strings.TrimSuffix(name, view.CompressedSuffix)
		}
		if strings.HasSuffix(name, ".txt") {
			// service name list
			err = cr.ReadHostsFile2(inputFilename, captureId)
			if err != nil {
				log.Errorf("unable to read hosts file '%s'. Error: %v", item.Name(), err)
			} else {
				log.Printf("read hosts file '%s' successfully", item.Name())
			}
		}
		if strings.HasSuffix(name, ".json") {
			md := cr.GetMetadataReader(captureId)
			if md != nil {
				err = md.ReadFile(inputFilename)
				if err != nil {
					log.Errorf("unable to process metadata file '%s'. Error: %v", item.Name(), err)
				} else {
					log.Printf("metadata file '%s' processed successfully", item.Name())
				}
			} else {
				log.Errorf("unable to create metadata reader for file %s:%v", item.Name(), err)
			}
		}
	}
	// read captures
	for _, item := range items {
		if item.IsDir() {
			continue
		}
		name := item.Name()
		if !strings.HasPrefix(name, captureId) {
			continue
		}
		if strings.HasSuffix(name, view.CompressedSuffix) {
			name = strings.TrimSuffix(name, view.CompressedSuffix)
		}
		inputFilename := path.Join(workDir, item.Name())
		if strings.HasSuffix(name, ".pcap") {
			packetCount, err := cr.ReadCaptureFile(captureId, inputFilename)
			if err != nil {
				log.Errorf("unable to process capture file '%s'. Error: %v", item.Name(), err)
			} else {
				log.Debugf("capture file '%s' read successfully, packets processed %d", item.Name(), packetCount)
			}
		}
	}
	return nil
}

func (cr *captureReaderImpl) GetMetadataReader(captureId string) MetadataReader {
	return NewMetadataReader(cr.db, captureId)
}

func (cr *captureReaderImpl) ReadHostsFile2(fileName, captureId string) error {
	if cr.hosts == nil {
		hosts, err := NewHostsReader(cr.workDir, captureId, cr.db)
		if err != nil {
			return err
		}
		cr.hosts = hosts
	}
	return cr.hosts.Read(fileName)
}

func (cr *captureReaderImpl) ReadHostsFile(fileName string) error {
	if cr.hosts == nil {
		log.Debugf("ReadHostsFile: %s", cr.captureId)
		hosts, err := NewHostsReader(cr.workDir, cr.captureId, cr.db)
		if err != nil {
			return err
		}
		cr.hosts = hosts
	}
	return cr.hosts.Read(fileName)
}

func (cr *captureReaderImpl) readPackets(captureId, fileName string) (int, error) {
	if cr.hosts == nil {
		return -1, fmt.Errorf("no hosts file for capture: %s", captureId)
	}
	handle, err := pcap.OpenOffline(fileName)
	if err != nil {
		log.Printf("unable to open offline capture %s. Error: %v", fileName, err)
		return 0, err
	}
	defer handle.Close()
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	i := -1
	df := decoders.DecodeFeedback{}
	processedPackets := 0
	httpPackets := 0
	for packet := range packetSource.Packets() {
		i++
		var (
			ls1                           = view.EmptyString
			ipPacket   *decoders.IPPacket = nil
			pktType                       = layers.EthernetTypeLLC
			ipPayLoad  []byte             = nil
			packetData []byte             = nil
		)
		packetData = packet.Data()
		if len(packetData) < decoders.ETHSize {
			continue // packet too small - skip
		}
		sll, err := decoders.DetectAndParseSll(packet.Data())
		if err != nil {
			if errors.Is(err, decoders.ErrorNotAnIpPacket) {
				continue // skip
			}
			etl := layers.Ethernet{}
			err = etl.DecodeFromBytes(packet.Data(), &df)
			if err != nil {
				log.Errorf("unable to decode ETH packet %d. Error: %v", i, err)
				continue
			} else {
				ipPayLoad = etl.Payload
				pktType = etl.EthernetType

			}
		} else {
			ipPayLoad = sll.Payload
			pktType = sll.EthernetType
		}
		if ipPayLoad != nil {
			switch pktType {
			case layers.EthernetTypeIPv4:
				ipPacket, err = decoders.DecodeIp(ipPayLoad, false, i)
			case layers.EthernetTypeIPv6:
				ipPacket, err = decoders.DecodeIp(ipPayLoad, true, i)
			default:
				{
					log.Tracef("unsupported packet type %d at %d", pktType, i)
					continue
				}
			}
		} else {
			log.Errorf("NIL packet with type %d at %d", pktType, i)
			continue
		}
		if err != nil {
			log.Errorf("unable to decode IP packet %d. Error: %v", i, err)
			continue
		}
		if ipPacket == nil {
			log.Errorf("no error, but IP packet is nil at %d", i)
			continue
		}
		if ipPacket.SrcIP == view.EmptyString && ipPacket.DstIP == view.EmptyString {
			log.Errorf("empty IP addresses in packet %d (IP layer length:%d, TCP layer length:%d)", i, len(ipPayLoad), len(ipPacket.TCP))
			continue
		}
		peers := make([]entities.ServiceAddress, 2)
		sourceService, err := cr.hosts.GetServiceByIp(ipPacket.SrcIP)
		if err == nil {
			peers[view.SourcePeer] = *sourceService
		} else {
			log.Debugf("source ip address %s not found: %v", ipPacket.SrcIP, err)
		}
		destService, err := cr.hosts.GetServiceByIp(ipPacket.DstIP)
		if err == nil {
			peers[view.DestPeer] = *destService
		} else {
			log.Debugf("dest ip address %s not found: %v", ipPacket.DstIP, err)
		}
		if ipPacket.TCP == nil {
			continue
		}
		tcp := layers.TCP{}
		err = tcp.DecodeFromBytes(ipPacket.TCP, &df)
		if err != nil {
			log.Tracef("unable to decode TCP packet %d. Error: %v\n", i, err)
			continue // not a TCP packet
		}
		headers := make(map[string]string)
		tcpPayLoad := tcp.Payload
		reqPath := view.EmptyString
		reqMethod := view.EmptyString
		if tcpPayLoad != nil {
			pt, parsedFields := decoders.DetectHttp(tcpPayLoad)
			bHttp := false
			switch pt {
			case decoders.PTHttpReq:
				{
					reqPath = fmt.Sprint(parsedFields[decoders.RequestPath])
					reqMethod = fmt.Sprint(parsedFields[decoders.RequestMethod])
					httpPackets++
					bHttp = true
					ls1 = string(tcpPayLoad)
					br := bytes.NewReader(tcpPayLoad)
					if br != nil {
						req, err := http.ReadRequest(bufio.NewReader(br))
						if err == nil {
							if req != nil {
								// headers
								for headerName, headerValue := range req.Header {
									headerString := strings.Join(headerValue, "\n")
									if headerString != view.EmptyString {
										headers[headerName] = headerString
									}
								}
								// request!
								reqMethod = req.Method
								reqPath = req.URL.Path
								bodyResult := decoders.BodyToString(req.Body, true)
								if bodyResult.Err == nil {
									if bodyResult.Body != nil {
										if ls1 == view.EmptyString {
											ls1 = string(bodyResult.Body)
										}
									}
								} else {
									if !strings.HasSuffix(bodyResult.Err.Error(), "unexpected EOF") {
										log.Tracef("unable to decode request %d body: %v", i, bodyResult.Err)
									}
								}
							} else {
								log.Tracef("request is nil")
							}
						} else {
							log.Tracef("unable to detect http request at %d:%v", i, err)
						}
					} else {
						log.Tracef("request bufreade is nil '%s' at %d:%v", ls1, i, err)
					}
				}
			case decoders.PTHttpResp:
				{
					httpPackets++
					reqMethod = fmt.Sprintf("%s %s", parsedFields[decoders.ResponseStatus], parsedFields[decoders.ResponseStatusText])
					bHttp = true
					ls1 = string(tcpPayLoad)
					br := bytes.NewReader(tcpPayLoad)
					resp, err := http.ReadResponse(bufio.NewReader(br), nil)
					if err == nil {
						if resp != nil {
							reqMethod = resp.Status
							// headers
							for headerName, headerValue := range resp.Header {
								headerString := strings.Join(headerValue, "\n")
								if len(headerString) > 0 {
									headers[headerName] = headerString
								}
							}
							// response!
							if resp.Request != nil {
								reqPath = resp.Request.URL.Path
							}
							bodyResult := decoders.BodyToString(resp.Body, resp.Uncompressed)
							if bodyResult.Err == nil {
								if bodyResult.Body != nil {
									if ls1 == view.EmptyString {
										ls1 = string(bodyResult.Body)
									}
								}
							} else {
								if !strings.HasSuffix(bodyResult.Err.Error(), "unexpected EOF") {
									log.Tracef("unable to decode response %d body: %v", i, bodyResult.Err)
								}
							}
						} else {
							log.Tracef("response is nil")
						}
					} else {
						if !errors.Is(err, io.ErrUnexpectedEOF) {
							log.Tracef("unable to detect http response at %d:%v", i, err)
						}
					}
				}
			case decoders.PTHttp:
				{
					pos := 0
					if iPos, bff := parsedFields[decoders.RHOffset]; bff {
						if pos, bff = iPos.(int); !bff {
							pos = 0
						} else {
							if pos < 0 {
								pos = 0
							}
						}
					}
					ls1 = string(tcpPayLoad[pos:])
				}
			case decoders.PTNotHttp:
				break // log.Debugf("not a HTTP packet %d\n", i)
			default:
				log.Debug("unhandled default case")
			}
			if ls1 == view.EmptyString && reqMethod == view.EmptyString && reqPath == view.EmptyString {
				if bHttp {
					log.Debugf("http packet %d rejected (%s,%s,%s)", i, reqMethod, reqPath, ls1)
				}
				continue
			}
			p2s := entities.ParsedPacket{
				Peers:         nil,
				Ports:         make([]int, 2),
				Timestamp:     packet.Metadata().Timestamp,
				SeqNo:         int(tcp.Seq),
				AckNo:         int(tcp.Ack),
				Payload:       tcpPayLoad,
				StrPayload:    ls1,
				RequestPath:   reqPath,
				RequestMethod: reqMethod,
				Headers:       headers,
			}
			p2s.Headers = headers
			p2s.Peers = peers
			p2s.Ports[view.SourcePeer] = int(tcp.SrcPort)
			p2s.Ports[view.DestPeer] = int(tcp.DstPort)
			err = cr.packets.StorePacket(p2s, cr.headers, captureId)
			if err == nil {
				processedPackets++
			} else {
				log.Debugf("unable to store packet %d:%v", i, err)
			}
		}
	}
	log.Debugf("total packets: %d, HTTP packets: %d,  processed packets: %d for capture %s", i, httpPackets, processedPackets, captureId)
	return processedPackets, nil
}

// Close
// performs an attempt to free underlying resources
func (cr *captureReaderImpl) Close() error {
	if cr.packets != nil {
		cr.packets.Close()
	}
	if cr.headers != nil {
		cr.headers.Close()
	}
	if cr.peers != nil {
		cr.peers.Close()
	}
	if cr.hosts != nil {
		return cr.hosts.Close()
	}
	return nil
}
