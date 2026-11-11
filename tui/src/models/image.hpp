#ifndef IMAGE_HPP
#define IMAGE_HPP
#include <string>
#include <vector>
#include <optional>
#include <nlohmann/json.hpp>
#include "settings.hpp"

namespace models
{
    struct ImageMeta
    {
        std::string uuid;
        std::string name;
        std::string created_at;
        int seed;
        GenerationSettingsResponse settings;
    };

    inline void from_json(const nlohmann::json& j, ImageMeta& m)
    {
        m.uuid = j.at("uuid").get< std::string >();
        m.name = j.at("name").get< std::string >();
        m.created_at = j.at("created_at").get< std::string >();
        m.seed = j.at("seed").get< int >();
        m.settings = j.at("settings").get<GenerationSettingsResponse>();
    }

    struct ListImagesResponse
    {
        std::vector< ImageMeta > items;
        int total;
        int limit;
        std::optional< std::string > next;
    };

    inline void from_json(const nlohmann::json& j, ListImagesResponse& res)
    {
        res.items = j.at("items").get< std::vector< ImageMeta > >();
        res.total = j.at("total").get< int >();
        res.limit = j.at("limit").get< int >();

        if (j.contains("next") && !j["next"].is_null()) 
        {
            res.next = j["next"].get< std::string >();
        }
    }

    struct CreateImageRequest
    {
        int seed; 
        std::string name;
        std::optional< GenerationSettingsRequest > settings;
    };

    inline void to_json(nlohmann::json& j, const CreateImageRequest& r)
    {
        j = nlohmann::json {
            { "seed", r.seed },
            { "name", r.name }
        };

        if (r.settings)
        {
            j["settings"] = *r.settings;
        }
    }

    struct UpdateImageRequest
    {
        std::string name;
    };

    inline void to_json(nlohmann::json& j, const UpdateImageRequest& r)
    {
        j = nlohmann::json({ "name", r.name });
    }
}

#endif
