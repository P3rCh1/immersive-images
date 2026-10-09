cmake_minimum_required(VERSION 3.20)
project(tui_client CXX)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

include(FetchContent)

FetchContent_Declare(
    ftxui
    GIT_REPOSITORY https://github.com/ArthurSonzogni/FTXUI.git
    GIT_TAG v5.0.0
)
FetchContent_MakeAvailable(ftxui)

add_executable(tui_client
    src/main.cpp
)

target_link_libraries(tui_client PRIVATE
    ftxui::screen
    ftxui::dom
    ftxui::component
)

enable_testing()