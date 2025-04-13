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
#ifndef _LITTLE_BIG_H_
#define _LITTLE_BIG_H_
#include <cstdint>
#if defined(__BYTE_ORDER) && __BYTE_ORDER == __BIG_ENDIAN || \
    defined(__BIG_ENDIAN__) || \
    defined(__ARMEB__) || \
    defined(__THUMBEB__) || \
    defined(__AARCH64EB__) || \
    defined(_MIBSEB) || defined(__MIBSEB) || defined(__MIBSEB__)
// It's a big-endian target architecture
#define __BIG_ENDIAN_DETECTED__
#elif defined(__BYTE_ORDER) && __BYTE_ORDER == __LITTLE_ENDIAN || \
    defined(__LITTLE_ENDIAN__) || \
    defined(__ARMEL__) || \
    defined(__THUMBEL__) || \
    defined(__AARCH64EL__) || \
    defined(_MIPSEL) || defined(__MIPSEL) || defined(__MIPSEL__)
// It's a little-endian target architecture
#ifndef __LITTLE_ENDIAN__LITTLE_BIG__
#define __LITTLE_ENDIAN__LITTLE_BIG__
#endif
#else
#pragma message ("Microsoft compiler makes your life more interesting. Usually little endian, but you can patch it there if not")
#ifndef __LITTLE_ENDIAN__LITTLE_BIG__
#define __LITTLE_ENDIAN__LITTLE_BIG__
#endif
#endif
// byte swap will be required
#ifdef __BIG_ENDIAN_DETECTED__
#define BYTESWAP_USHORT
#define BYTESWAP_UINT32
#define BYTESWAP_UINT64
// don't trust Microsoft and provide word builder
#ifndef MAKE_SHORT
#define MAKE_SHORT(_HIGH_,_LOW_)	(((_HIGH_) << 8) | (_LOW_)) & UINT16_MAX
#endif // !MAKE_SHORT
#endif
// no byte swap, but MAKE_SHORT macro
#ifdef __LITTLE_ENDIAN__LITTLE_BIG__
#if defined(_MSC_VER)
#include <intrin.h>
#define BYTESWAP_USHORT _byteswap_ushort
#define BYTESWAP_UINT32 _byteswap_ulong
#define BYTESWAP_UINT64 _byteswap_uint64
#else
#define BYTESWAP_USHORT(__VAL__)  (((__VAL__ >> 8) & 0x00FF) | ((__VAL__ << 8) & 0xFF00))
#define BYTESWAP_UINT32 __builtin_bswap32
#define BYTESWAP_UINT64 __builtin_bswap64
#endif
#ifndef MAKEWORD
#define MAKE_SHORT(_LOW_,_HIGH_)	((((_HIGH_) << 8) | (_LOW_)) & UINT16_MAX)
#else
#define MAKE_SHORT(_LOW_,_HIGH_)	MAKEWORD(_LOW_,_HIGH_)
#endif // !MAKEWORD
#endif
#endif  // _LITTLE_BIG_H_