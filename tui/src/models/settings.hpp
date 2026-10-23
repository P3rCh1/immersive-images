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
        j = nlohmann::json::object();

        if (s.style) j["style"] = geometryToStr(*s.style);
        if (s.palette) j["palette"] = paletteToStr(*s.palette);
        if (s.additional) j["additional"] = additionalToStr(*s.additional);
        if (s.width) j["width"] = *s.width;
        if (s.height) j["height"] = *s.height;
        if (s.scale) j["scale"] = *s.scale;
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

    void from_json(const nlohmann::json& j, GenerationSettingsResponse& s)
    {
        s.style = strToGeometry(j.at("style").get< std::string >());
        s.palette = strToPalette(j.at("palette").get< std::string >());
        s.width = j.at("width").get< int >();
        s.height = j.at("height").get< int >();
        s.scale = j.at("scale").get< double >();

        if (j.contains("additional") && !j["additional"].is_null())
        {
            s.additional = strToAdditional(j["additional"].get< std::string >());
        }
    }
}

#endif