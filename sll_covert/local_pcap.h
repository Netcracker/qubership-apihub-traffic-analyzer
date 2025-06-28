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

#ifndef LOCAL_PCAP_H
#define LOCAL_PCAP_H
#if defined(NO_PCAP_DETECTED) || defined(__CYGWIN__)
typedef struct {
    uint32_t magic;
    uint16_t version_major;
    uint16_t version_minor;
    uint32_t thiszone;	/* not used - SHOULD be filled with 0 */
    uint32_t sigfigs;	/* not used - SHOULD be filled with 0 */
    uint32_t snaplen;	/* max length saved portion of each pkt */
    uint32_t linktype;
}FileHeader;

typedef struct {
    uint32_t	tv_sec;		/* seconds */
    uint32_t	tv_usec;	/* and microseconds */
}PCTimeVal;

typedef struct  {
    PCTimeVal ts;	/* time stamp */
    uint32_t caplen;	/* length of portion present in data */
    uint32_t len;	    /* length of this packet prior to any slicing */
}PacketHeader;
#ifndef SLL_HDR_LEN
#define SLL_HDR_LEN	16		/* total header length */
#endif // SLL_HDR_LEN
#ifndef SLL_ADDRLEN
#define SLL_ADDRLEN	8		/* length of address field */
#endif // SLL_ADDRLEN

typedef struct {
    uint16_t sll_pkttype;		/* packet type */
    uint16_t sll_hatype;		/* link-layer address type */
    uint16_t sll_halen;		/* link-layer address length */
    uint8_t  sll_addr[SLL_ADDRLEN];	/* link-layer address */
    uint16_t sll_protocol;		/* protocol */
}  SLLHeader;
#ifndef ETH_ALEN
#define ETH_ALEN	6
#endif
#ifndef ETH_HLEN
#define ETH_HLEN	14
#endif

typedef struct {
    unsigned char	h_dest[ETH_ALEN];	/* destination eth addr	*/
    unsigned char	h_source[ETH_ALEN];	/* source ether addr	*/
    uint16_t		h_proto;		/* packet type ID field	*/
}  EthHeader;

#else
#include <pcap.h>
#include <pcap/sll.h> // for struct sll_header
#include <linux/if_ether.h> // for struct ethhdr
typedef struct pcap_file_header FileHeader;
typedef struct pcap_pkthdr      PacketHeader;
typedef struct sll_header       SLLHeader;
typedef struct ethhdr           EthHeader;
#endif
#ifndef FILE_HEADER_MAGIC
#define FILE_HEADER_MAGIC  0xA1B23C4D
#endif
#ifndef FILE_HEADER_VERSION_MAJOR
#define FILE_HEADER_VERSION_MAJOR  2
#endif
#ifndef FILE_HEADER_VERSION_MINOR
#define FILE_HEADER_VERSION_MINOR  4
#endif
#ifndef FILE_HEADER_LINK_TYPE
#define FILE_HEADER_LINK_TYPE  1
#endif
#ifndef DEF_SNAP_LEN
#define DEF_SNAP_LEN 262144
#endif
#ifndef LINUX_SLL_HOST
#define LINUX_SLL_HOST		0
#endif
#ifndef LINUX_SLL_BROADCAST
#define LINUX_SLL_BROADCAST	1
#endif
#ifndef LINUX_SLL_MULTICAST
#define LINUX_SLL_MULTICAST	2
#endif
#ifndef LINUX_SLL_OTHERHOST
#define LINUX_SLL_OTHERHOST	3
#endif
#ifndef LINUX_SLL_OUTGOING
#define LINUX_SLL_OUTGOING	4
#endif
#ifndef LINUX_SLL_LOOPBACK
#define LINUX_SLL_LOOPBACK	5
#endif
#ifndef LINUX_SLL_FASTROUTE
#define LINUX_SLL_FASTROUTE	6
#endif
#ifndef ETHERNET_TYPE_IPV4
#define ETHERNET_TYPE_IPV4 0x0800
#endif
#ifndef ETHERNET_TYPE_IPV6
#define ETHERNET_TYPE_IPV6 0x86DD
#endif
#endif //LOCAL_PCAP_H
