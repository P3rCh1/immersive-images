cmake_minimum_required(VERSION 3.20)
project(tui_client CXX)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

set(CMAKE_EXPORT_COMPILE_COMMANDS ON)

include(FetchContent)
target_include_directories(tui_client PRIVATE src)

FetchContent_Declare(
    ftxui
    GIT_REPOSITORY https://github.com/ArthurSonzogni/FTXUI.git
    GIT_TAG v6.0.0
)
FetchContent_MakeAvailable(ftxui)

FetchContent_Declare(
    json
    GIT_REPOSITORY https://github.com/nlohmann/json.git
    GIT_TAG v3.11.3
)
FetchContent_MakeAvailable(json)

set(CPR_FORCE_USE_SYSTEM_CURL ON CACHE BOOL "Use system curl" FORCE)

FetchContent_Declare(
    cpr
    GIT_REPOSITORY https://github.com/libcpr/cpr.git
    GIT_TAG 1.10.5
)
FetchContent_MakeAvailable(cpr)

file(GLOB_RECURSE SOURCES "src/*.cpp")

add_executable(tui_client ${SOURCES})

target_link_libraries(tui_client PRIVATE
    ftxui::screen
    ftxui::dom
    ftxui::component
    cpr::cpr
    nlohmann_json::nlohmann_json
)