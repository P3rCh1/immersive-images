#include "app.hpp"
#include <ftxui/component/event.hpp>
#include <ftxui/component/mouse.hpp>
#include <ftxui/dom/elements.hpp>

#include <algorithm>
#include <cctype>
#include <chrono>
#include <cstdlib>
#include <filesystem>
#include <fstream>
#include <functional>
#include <stdexcept>

namespace
{
    constexpr std::size_t kMinMousePoints = 20;

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

namespace ui
{
    using namespace ftxui;

    App::App(api::ApiClient& client): client_(client)
    {
        loadImages();

        MenuOption opt = MenuOption::Vertical();
        opt.on_enter = [this] { openSelected(); };
        menu_ = Menu(&display_items_, &selected_index_, opt);
    }

    void App::loadImages()
    {
        try
        {
            auto response = client_.listImages(50);
            images_ = response.items;
            display_items_.clear();

            for (const auto& img : images_)
            {
                display_items_.push_back(img.name + " [ " + img.created_at + " ]");
            }

            if (display_items_.empty())
            {
                display_items_.push_back("Список пуст. Нажмите N для создания.");
            }
        }
        catch (const std::exception& e)
        {
            images_.clear();
            display_items_.clear();
            display_items_.push_back("Ошибка загрузки: " + std::string(e.what()));
        }

        const int count = static_cast<int>(display_items_.size());
        selected_index_ = std::clamp(selected_index_, 0, count - 1);
    }

    void App::deleteImage()
    {
        try
        {
            client_.deleteImage(selected_uuid_);
            loadImages();
            setScreen(Screen::List);
        }
        catch (const std::exception& e)
        {
            error_text_ = e.what();
            setScreen(Screen::Error);
        }
    }

    void App::createImage()
    {
        if (new_name_.empty())
        {
            error_text_ = "Имя не может быть пустым!";
            setScreen(Screen::Error);
            return;
        }

        try
        {
            uint32_t seed = mouse_.getSeed();
            if (mouse_.size() < kMinMousePoints)
            {
                seed ^= static_cast<uint32_t>(std::hash<std::string>{}(new_name_));
                seed ^= static_cast<uint32_t>(
                    std::chrono::steady_clock::now().time_since_epoch().count());
                if (seed == 0) seed = 1;
            }

            models::CreateImageRequest req;
            req.seed = seed;
            req.name = new_name_;

            client_.createImage(req);

            new_name_.clear();
            mouse_.clear();
            loadImages();
            setScreen(Screen::List);
        }
        catch (const std::exception& e)
        {
            error_text_ = e.what();
            setScreen(Screen::Error);
        }
    }

    void App::openImage(const std::string& uuid)
    {
        try
        {
            std::string image_data = client_.downloadImage(uuid);

            std::filesystem::path path =
                std::filesystem::temp_directory_path() /
                ("image_" + sanitizeFileName(uuid) + ".png");

            {
                std::ofstream file(path, std::ios::binary);
                if (!file.is_open()) throw std::runtime_error("Не удалось создать файл");
                file.write(image_data.data(), static_cast<std::streamsize>(image_data.size()));
            }

            openWithSystemViewer(path);
            status_text_ = "Файл: " + path.string();
        }
        catch (const std::exception& e)
        {
            error_text_ = "Ошибка скачивания: " + std::string(e.what());
            setScreen(Screen::Error);
        }
    }

    void App::openSelected()
    {
        if (images_.empty()) return;
        if (selected_index_ < 0 || selected_index_ >= static_cast<int>(images_.size())) return;

        selected_uuid_ = images_[selected_index_].uuid;
        status_text_.clear();
        setScreen(Screen::Details);
    }

    void App::openCreate()
    {
        new_name_.clear();
        mouse_.clear();
        setScreen(Screen::Create);
        if (input_name_) input_name_->TakeFocus();
    }

    void App::initComponents(ScreenInteractive& screen)
    {
        auto btn_refresh = Button(" Обновить (R) ", [this] { loadImages(); });
        auto btn_open_d  = Button(" Детали (Enter) ", [this] { openSelected(); });
        auto btn_create  = Button(" Создать (N) ", [this] { openCreate(); });
        auto btn_quit    = Button(" Выход (Q) ", [&screen] { screen.Exit(); });

        auto list_buttons = Container::Horizontal({ btn_refresh, btn_open_d, btn_create, btn_quit });
        list_container_ = Container::Vertical({ menu_, list_buttons });

        auto btn_back   = Button(" Назад (ESC) ", [this] { setScreen(Screen::List); });
        auto btn_delete = Button(" Удалить (DEL) ", [this] { deleteImage(); });
        auto btn_open   = Button(" Открыть (O) ", [this] { openImage(selected_uuid_); });

        auto details_buttons = Container::Horizontal({ btn_back, btn_open, btn_delete });
        details_container_ = Container::Vertical({ details_buttons });

        input_name_ = Input(&new_name_, "Введите название (сид создастся от движения мыши)...");

        btn_generate_ = Button(" Сгенерировать ", [this] { createImage(); });
        auto btn_cancel = Button(" Отмена (ESC) ", [this] {
            new_name_.clear();
            mouse_.clear();
            setScreen(Screen::List);
        });

        auto create_buttons = Container::Horizontal({ btn_generate_, btn_cancel });
        create_container_ = Container::Vertical({ input_name_, create_buttons });

        auto btn_error_ok = Button(" Понятно (OK) ", [this] { setScreen(Screen::List); });
        auto error_buttons = Container::Horizontal({ btn_error_ok });
        error_container_ = Container::Vertical({ error_buttons });
    }

