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
#include <cstring>
#include <ctime>
#include <cstdlib>
#if defined(_MSC_VER) || defined(__BORLANDC__)
#include "getopt_win32.h"
#define PATH_SEP '\\'
#else
#include <getopt.h>
#include <unistd.h>
#ifndef errno_t
#define errno_t int
#endif // errno_t
#define PATH_SEP '/'
#endif
#include "appConfig.h"

using namespace utils;
std::unique_ptr<appConfig> appConfig::pThis(new appConfig());
static constexpr option modOpts[] = {
    { appConfig::OPT_HELP,   no_argument,       nullptr, 'h' },
    { appConfig::OPT_INPUT,  required_argument, nullptr, 'i' },
    { appConfig::OPT_OUTPUT, required_argument, nullptr, 'o' },
    { nullptr,  0, nullptr, 0 }
    };

/**
 * private ctor
 */
appConfig::appConfig() {
    const auto* pOpt = reinterpret_cast<const option *>(&modOpts);
    while (pOpt->name != nullptr && pOpt->val != 0)
    {
        char c[3] = {};
        c[0] = static_cast<char>(pOpt->val);
        if (pOpt->has_arg == required_argument)
        {
            c[1] = ':';
            c[2] = '\0';
        }
        else
            c[1] = '\0';
        optString += c;
        pOpt++;
    }
}

/**
 * reference access method
 * @return instance reference
 */
appConfig & appConfig::getInstance() {
    return *pThis;
}

/**
 * parses command line arguments
 * @param argc argument count
 * @param argv argument values
 * @return EXIT_SUCCESS if arguments parsed, otherwise - EXIT_FAILURE
 */
int appConfig::parseCmdLine(const int argc, char **argv) {
    if(const char* ptr = strrchr(argv[0],PATH_SEP); ptr!=nullptr)    {
        progName.assign(ptr + 1);
    }
    else {
        progName.assign(argv[0]);
    }
    int longIndex;
    int needHelp = 0;
    const auto* pModOpts = reinterpret_cast<const struct option *>(&modOpts);
    int opt = getopt_long(argc, argv, optString.c_str(), pModOpts, &longIndex);
    while (opt!=-1 && needHelp==0)
    {
        std::string longOption(getOptName(opt));
        switch (opt)
        {
            case 'h':
                needHelp = NEED_HELP;
                break;
            case 'i':
            case 'o':
                if (optarg)
                {
                    if (!longOption.empty())
                    {
                        config[longOption]= optarg;
                    }
                }
                break;
            default:
                needHelp = opt;
                break;
        }
        //idx ++;
        opt = getopt_long(argc, argv, optString.c_str(), pModOpts, &longIndex);
    }
    if (needHelp == 0)
    {
        return EXIT_SUCCESS;
    }
    if (needHelp != NEED_HELP)
    {
        // unknown option
        std::cerr << "unknown option '" << static_cast<char>(needHelp) << "'" << std::endl;
    }
    printUsage();
    return EXIT_FAILURE;
}

/**
 * returns long option by its short equivalent
 * @param opt short option
 * @return long option or empty string
 */
std::string appConfig::getOptName(const int opt) {
    const auto* pModOpts = reinterpret_cast<const struct option *>(&modOpts);
    while (pModOpts->name != nullptr) {
        if (opt == pModOpts->val) {
            return pModOpts->name;
        }
        pModOpts ++;
    }
    return {};
}

/**
 * returns shot option by its long equivalent
 * @param opt long option
 * @return short option or -1 if no option found
 */
int appConfig::getShortOpt(const std::string& opt) {
    const auto* pModOpts = reinterpret_cast<const struct option *>(&modOpts);
    while (pModOpts->name != nullptr) {
        if (opt==pModOpts->name) {
            return pModOpts->val;
        }
        pModOpts ++;
    }
    return -1;
}

/**
 * set option value
 * @param optName long option
 * @param optVal option value
 */
void appConfig::setOpt(const std::string &optName, const std::string &optVal) {
    config[optName] = optVal;
}

static std::map<std::string, std::string> helpMap{
            { appConfig::OPT_HELP,   "To see this message" },
            { appConfig::OPT_INPUT,  "Input file" },
            { appConfig::OPT_OUTPUT,  "Output file" },
    };

/**
 * returns the option value
 * @param optName long option
 * @return option value or empty string
 */
std::string appConfig::getOptVal(const std::string &optName) const {
    try {
        return config.at(optName);
    } catch (...) {
    }
    return {};
}

/**
 * prints program usage message
 */
void appConfig::printUsage() const {
    std::cerr << "Usage: " << progName << " [options]" << std::endl;
    const auto* pModOpts = reinterpret_cast<const struct option *>(&modOpts);
    std::cerr << std::endl;
    std::cerr << "Options:" << std::endl;
    while (pModOpts->name != nullptr) {
        if (helpMap.contains(pModOpts->name)) {
            const int chr = pModOpts->val;
            std::cerr << "    -" << static_cast<char>(chr) << ", --" << pModOpts->name;
            if (pModOpts->has_arg & required_argument) {
                std::cerr << " <" << pModOpts->name << ">";
            }
            if (pModOpts->has_arg & optional_argument) {
                std::cerr << " [" << pModOpts->name << "]";
            }
            std::cerr << "    " << helpMap[pModOpts->name] << std::endl;
        }
        pModOpts++;
    }
}
