#ifndef SETTINGS_HPP
#define SETTINGS_HPP
#include <optional>
#include <nlohmann/json.hpp>
#include "enums.hpp"

namespace models
{
    struct GenerationSettingsRequest
    {
        std::optional< Geometry > style;
        std::optional< Palette > palette;
        std::optional< Additional > additional;
        std::optional< int > width;
        std::optional< int > height;
        std::optional< double > scale;
    };

    void to_json(nlohmann::json& j, const GenerationSettingsRequest& s)
    {

    }

    struct GenerationSettingsResponse
    {
        Geometry style;
        Palette palette;
        std::optional< Additional > additional;
        int width;
        int height;
        double scale;
    };
}

#endif