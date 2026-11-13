#include "mouse_tracker.hpp"

namespace utils
{
    void MouseTracker::addPoint(int x, int y)
    {
        const std::pair< int, int > p{ x, y };

        if (!points_.empty() && points_.back() == p)
        {
            return;
        }

        if (points_.size() >= kMaxPoints)
        {
            points_.pop_front();
        }
        points_.push_back(p);
    }

    uint32_t MouseTracker::getSeed() const
    {
        if (points_.empty())
        {
            return 12345;
        }

        uint32_t h = 2166136261u;
        auto mix = [&h](uint32_t v)
        {
            for (int i = 0; i < 4; ++i)
            {
                h ^= (v >> (8 * i)) & 0xFFu;
                h *= 16777619u;
            }
        };

        for (const auto& p : points_)
        {
            mix(static_cast< uint32_t >(p.first));
            mix(static_cast< uint32_t >(p.second));
        }

        return h == 0 ? 1 : h;
    }

    void MouseTracker::clear()
    {
        points_.clear();
    }

    std::size_t MouseTracker::size() const
    {
        return points_.size();
    }
}
