#include "api_client.hpp"
#include "models/error.hpp"

namespace api
{
    ApiClient::ApiClient(const std::string& base_url):
        base_url_(base_url)
    {}


}