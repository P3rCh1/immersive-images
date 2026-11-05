#include "app.hpp"
#include <ftxui/dom/elements.hpp>

namespace ui
{
    using namespace ftxui;

    App::App(api::ApiClient& client):
        client_(client)
    {
        loadImages();

        menu_ = Menu(&display_items_, &selected_index_);       
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
            display_items_.clear();
            display_items_.push_back("Ошибка загрузки: " + std::string(e.what()));
        }
    }

    Component App::getComponent(ScreenInteractive& screen)
    {
        btn_refresh_ = Button(" Обновить (R) ", [&] { loadImages(); });
        btn_create_  = Button(" Создать (N) ", [&] {  });
        btn_quit_    = Button(" Выход (Q) ", [&screen] { screen.Exit(); });

        auto buttons = Container::Horizontal({
            btn_refresh_,
            btn_create_,
            btn_quit_
        });

        main_container_ = Container::Vertical({
            menu_,
            buttons
        });

        auto renderer = Renderer(main_container_, [&] 
        {
            return vbox({
                text(" ГАЛЕРЕЯ ИЗОБРАЖЕНИЙ ") | bold | center,
                separator(),
                menu_->Render() | border | flex, 
                separator(),
                hbox({
                    btn_refresh_->Render() | flex,
                    btn_create_->Render() | flex,
                    btn_quit_->Render() | flex
                }),
                separator(),
                text(" ↑↓ - навигация | Tab - переход к кнопкам ") | dim | center
            }) | border;
        });

        auto component = CatchEvent(renderer, [&](Event event)
        {
            if (event.is_mouse())
            {
                return false;
            }

            if (event == Event::Character('q') || event == Event::Character('Q'))
            {
                screen.Exit();
                return true;
            }
            if (event == Event::Character('r') || event == Event::Character('R'))
            {
                loadImages();
                return true;
            }
            if (event == Event::Character('n') || event == Event::Character('N'))
            {
                return true;
            }

            return false;
        });

        return component;
    }

}