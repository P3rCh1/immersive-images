#ifndef ENUMS_HPP
#define ENUMS_HPP
#include <string>

namespace models
{
    enum class Geometry
    {
        PATCHES,
        SPOTS,
        RINGS,
        TERRACES,
        WAVES,
        MOTTLED
    };

    enum class Palette
    {
        NEON,
        DARK,
        TOXIC,
        RAINBOW,
        OCEAN,
        LAVA,
        MONO
    };

    enum class Additional
    {
        PIXEL,
        BLUR
    };

    Geometry strToGeometry(const std::string&);
    Palette strToPalette(const std::string&);
    Additional strToAdditional(const std::string&);

    std::string geometryToStr(const Geometry&);
    std::string paletteToStr(const Palette&);
    std::string additionalToStr(const Additional&);
}

#endif