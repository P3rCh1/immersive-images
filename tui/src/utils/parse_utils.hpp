#ifndef PARSE_UTILS_HPP
#define PARSE_UTILS_HPP
#include <string>

namespace utils
{
    constexpr const char* PARSE_ERROR_MSG = "Введено нечисловое значение параметра!";

    int strict_stoi(const std::string&);
    double strict_stod(const std::string&);
}

#endif
