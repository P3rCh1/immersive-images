#ifndef IMAGE_HPP
#define IMAGE_HPP
#include <string>
#include <vector>
#include <optional>
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

    struct ListImagesResponse
    {
        std::vector< ImageMeta > items;
        int total;
        int limit;
        std::optional< std::string > next;
    };

    struct CreateImagesRequest
    {
        int seed; 
        std::string name;
        std::optional< GenerationSettingsRequest > settings;
    };

    struct UpdateImageRequest
    {
        std::string name;
    };
}

#endif