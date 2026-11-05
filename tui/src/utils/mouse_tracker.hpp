#ifndef MOUSE_TRACKER_HPP
#define MOUSE_TRACKER_HPP
#include <cstdint>
#include <vector>

namespace utils
{
    class MouseTracker
    {
    public:
        void addPoint(int x, int y);
        uint32_t getSeed() const;
        void clear();
        size_t size() const;

    private:
        std::vector< std::pair< int, int > > points_;
    };
}

#endif