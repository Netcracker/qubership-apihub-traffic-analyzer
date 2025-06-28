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
#include <iostream>
#include <sstream>
#include <cstdio>
#include <fcntl.h>
#include <cerrno>
#include <cstdlib>
#if defined(_MSC_VER) || defined(__BORLANDC__)
#else
#include <unistd.h>
#define errno_t int
#endif

#include "fileConverter.h"

#include <cstring>
#include <iomanip>

#include "local_pcap.h"
#include "little_big.h"
#if defined(_MSC_VER) || defined(__BORLANDC__)
#define READ_FILE_BINARY    "rb"
#define WRITE_FILE_BINARY    "wb"
#else
#define READ_FILE_BINARY    "r"
#define WRITE_FILE_BINARY    "w"
/**
 * retrives error text by error code
 * @param buf buffer to receive error text
 * @param bufLen length of th buffer
 * @param err error code
 * @return error or 0
 */
errno_t strerror_s(char* buf, size_t bufLen, errno_t err)
{
# ifdef __GNUC__
    if (strerror_r(err, buf, bufLen) == 0)
#else
    if (strerror_r(err, buf, bufLen) != nullptr)
#endif
        return 0;
    return EINVAL;
}

/**
 * opens file (suppress MSVC warnings)
 * @param ppF pointer to recive opened file structure (FILE*)
 * @param fileName file name
 * @param fileOpenMode file open mode (see fopen)
 * @return 0 if succeeded or error code
 */
errno_t fopen_s(FILE** ppF, const char* fileName, const char* fileOpenMode)
{
    *ppF = fopen(fileName, fileOpenMode);
    if (*ppF == nullptr) {
        return errno;
    }
    return 0;
}
#endif

/**
 * reads bytes from input file
 * @param ptr buffer to receive bytes from file
 * @param size requested byte count
 * @return bytes count that actually read
 */
size_t fileConverter::readInput(void *ptr, const size_t size) {
    if(feof(fdIn)) return 0;
    const auto res = fread(ptr, 1, size, fdIn);
    if (res == -1) {
        if (char text[1024]; strerror_s(text, 1024, errno)==0) {
            lastErrorText = "Error reading file";
        } else {
            lastErrorText = text;
        }
    }
    return res;
}

/**
 * writes bytes to the output file
 * @param ptr bytes to be written
 * @param size bytes count
 * @return bytes count that actually written
 */
size_t fileConverter::writeOutput(const void *ptr, const size_t size) {
    const auto res = fwrite(ptr, 1, size, fdOut);
    if (res == -1) {
        if (char text[1024]; strerror_s(text, 1024, errno) == 0) {
            lastErrorText = "Error writing file";
        } else {
            lastErrorText = text;
        }
    }
    return res;
}

/**
 * validates PCAP file header
 * @param fileHeader header that has been read
 * @return true if header was valid, otherwise - false
 */
bool fileConverter::validateFileHeader(const FileHeader &fileHeader) {
    return fileHeader.magic==FILE_HEADER_MAGIC &&
           fileHeader.version_major==FILE_HEADER_VERSION_MAJOR &&
           fileHeader.version_minor==FILE_HEADER_VERSION_MINOR &&
           fileHeader.thiszone ==0 && fileHeader.sigfigs == 0 &&
           fileHeader.linktype == FILE_HEADER_LINK_TYPE;
}

/**
 * resizes internal buffer
 * @param size requested size
 * @return true if resized, false - if not
 */
bool fileConverter::resizeBuffer(const size_t size) {
    if(size>bytesSize) {
        delete[] bytes;
        try {
            bytes = new unsigned char[size];
            if(bytes==nullptr) {
                bytesSize = 0;
                return false;
            }
        } catch (const std::bad_alloc&) {
            bytesSize = 0;
            return false;
        }
        bytesSize = size;
    }
    return true;
}

/**
 * print bytes in a hex dump manner
 * @param bytes bytes to print
 * @param size count to print
 */
