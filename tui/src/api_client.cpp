#include "api_client.hpp"
#include <cpr/cpr.h>
#include <nlohmann/json.hpp>
#include "models/error.hpp"

namespace
{
    using json = nlohmann::json;
}

namespace api
{
    ApiClient::ApiClient(const std::string& base_url):
        base_url_(base_url)
    {}

    std::string ApiClient::extractError(const std::string& response_body, int status_code)
    {
        try 
        {
            auto j = json::parse(response_body);
            models::ApiError err = j.get< models::ApiError >();

            return err.message;
        }
        catch (...)
        {
            return "HTTP error: " + std::to_string(status_code);
        }
    }

    bool ApiClient::checkHealth()
    {
        try
        {
            auto r = cpr::Get(cpr::Url{ base_url_ + "/healthz" }, cpr::Timeout{ 3000 });
        
            return r.status_code == 200;
        }
        catch (...)
        {
            return false;
        }
    }

    models::ImageMeta ApiClient::createImage(const models::CreateImageRequest& req)
    {
        std::string json_body = json(req).dump();

        auto r = cpr::Post(
            cpr::Url{ base_url_ + "/images" },
            cpr::Header{ { "Content-Type", "application/json" } },
            cpr::Body{ json_body },
            cpr::Timeout{ 60000 }
        );

        if (r.status_code != 201)
        {
            throw std::runtime_error(extractError(r.text, r.status_code));
        }

        return json::parse(r.text).get< models::ImageMeta >();
    }

    models::ListImagesResponse ApiClient::listImages(int limit, const std::string& start_uuid)
    {
        cpr::Parameters params{ { "limit", std::to_string(limit) } };
        if (!start_uuid.empty())
        {
            params.Add({ "start", start_uuid });
        }

        auto r = cpr::Get(
            cpr::Url{ base_url_ + "/images" }, 
            params, 
            cpr::Timeout{ 5000 }
        );
        if (r.status_code != 200)
        {
            throw std::runtime_error(extractError(r.text, r.status_code));
        }

        return json::parse(r.text).get< models::ListImagesResponse >();
    }

    models::ImageMeta ApiClient::getImage(const std::string& uuid)
    {
        auto r = cpr::Get(
            cpr::Url({ base_url_ + "/images/" + uuid }),
            cpr::Timeout{ 5000 }
        );

        if (r.status_code != 200)
        {
            throw std::runtime_error(extractError(r.text, r.status_code));
        }

        return json::parse(r.text).get< models::ImageMeta >();
    }

    models::ImageMeta ApiClient::updateImage(const std::string& uuid, const models::UpdateImageRequest& req)
    {
        std::string json_body = json(req).dump();

        auto r = cpr::Put(
            cpr::Url{ base_url_ + "/images/" + uuid },
            cpr::Header{ { "Content-Type", "application/json" } },
            cpr::Body{ json_body },
            cpr::Timeout{ 5000 }
        );

        if (r.status_code != 200)
        {
            throw std::runtime_error(extractError(r.text, r.status_code));
        }

        return json::parse(r.text).get< models::ImageMeta >();
    }

    void ApiClient::deleteImage(const std::string& uuid)
    {
        auto r = cpr::Delete(
            cpr::Url{ base_url_ + "/images/" + uuid },
            cpr::Timeout{ 5000 }
        );

        if (r.status_code != 204)
        {
            throw std::runtime_error(extractError(r.text, r.status_code));
        }
    }

    std::string ApiClient::downloadImage(const std::string& uuid)
    {
        auto r = cpr::Get(
            cpr::Url{ base_url_ + "/images/" + uuid + "/content" },
            cpr::Timeout{ 15000 }
        );

        if (r.status_code != 200)
        {
            throw std::runtime_error(extractError(r.text, r.status_code));
        }

        return r.text;
    }
}
