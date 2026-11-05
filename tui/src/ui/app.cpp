#include "app.hpp"
#include <ftxui/dom/elements.hpp>
#include <iostream>

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
            if (selected_index_ >= display_items_.size()) 
            {
                selected_index_ = 0;
            }
        }
        catch (const std::exception& e)
        {
            display_items_.clear();
            display_items_.push_back("Ошибка загрузки: " + std::string(e.what()));
        }
    }

    void App::deleteImage()
    {
        try
        {
            client_.deleteImage(selected_uuid_);
            loadImages();
            current_screen_ = Screen::List;
        }
        catch (const std::exception& e)
        {
            std::cerr << "Ошибка удаления: " << e.what() << "\n";
            error_text_ = e.what();
            current_screen_ = Screen::Error;
        }
    }

    void App::initComponents(ScreenInteractive& screen)
    {
        auto btn_refresh = Button(" Обновить (R) ", [&] { loadImages(); });
        auto btn_create  = Button(" Создать (N) ", [&] { /* TODO */ });
        auto btn_quit    = Button(" Выход (Q) ", [&screen] { screen.Exit(); });
        
        auto list_buttons = Container::Horizontal({ btn_refresh, btn_create, btn_quit });
        list_container_ = Container::Vertical({ menu_, list_buttons });

        auto btn_back   = Button(" Назад (ESC) ", [&] { current_screen_ = Screen::List; });
        auto btn_delete = Button(" Удалить (DEL) ", [&] { deleteImage(); });
        
        auto details_buttons = Container::Horizontal({ btn_back, btn_delete });
        details_container_ = Container::Vertical({ details_buttons });
    }


    Element App::renderListScreen()
    {
        return vbox({
            text(" ГАЛЕРЕЯ ИЗОБРАЖЕНИЙ ") | bold | center,
            separator(),
            menu_->Render() | border | flex, 
            separator(),
            list_container_->ChildAt(1)->Render(), 
            separator(),
            text(" ↑↓ - навигация | Enter - детали ") | dim | center
        }) | border;
    }

    Element App::renderDetailsScreen()
    {
        std::string title = "Неизвестно";
        for (const auto& img : images_)
        {
            if (img.uuid == selected_uuid_)
            {
                title = img.name;
                break;
            }
        }

        return vbox({
            text(" ДЕТАЛИ ИЗОБРАЖЕНИЯ ") | bold | center,
            separator(),
            text(" Название: " + title) | center,
            text(" UUID: " + selected_uuid_) | center | dim,
            separator(),
            details_container_->ChildAt(0)->Render(), 
            separator(),
            text(" (Удаление необратимо!) ") | center | color(Color::Red)
        }) | border;
    }

    Element App::renderErrorScreen()
    {
        return vbox({
            text(" ОШИБКА! ") | bold | center,
            separator(),
            text(error_text_) | center | color(Color::Red),
            separator(),
            text(" Нажмите ESC для возврата ") | dim | center
        }) | border;
    }

    Component App::getComponent(ScreenInteractive& screen)
    {
        initComponents(screen);

        auto tab = Container::Tab({ list_container_, details_container_ }, (int*)&current_screen_);

        auto renderer = Renderer(tab, [&] 
        {
            if (current_screen_ == Screen::List)
            {
                return renderListScreen();
            }
            else if (current_screen_ == Screen::Details)
            {
                return renderDetailsScreen();
            }
            else
            {
                return renderErrorScreen();
            }
        });

        return CatchEvent(renderer, [&](Event event)
        {
            if (event.is_mouse() || event == Event::ArrowUp || event == Event::ArrowDown || event == Event::Tab)
            {
                return false;
            }

            if (event == Event::Character('q') || event == Event::Character('Q'))
            {
                screen.Exit();
                return true;
            }

            if (current_screen_ == Screen::List)
            {
                if (event == Event::Character('r') || event == Event::Character('R')) 
                { 
                    loadImages(); 
                    return true;
                }
                
                if (event == Event::Return && menu_->Focused() && !images_.empty())
                {
                    selected_uuid_ = images_[selected_index_].uuid;
                    current_screen_ = Screen::Details;
                    return true;
                }
            }
            else if (current_screen_ == Screen::Details)
            {
                if (event == Event::Escape) { current_screen_ = Screen::List; return true; }
                if (event == Event::Delete) { deleteImage(); return true; }
            }
            else if (current_screen_ == Screen::Error)
            {
                if (event == Event::Escape) 
                { 
                    current_screen_ = Screen::List; 
                    return true; 
                }
            }
            return false;
        });
    }
}