#include <iostream>
#include <cstdlib>
#include <ctime>
#if defined(_MSC_VER) || defined(__BORLANDC__)
#else
#include <unistd.h>
#ifndef errno_t
#define errno_t int
#endif // errno_t
#endif
#include "appConfig.h"
#include "fileConverter.h"

int main(int argc, char* argv[]) {
    auto cfg = utils::appConfig::getInstance();
    int nRet = cfg.parseCmdLine(argc,argv);
    if(nRet!=EXIT_FAILURE) {
        if (auto cf = fileConverter(cfg); cf.valid()) {
            nRet = cf.convertFile();
            if (nRet != EXIT_SUCCESS) {
                std::cerr << "Convert failed:" << cf.getLastErrorText() << std::endl;
            }
        }
        else {
            std::cerr << "Unable to create converter:" << cf.getLastErrorText() << std::endl;
            nRet = EXIT_FAILURE;
        }
    }
    return nRet;
}
