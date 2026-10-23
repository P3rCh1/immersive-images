#ifndef MODELS_HPP
#define MODELS_HPP
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
}


#endif