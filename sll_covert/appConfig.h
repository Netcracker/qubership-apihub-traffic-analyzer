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

#ifndef APPCONFIG_H
#define APPCONFIG_H
#include <memory>
#include <string>
#include <map>

namespace utils {

class appConfig {
    static std::unique_ptr<appConfig> pThis;
    std::string progName;
    std::string optString;
    std::map<std::string, std::string> config;
    appConfig();
public:
    static constexpr const char* OPT_HELP = "help";
    static constexpr const char* OPT_INPUT = "input";
    static constexpr const char* OPT_OUTPUT = "output";
    static constexpr int NEED_HELP = -1;
    static appConfig& getInstance();
    int parseCmdLine(int argc, char** argv);
    static std::string getOptName(int opt);
    static int getShortOpt(const std::string& opt);
    void setOpt(const std::string& optName, const std::string& optVal);
    bool isOptSet(const std::string& optName) const {return config.contains(optName);};
    std::string getOptVal(const std::string& optName) const;
    void printUsage() const;
};

} // utils

#endif //APPCONFIG_H