void printBytes(const unsigned char *bytes, const size_t size) {
    unsigned char x[17];
    for(size_t i=0; i<size; i++) {
        constexpr unsigned char MinPrintable = 19;
        if (constexpr unsigned char MaxPrintable = 127; bytes[i]>MinPrintable && bytes[i]<MaxPrintable) {
            x[i%16] = bytes[i];
        } else {
            x[i%16] = '.';
        }

        std::cout << std::hex << std::setw(2) << std::setfill('0') << static_cast<uint16_t>(bytes[i]) << " ";
        if (i % 16 == 15) {
            x[16] = 0;
            std::cout << " | " << x << std::endl;
        }
    }
    if (size % 16 != 0) {
        x[size % 16] = 0;
        std::cout << " | " << x << std::endl;
    }
}

/**
 * validates Linux SLL packet header
 * @param header header values will be put here
 * @param ptr bytes that pissibly contains a header
 * @return true if header was valid, otherwise - false
 */
bool fileConverter::validateSLLHeader(SLLHeader& header, const unsigned char* ptr) {
    //auto pHead = reinterpret_cast<const SLLHeader *>(ptr);
    header.sll_pkttype  = MAKE_SHORT(ptr[1], ptr[0]);
    header.sll_hatype   = MAKE_SHORT(ptr[3], ptr[2]);
    header.sll_halen    = MAKE_SHORT(ptr[5], ptr[4]);
    memcpy(&header.sll_addr, ptr+6, SLL_ADDRLEN);
    header.sll_protocol = MAKE_SHORT(ptr[15], ptr[14]);
    switch (header.sll_pkttype) {
    case LINUX_SLL_HOST:
    case LINUX_SLL_BROADCAST:
    case LINUX_SLL_MULTICAST:
    case LINUX_SLL_OTHERHOST:
    case LINUX_SLL_OUTGOING:
        //break;// ok
    case LINUX_SLL_LOOPBACK:
    case LINUX_SLL_FASTROUTE:
        break;// strange, but ok
        default:
        //std::cerr << "Unsupported sll type:" << header.sll_pkttype << std::endl;
        return false;
    }
    switch (header.sll_hatype) {
    case 1:
    case 3:
    case 0x100:
    case 0x300:
    case 0x304:
        break;
    default:
        return false;
    }
    if(/*header.sll_halen==0 ||*/ header.sll_halen>SLL_ADDRLEN) {
        //std::cerr << "Invalid address len: " << std::hex << static_cast<uint16_t>(ptr[4]) << ":" << static_cast<uint16_t>(ptr[5]) << "=" << std::dec << header.sll_halen << std::endl;
        return false;
    }
    switch (header.sll_protocol) {
    case ETHERNET_TYPE_IPV4:
    case ETHERNET_TYPE_IPV6:
        break;
        default:
        //std::cerr << "Invalid protocol: " << std::hex << static_cast<uint16_t>(ptr[14]) << ":" << static_cast<uint16_t>(ptr[15]) << "=" << std::dec << header.sll_protocol << std::endl;
        return false;
    }
    return true;
}

/**
 * constructs file converter
 * @param cfg an application config
 */
fileConverter::fileConverter(const utils::appConfig& cfg) : fdIn(nullptr), fdOut(nullptr), bytes(nullptr), bytesSize(0)
{
    if(cfg.isOptSet(utils::appConfig::OPT_INPUT)) {
        if (fopen_s(&fdIn, cfg.getOptVal(utils::appConfig::OPT_INPUT).c_str(), READ_FILE_BINARY) != 0) {
            fdIn = nullptr;
        }
    } else {
        fdIn = stdin;
    }
    if(fdIn == nullptr) {
        lastErrorText = "Error opening input file";
    } else {
        if(cfg.isOptSet(utils::appConfig::OPT_OUTPUT)) {
            if (fopen_s(&fdOut, cfg.getOptVal(utils::appConfig::OPT_OUTPUT).c_str(), WRITE_FILE_BINARY) != 0) {
                fdOut = nullptr;
            }
        }
        if(fdOut == nullptr) {
            lastErrorText = "Error opening output file";
        }
    }
    bytes = new unsigned char[DEF_SNAP_LEN];
    bytesSize = DEF_SNAP_LEN;
}

/**
 * destructs file converter
 */
fileConverter::~fileConverter() {
    if (fdIn != nullptr) {
        fclose(fdIn);
        fdIn = nullptr;
    }
    if (fdOut != nullptr) {
        fclose(fdOut);
        fdOut = nullptr;
    }
    delete[] bytes;
    bytesSize = 0;
}

