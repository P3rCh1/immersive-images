#ifndef MOUSE_TRACKER_HPP
#define MOUSE_TRACKER_HPP

#include <cstddef>
#include <cstdint>
#include <deque>
#include <utility>

namespace utils
{
    constexpr std::size_t kMinMousePoints = 20;

    class MouseTracker
    {
    public:
        static constexpr std::size_t kMaxPoints = 500;

        void addPoint(int x, int y);
        uint32_t getSeed() const;
        void clear();
        std::size_t size() const;

    private:
        std::deque< std::pair< int, int > > points_;
    };
}

#endif
