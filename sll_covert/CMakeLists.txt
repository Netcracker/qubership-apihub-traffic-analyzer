cmake_minimum_required(VERSION 3.18)
project(sll_covert)

set(CMAKE_CXX_STANDARD 20)

if (WIN32)
    set(CMAKE_EXE_LINKER_FLAGS "-static-libgcc -static-libstdc++")
    set(PRJ_ROOT, "C:/Projects/unsorted/C_C++")
    include_directories(${PRJ_ROOT}/shared_src/c)
    add_executable(sll_convert sll_convert.cpp appConfig.h appConfig.cpp getopt_win32.h getopt_win32.c
            fileConverter.cpp
            fileConverter.h)
else(WIN32)
    add_executable(sll_convert sll_convert.cpp appConfig.h appConfig.cpp
            little_big.h
            fileConverter.cpp
            fileConverter.h
            local_pcap.h)
    target_link_libraries(sll_convert rt)
endif(WIN32)

find_package(PCAP)
if (PCAP_FOUND)
    message(STATUS "Using detected libpcap-dev package")
    target_link_libraries(sll_convert  ${PCAP_LIBRARY})
else(PCAP_FOUND)
    message(STATUS "Unable to detect libpcap-dev package - falling back to pkg-config...")
    find_package(PkgConfig REQUIRED)
    pkg_check_modules(PCAP libpcap-dev)
    if (PCAP)
        target_link_libraries(sll_convert  ${PCAP_LIBRARIES})
        target_include_directories(sll_convert PUBLIC ${PCAP_INCLUDE_DIRS})
    else(PCAP)
        message(STATUS "Unable to detect libpcap-dev package - notify compiler...")
        target_compile_definitions(sll_convert PUBLIC NO_PCAP_DETECTED=1)
    endif (PCAP)
endif(PCAP_FOUND)
