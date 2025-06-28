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

type PacketCache interface {
	GetPacketCount(captureId string) (int, error)
	StorePacket(packet entities.ParsedPacket, headersCache HttpHeadersCache, captureId string) error
	Close()
}

type packetCacheImpl struct {
	db          db.ConnectionProvider
	addrRepo    ServiceAddressRepository
	headersRepo HttpHeadersCache
}

func NewPacketCache(db db.ConnectionProvider, peersCache ServiceAddressRepository, headersCache HttpHeadersCache) PacketCache {
	return &packetCacheImpl{
		db:          db,
		addrRepo:    peersCache,
		headersRepo: headersCache,
	}
}

func (p *packetCacheImpl) StorePacket(packet entities.ParsedPacket, headersCache HttpHeadersCache, captureId string) error {
	var err error = nil
	servicePacket := entities.MakeDbPacket(packet, captureId)
	servicePacket.PacketId, err = p.acquirePacketId(servicePacket)
	if err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			result := servicePacket
			_, err = p.db.GetConnection().Model(&servicePacket).Returning("packet_id").Insert(&result)
			if err == nil {
				servicePacket.PacketId = result.PacketId
			} else {
				err = fmt.Errorf("unable to insert packet: %v (%d %s/%d %s)", err, result.SourceId, packet.Peers[view.SourcePeer].Address, result.DestId, packet.Peers[view.DestPeer].Address)
			}
		} else {
			err = fmt.Errorf("unable to get packet: %v", err)
		}
		if err != nil {
			return err
		}
	}

	for k := range packet.Headers {
		hid, err := headersCache.GetPrimaryKeyValue(k, packet.Headers[k])
		if err != nil || hid == view.EmptyString {
			log.Debugf("unable to get header id for header %s: %v", k, err)
			continue
		}
		ph := entities.PacketHeader{
			HeaderId: hid,
			PacketId: servicePacket.PacketId}
		pr := ph
		err = p.db.GetConnection().Model(&pr).Where("header_Id=? and packet_id=?", ph.HeaderId, ph.PacketId).Select()
		if err != nil {
			if errors.Is(err, pg.ErrNoRows) {
				pr = ph
				_, err = p.db.GetConnection().Model(&ph).Insert(&pr)
				if err != nil {
					log.Debugf("unable to insert packet header: %v for packet %d:%s", err, servicePacket.PacketId, hid)
					err = fmt.Errorf("unable to insert packet header: %v for packet %d:%s", err, servicePacket.PacketId, hid)
				}
			} else {
				if err != nil {
					log.Debugf("unable to query packet header: %v for packet %d", err, servicePacket.PacketId)
					err = fmt.Errorf("unable to query packet header: %v for packet %d", err, servicePacket.PacketId)
				}
			}
		}
	}
	return nil
}

func (p *packetCacheImpl) acquirePacketId(packet entities.ServicePacket) (int, error) {
	result := new(entities.ServicePacket)
	err := p.db.GetConnection().Model(result).Where(
		//"source_id=? and source_port=? and dest_id=? and dest_port=? and Seq_No=? and Ack_no=? and Time_Stamp=? and Body=? and capture_id=?",
		//packet.SourceId, packet.SourcePort, packet.DestId, packet.DestPort, packet.SeqNo, packet.AckNo, packet.TimeStamp, packet.Body, packet.CaptureId).First()
		"source_id=? and source_port=? and dest_id=? and dest_port=? and Seq_No=? and Ack_no=? and Time_Stamp=? and capture_id=?",
		packet.SourceId, packet.SourcePort, packet.DestId, packet.DestPort, packet.SeqNo, packet.AckNo, packet.TimeStamp, packet.CaptureId).First()
	return result.PacketId, err
}

func (p *packetCacheImpl) GetPacketCount(captureId string) (int, error) {
	mr := new(entities.ServicePacket)
	recCount, err := p.db.GetConnection().Model(mr).Where("capture_id=?", captureId).Count()
	if err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			return 0, nil
		}
	}
	return recCount, err
}

func (p *packetCacheImpl) Close() {
}
