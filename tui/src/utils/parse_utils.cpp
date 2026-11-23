#include "parse_utils.hpp"
#include <string>
#include <stdexcept>

namespace 
{

    template <typename T, typename Parser>
    T strict_parse(const std::string& str, Parser parser)
    {
        if (str.empty())
        {
            throw std::invalid_argument("Пустое значение параметра!");
        }

        size_t pos = 0;
        T result;
        try
        {
            result = parser(str, &pos);
        }
        catch (const std::invalid_argument&)
        {
            throw std::invalid_argument(utils::PARSE_ERROR_MSG);
        }
        catch (const std::out_of_range&)
        {
            throw std::invalid_argument("Число вне допустимого диапазона!");
        }

        if (pos != str.size())
        {
            throw std::invalid_argument(utils::PARSE_ERROR_MSG);
        }

        return result;
    }

} 

int utils::strict_stoi(const std::string& str)
{
    return strict_parse<int>(str, [](const std::string& s, size_t* p) {
        return std::stoi(s, p);
    });
}

double utils::strict_stod(const std::string& str)
{
    return strict_parse<double>(str, [](const std::string& s, size_t* p) {
        return std::stod(s, p);
    });
}
