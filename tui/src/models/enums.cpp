#include "enums.hpp"
#include <stdexcept>

namespace models {

    Geometry stringToGeometry(const std::string& str) 
    {
        if (str == "PATCHES") return Geometry::PATCHES;
        if (str == "SPOTS") return Geometry::SPOTS;
        if (str == "RINGS") return Geometry::RINGS;
        if (str == "TERRACES") return Geometry::TERRACES;
        if (str == "WAVES") return Geometry::WAVES;
        if (str == "MOTTLED") return Geometry::MOTTLED;

        throw std::invalid_argument("Unknown Geometry: " + str);
    }

    Palette stringToPalette(const std::string& str)
    {
        if (str == "NEON") return Palette::NEON;
        if (str == "DARK") return Palette::DARK;
        if (str == "TOXIC") return Palette::TOXIC;
        if (str == "RAINBOW") return Palette::RAINBOW;
        if (str == "OCEAN") return Palette::OCEAN;
        if (str == "LAVA") return Palette::LAVA;
        if (str == "MONO") return Palette::MONO;

        throw std::invalid_argument("Unknown Palette: " + str);
    }

    Additional stringToAdditional(const std::string& str) 
    {
        if (str == "PIXEL") return Additional::PIXEL;
        if (str == "BLUR") return Additional::BLUR;

        throw std::invalid_argument("Unknown Additional: " + str);
    }

    std::string geometryToString(Geometry g) 
    {
        switch (g) 
        {
        case Geometry::PATCHES: return "PATCHES";
        case Geometry::SPOTS: return "SPOTS";
        case Geometry::RINGS: return "RINGS";
        case Geometry::TERRACES: return "TERRACES";
        case Geometry::WAVES: return "WAVES";
        case Geometry::MOTTLED: return "MOTTLED";
        }

        return "";
    }

    std::string paletteToString(Palette p) 
    {
        switch (p) 
        {
        case Palette::NEON: return "NEON";
        case Palette::DARK: return "DARK";
        case Palette::TOXIC: return "TOXIC";
        case Palette::RAINBOW: return "RAINBOW";
        case Palette::OCEAN: return "OCEAN";
        case Palette::LAVA: return "LAVA";
        case Palette::MONO: return "MONO";
        }

        return "";
    }

    std::string additionalToString(Additional a) 
    {
        switch (a) 
        {
        case Additional::PIXEL: return "PIXEL";
        case Additional::BLUR: return "BLUR";
        }
        
        return "";
    }
}