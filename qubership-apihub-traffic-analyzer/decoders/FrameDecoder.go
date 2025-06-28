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

package decoders

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/google/gopacket/layers"
)

type IPVersion int

const (
	// SLLSize Linux cooked-mode capture (SLL) minimal packet size
	SLLSize = 16
	// ETHSize ethernet header size to drop packet without a proper data inside
	ETHSize = 14
	// IPvUnknown IP version is not known
	IPvUnknown IPVersion = 0
	// IPv4 IP version 4 (32 bit/four-bytes IPs)
	IPv4 = 4
	// IPv6 IP version 6 (128-bit/eight bytes IPs)
	IPv6 = 6
)

type IPPacket struct {
	// Version an IP version
	Version IPVersion
	// SrcIP Source IP address
	SrcIP string
	// DstIP Destination IP address
	DstIP string
	// TCP a tcp data (IP payload)
	TCP []byte
	// TCPExpected true if IP completely parsed
	TCPExpected bool
}

// ErrorNotAnIpPacket no TCP/IP data in Linux SLL packet
var ErrorNotAnIpPacket error = errors.New("not an IP packet")

// DecodeIp
// decodes IPv4 and IPv6 packet from bytes into common structure
func DecodeIp(data []byte, ipv6 bool, i int) (*IPPacket, error) {
	var err error = nil
	ip := IPPacket{
		Version:     IPvUnknown,
		SrcIP:       "",
		DstIP:       "",
		TCP:         nil,
		TCPExpected: false,
	}
	df := DecodeFeedback{}
	if ipv6 {
		ipv6 := layers.IPv6{}
		err = ipv6.DecodeFromBytes(data, &df)
		if err != nil {
			return nil, fmt.Errorf("unable to decode IPv6 packet %d. Error: %v\n", i, err)
		}
		ip.SrcIP = ipv6.SrcIP.String()
		ip.DstIP = ipv6.DstIP.String()
		ip.TCP = ipv6.Payload
		ip.Version = IPv6
		ip.TCPExpected = true
	} else {
		ipv4 := layers.IPv4{}
		err = ipv4.DecodeFromBytes(data, &df)
		if err != nil {
			return nil, fmt.Errorf("unable to decode IPv4 packet %d. Error: %v\n", i, err)
		}
		ip.SrcIP = ipv4.SrcIP.String()
		ip.DstIP = ipv4.DstIP.String()
		ip.Version = IPv4
		if ipv4.Protocol == layers.IPProtocolTCP {
			ip.TCP = ipv4.Payload
			ip.TCPExpected = true
		}
	}
	return &ip, err
}

// DetectAndParseSll
// an error exposing parser of Linux cooked-mode capture (SLL) packet
func DetectAndParseSll(data []byte) (*layers.LinuxSLL, error) {
	nSize := SLLSize
	if len(data) < nSize {
		return nil, errors.New("not enough data to parse SLL")
	}
	sll := layers.LinuxSLL{}
	sll.PacketType = layers.LinuxSLLPacketType(binary.BigEndian.Uint16(data[0:2]))
	switch sll.PacketType {
	case layers.LinuxSLLPacketTypeHost, layers.LinuxSLLPacketTypeBroadcast, layers.LinuxSLLPacketTypeMulticast, layers.LinuxSLLPacketTypeOtherhost, layers.LinuxSLLPacketTypeOutgoing:
		break // expected packet types
	case layers.LinuxSLLPacketTypeLoopback, layers.LinuxSLLPacketTypeFastroute:
		break //log.Tracef("unexpected SLL packet type: %v", sll.PacketType) // known but not expected
	default:
		return nil, fmt.Errorf("packet type %x not supported", sll.PacketType) // shouldn't be, but...
	}
	sll.AddrType = binary.BigEndian.Uint16(data[2:4])
	switch sll.AddrType {
	case 1, 0x100, 0x300, 0x304:
		break // a skeleton fo further validations
	default:
		break //log.Tracef("ARPHRD_ type: %x", sll.AddrType)
	}
	sll.AddrLen = binary.BigEndian.Uint16(data[4:6])
	sll.Addr = data[6 : sll.AddrLen+6]
	sll.EthernetType = layers.EthernetType(binary.BigEndian.Uint16(data[14:16]))
	sll.BaseLayer = layers.BaseLayer{Contents: data[:SLLSize], Payload: data[SLLSize:]}
	switch sll.EthernetType {
	case layers.EthernetTypeIPv4, layers.EthernetTypeIPv6:
		return &sll, nil
	}
	return nil, ErrorNotAnIpPacket
}
