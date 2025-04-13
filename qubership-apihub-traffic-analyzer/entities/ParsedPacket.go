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

package entities

import (
	"time"

	"github.com/Netcracker/qubership-apihub-traffic-analyzer/qubership-apihub-traffic-analyzer/view"
)

type ServicePacket struct {
	tableName struct{} `pg:"service_packets, alias:service_packets"`

	PacketId      int       `pg:"packet_id,pk,type:bigint"`
	SourceId      int       `pg:"source_id,type:bigint"`
	SourcePort    int       `pg:"source_port,type:bigint"`
	DestId        int       `pg:"dest_id,type:bigint"`
	DestPort      int       `pg:"dest_port,type:bigint"`
	TimeStamp     time.Time `pg:"time_stamp,type:timestamptz"`
	SeqNo         int       `pg:"seq_no,type:bigint"`
	AckNo         int       `pg:"ack_no,type:bigint"`
	Body          string    `pg:"body,type:text"`
	CaptureId     string    `pg:"capture_id,type:varchar"`
	RequestPath   string    `pg:"request_path,type:varchar"`
	RequestMethod string    `pg:"request_method,type:varchar"`
}

type ParsedPacket struct {
	Peers         []ServiceAddress
	Ports         []int
	Timestamp     time.Time
	SeqNo         int
	AckNo         int
	Payload       []byte
	StrPayload    string
	ServiceName   string
	Headers       map[string]string
	RequestPath   string
	RequestMethod string
}

func MakeDbPacket(p ParsedPacket, captureId string) ServicePacket {
	return ServicePacket{
		SourceId:      p.Peers[view.SourcePeer].Id,
		SourcePort:    p.Ports[view.SourcePeer],
		DestId:        p.Peers[view.DestPeer].Id,
		DestPort:      p.Ports[view.DestPeer],
		TimeStamp:     p.Timestamp,
		SeqNo:         p.SeqNo,
		AckNo:         p.AckNo,
		Body:          p.StrPayload,
		CaptureId:     captureId,
		RequestPath:   p.RequestPath,
		RequestMethod: p.RequestMethod,
	}
}