    Element App::renderListScreen()
    {
        return vbox({
            text(" ГАЛЕРЕЯ ИЗОБРАЖЕНИЙ ") | bold | center,
            separator(),
            menu_->Render() | border | flex,
            separator(),
            list_container_->ChildAt(1)->Render() | center,
            separator(),
            text(" ↑↓ - навигация | Enter - детали | Tab - смена фокуса ") | dim | center
        }) | border;
    }

    Element App::renderDetailsScreen()
    {
        std::string title = "Неизвестно";
        for (const auto& img : images_)
        {
            if (img.uuid == selected_uuid_) { title = img.name; break; }
        }

        Elements rows = {
            text(" ДЕТАЛИ ИЗОБРАЖЕНИЯ ") | bold | center,
            separator(),
            text(" Название: " + title) | center,
            text(" UUID: " + selected_uuid_) | center | dim,
            separator(),
            details_container_->ChildAt(0)->Render() | center,
            separator()
        };
        if (!status_text_.empty())
        {
            rows.push_back(text(" " + status_text_ + " ") | center | dim);
        }
        rows.push_back(text(" (Удаление необратимо!) ") | center | color(Color::Red));

        return vbox(rows) | border;
    }

    Element App::renderCreateScreen()
    {
        const std::size_t n = mouse_.size();
        const std::size_t max = utils::MouseTracker::kMaxPoints;
        const bool enough = n >= kMinMousePoints;

        return vbox({
            text(" СОЗДАНИЕ ИЗОБРАЖЕНИЯ ") | bold | center,
            separator(),
            text(" Введите название: ") | center,
            create_container_->ChildAt(0)->Render() | border | center,
            separator(),
            text(" Водите мышью по экрану - так создаётся сид ") | center,
            hbox({
                text(" Движений: " + std::to_string(n) + "/" + std::to_string(max) + " "),
                gauge(static_cast<float>(n) / static_cast<float>(max)) | flex
                    | color(enough ? Color::Green : Color::Yellow),
            }) | border,
            separator(),
            create_container_->ChildAt(1)->Render() | center,
            separator(),
            text(enough ? " Сид готов. " : " Подвигайте мышью (минимум " +
                 std::to_string(kMinMousePoints) + " точек). ") | dim | center
        }) | border;
    }

    Element App::renderErrorScreen()
    {
        return vbox({
            text(" ОШИБКА! ") | bold | center,
            separator(),
            text(error_text_) | center | color(Color::Red),
            separator(),
            error_container_->ChildAt(0)->Render() | center,
            separator(),
            text(" ESC - тоже работает ") | dim | center
        }) | border;
    }

    Component App::getComponent(ScreenInteractive& screen)
    {
        initComponents(screen);

        auto tab = Container::Tab({
            list_container_,
            details_container_,
            create_container_,
            error_container_
        }, &screen_index_);

        auto renderer = Renderer(tab, [this] {
            switch (this->screen()) {
                case Screen::List:    return renderListScreen();
                case Screen::Details: return renderDetailsScreen();
                case Screen::Create:  return renderCreateScreen();
                case Screen::Error:   return renderErrorScreen();
            }
            return text("Ошибка") | center;
        });

        return CatchEvent(renderer, [this, &screen](Event event) {
            if (event.is_mouse()) 
            {
                if (this->screen() == Screen::Create) 
                {
                    mouse_.addPoint(event.mouse().x, event.mouse().y);
                }
                return false;
            }

            if (event == Event::ArrowUp || event == Event::ArrowDown || event == Event::Tab) 
            {
                return false;
            }
            const Screen s = this->screen();

            if (s != Screen::Create && (event == Event::Character('q') || event == Event::Character('Q'))) 
            {
                screen.Exit(); return true;
            }

            if (s == Screen::List) 
            {
                if (event == Event::Character('r') || event == Event::Character('R')) { loadImages(); return true; }
                if (event == Event::Character('n') || event == Event::Character('N')) { openCreate(); return true; }
            }
            else if (s == Screen::Details) 
            {
                if (event == Event::Escape) 
                { 
                    setScreen(Screen::List); return true; 
                }
                if (event == Event::Delete) 
                { 
                    deleteImage(); return true; 
                }
                if (event == Event::Character('o') || event == Event::Character('O')) 
                {
                    openImage(selected_uuid_); return true;
                }
            }
            else if (s == Screen::Create) 
            {
                if (event == Event::Escape) 
                {
                    new_name_.clear(); mouse_.clear(); setScreen(Screen::List); return true;
                }
            }
            else if (s == Screen::Error) 
            {
                if (event == Event::Escape) { setScreen(Screen::List); return true; }
            }
            return false;
        });
    }
}
