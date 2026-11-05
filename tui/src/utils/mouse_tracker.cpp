#include "mouse_tracker.hpp"
#include <functional>
#include <string>
#include <sstream>

namespace utils
{
    void MouseTracker::addPoint(int x, int y)
    {
        if (points_.size() >= 100)
        {
            points_.erase(points_.begin());
            points_.push_back({ x, y });
        }
    }

    uint32_t MouseTracker::getSeed() const
    {
        if (points_.empty())
        {
            return 12345;
        }

        std::stringstream ss;
        for (const auto& p : points_)
        {
            ss << p.first << "," << p.second << ";";
        }

        std::hash< std::string > hasher;
        size_t hash = hasher(ss.str());

        uint32_t seed = static_cast< uint32_t >(hash & 0xFFFFFFFF);

        return seed == 0 ? 1 : seed;
    }

    void MouseTracker::clear()
    {
        points_.clear();
    }

    size_t MouseTracker::size() const
    {
        return points_.size();
    }
}