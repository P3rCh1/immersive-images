cmake_minimum_required(VERSION 3.20)
project(tui_client CXX)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

# Генерируем compile_commands.json для корректной подсветки в VS Code
set(CMAKE_EXPORT_COMPILE_COMMANDS ON)

# ============================================================
# Подключение зависимостей через FetchContent
# ============================================================
include(FetchContent)

# --- 1. FTXUI (Интерфейс) ---
FetchContent_Declare(
    ftxui
    GIT_REPOSITORY https://github.com/ArthurSonzogni/FTXUI.git
    GIT_TAG v5.0.0
)
FetchContent_MakeAvailable(ftxui)

# --- 2. nlohmann/json (Парсинг JSON) ---
FetchContent_Declare(
    json
    GIT_REPOSITORY https://github.com/nlohmann/json.git
    GIT_TAG v3.11.3
)
FetchContent_MakeAvailable(json)

# --- 3. cpr (HTTP-запросы) ---
# Важно: используем системный libcurl, чтобы избежать ошибок скачивания
set(CPR_FORCE_USE_SYSTEM_CURL ON CACHE BOOL "Use system curl" FORCE)

FetchContent_Declare(
    cpr
    GIT_REPOSITORY https://github.com/libcpr/cpr.git
    GIT_TAG 1.10.5
)
FetchContent_MakeAvailable(cpr)

# ============================================================
# Сборка исполняемого файла
# ============================================================

# Автоматически находим ВСЕ .cpp файлы в папке src и её подпапках
file(GLOB_RECURSE SOURCES "src/*.cpp")

# Создаем исполняемый файл
add_executable(tui_client ${SOURCES})

# Настраиваем пути для инклюдов (чтобы не писать ../)
target_include_directories(tui_client PRIVATE src)

# Подключаем библиотеки к нашему проекту
target_link_libraries(tui_client PRIVATE
    ftxui::screen
    ftxui::dom
    ftxui::component
    cpr::cpr
    nlohmann_json::nlohmann_json
)

# ============================================================
# Тесты (на будущее)
# ============================================================
# enable_testing()