/**
 * reads input file, writes converted packets into output file
 * @return EXIT_SUCCESS|EXIT_FAILURE
 */
int fileConverter::convertFile() {
    FileHeader header;
    if(readInput(&header, sizeof(FileHeader)) != sizeof(FileHeader)) {
        lastErrorText = "Header was not read";
        return EXIT_FAILURE;
    }
    if(!validateFileHeader(header)) {
        lastErrorText = "Header was not validated";
        return EXIT_FAILURE;
    }
    if(writeOutput(&header, sizeof(FileHeader)) != sizeof(FileHeader)) {
        lastErrorText = "Header was not written";
        return EXIT_FAILURE;
    }
    PacketHeader packetHeader;
    auto capMaxLen = header.snaplen;
    if (!resizeBuffer(capMaxLen)) {
        lastErrorText = "Unable to resize input buffer to " + std::to_string(capMaxLen);
        return EXIT_FAILURE;
    }
    size_t idx = 0;
    while(readInput(&packetHeader, sizeof(PacketHeader)) == sizeof(PacketHeader)) {
        idx ++;
        std::cerr.flush();
        if(packetHeader.caplen>capMaxLen || packetHeader.len>capMaxLen) {
            std::stringstream oss;
            oss << "improper capture len:" << std::dec << packetHeader.caplen << " or length:" << packetHeader.len;
            lastErrorText = oss.str();
            return EXIT_FAILURE;
        }
        if (packetHeader.caplen == 0) {
            continue;
        }
        auto rl = readInput(bytes, packetHeader.caplen);
        if (packetHeader.caplen < 12) {
            continue;
        }
        if(rl!=packetHeader.caplen) {
            std::stringstream oss;
            oss << "unable to read packet body:" << rl << " instead of " << packetHeader.caplen << std::endl;
            lastErrorText = oss.str();
            return EXIT_FAILURE;
        }
        // make some dirty things
        auto toBeWritten = rl;
        // validate SLL header
        SLLHeader sll_header;
        if(validateSLLHeader(sll_header, bytes)) {
            packetHeader.caplen -= SLL_HDR_LEN - ETH_HLEN;
            packetHeader.len -= SLL_HDR_LEN - ETH_HLEN;
            if (writeOutput(&packetHeader, sizeof(PacketHeader)) != sizeof(PacketHeader)) {
                lastErrorText = "unable to write packet header";
                return EXIT_FAILURE;
            }
            auto pHead = &sll_header;
            //size_t addrLen = pHead->sll_halen;
            // if(addrLen >= ETH_ALEN)
            //     addrLen = ETH_ALEN - 2;
            EthHeader ethHeader;
            memcpy(ethHeader.h_source, pHead->sll_addr, ETH_ALEN-1);
            ethHeader.h_source[ETH_ALEN-1] = 1;
            memcpy(ethHeader.h_dest, pHead->sll_addr, ETH_ALEN-1);
            ethHeader.h_dest[ETH_ALEN-1] = 2;
            ethHeader.h_proto = BYTESWAP_USHORT(pHead->sll_protocol);
            //printBytes((const unsigned char*)&ethHeader, sizeof(ethHeader));
            if(writeOutput(&ethHeader, sizeof(EthHeader))!=sizeof(EthHeader)) {
                return EXIT_FAILURE;
            }
            toBeWritten -= SLL_HDR_LEN;
            if(writeOutput(bytes + SLL_HDR_LEN, toBeWritten)!=toBeWritten) {
                lastErrorText = "unable to write rest of the packet";
                return EXIT_FAILURE;
            }
            //return EXIT_SUCCESS;
            continue;
        }
        // write unmodified data
        if(writeOutput(&packetHeader, sizeof(PacketHeader))!=sizeof(PacketHeader)) {
            lastErrorText = "unable to write packet header";
            return EXIT_FAILURE;
        }
        if(writeOutput(bytes, toBeWritten)!=toBeWritten) {
            lastErrorText = "unable to write packet body";
            return EXIT_FAILURE;
        }
    }
    return EXIT_SUCCESS;
}
