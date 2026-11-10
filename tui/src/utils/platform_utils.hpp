#ifndef PLATFORM_UTILS_HPP
#define PLATFORM_UTILS_HPP

#include <string>
#include <filesystem>

namespace utils
{
    std::string sanitizeFileName(const std::string& s);
    void openWithSystemViewer(const std::filesystem::path& file);
}

#endif 