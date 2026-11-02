#ifndef ERROR_HPP
#define ERROR_HPP
#include <string>
#include <nlohmann/json.hpp>

namespace models
{
    struct ApiError
    {
        std::string message;
    };

    inline void from_json(const nlohmann::json& j, ApiError& e)
    {
        e.message = j.at("error").at("message").get< std::string >();
    }
}

#endif