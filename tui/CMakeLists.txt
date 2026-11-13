cmake_minimum_required(VERSION 3.20)
project(tui_client CXX)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

set(CMAKE_EXPORT_COMPILE_COMMANDS ON)

include(FetchContent)

FetchContent_Declare(
    ftxui
    URL https://codeload.github.com/ArthurSonzogni/FTXUI/tar.gz/refs/tags/v6.0.0
    DOWNLOAD_EXTRACT_TIMESTAMP TRUE
)
FetchContent_MakeAvailable(ftxui)

FetchContent_Declare(
    json
    URL https://codeload.github.com/nlohmann/json/tar.gz/refs/tags/v3.11.3
    DOWNLOAD_EXTRACT_TIMESTAMP TRUE
)
FetchContent_MakeAvailable(json)

set(CPR_USE_SYSTEM_CURL ON CACHE BOOL "Use system curl" FORCE)

FetchContent_Declare(
    cpr
    URL https://codeload.github.com/libcpr/cpr/tar.gz/refs/tags/1.10.5
    DOWNLOAD_EXTRACT_TIMESTAMP TRUE
)
FetchContent_MakeAvailable(cpr)

file(GLOB_RECURSE SOURCES "src/*.cpp")

add_executable(tui_client ${SOURCES})
target_include_directories(tui_client PRIVATE src)

target_link_libraries(tui_client PRIVATE
    ftxui::screen
    ftxui::dom
    ftxui::component
    cpr::cpr
    nlohmann_json::nlohmann_json
)
