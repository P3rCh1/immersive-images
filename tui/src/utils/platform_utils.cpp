#include "platform_utils.hpp"

#include <algorithm>
#include <cctype>
#include <cstdlib>
#include <fstream>

namespace utils
{
    namespace
    {
#if !defined(_WIN32)
        std::string shellQuote(const std::string& s)
        {
            std::string r = "'";
            for (char c : s)
            {
                if (c == '\'') r += "'\\''";
                else r += c;
            }
            r += "'";
            return r;
        }

        bool isWsl()
        {
#if defined(__linux__)
            std::ifstream f("/proc/version");
            std::string line;
            std::getline(f, line);
            std::transform(line.begin(), line.end(), line.begin(),
                           [](unsigned char c) { return static_cast<char>(std::tolower(c)); });
            return line.find("microsoft") != std::string::npos ||
                   line.find("wsl") != std::string::npos;
#else
            return false;
#endif
        }
#endif
    } 

    std::string sanitizeFileName(const std::string& s)
    {
        std::string r;
        for (unsigned char c : s)
        {
            if (std::isalnum(c) || c == '-' || c == '_')
            {
                r.push_back(static_cast<char>(c));
            }
        }
        return r.empty() ? "image" : r;
    }

    void openWithSystemViewer(const std::filesystem::path& file)
    {
#if defined(_WIN32)
        std::filesystem::path p = file;
        p.make_preferred();
        std::string cmd = "start \"\" \"" + p.string() + "\"";
        std::system(cmd.c_str());
#elif defined(__APPLE__)
        std::string cmd = "open " + shellQuote(file.string()) + " </dev/null >/dev/null 2>&1 &";
        std::system(cmd.c_str());
#else
        const std::string q = shellQuote(file.string());
        std::string cmd;
        if (isWsl())
        {
            cmd = "( if command -v explorer.exe >/dev/null 2>&1; then "
                  "explorer.exe \"$(wslpath -w " + q + ")\"; "
                  "else xdg-open " + q + "; fi ) </dev/null >/dev/null 2>&1 &";
        }
        else
        {
            cmd = "( xdg-open " + q + " || gio open " + q + " ) </dev/null >/dev/null 2>&1 &";
        }
        std::system(cmd.c_str());
#endif
    }
} 