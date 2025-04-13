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

#ifndef FILE_CONVERTER_H
#define FILE_CONVERTER_H
#include <string>

#include "local_pcap.h"
#include "appConfig.h"


class fileConverter {
    FILE* fdIn;
    FILE* fdOut;
    std::string lastErrorText;
    unsigned char* bytes;
    size_t bytesSize;
    size_t readInput(void* ptr, size_t size);
    size_t writeOutput(const void* ptr, size_t size);

    static bool validateFileHeader(const FileHeader& fileHeader);
    bool resizeBuffer(size_t size);
    static bool validateSLLHeader(SLLHeader& header, const unsigned char* ptr);
public:
    explicit fileConverter(const utils::appConfig& cfg);
    ~fileConverter();
    int convertFile();

    bool valid() const {return fdIn!=nullptr && fdOut!=nullptr && bytes!=nullptr && bytesSize>0;};
    const std::string& getLastErrorText() const {return lastErrorText;}
};



#endif //FILE_CONVERTER_H
