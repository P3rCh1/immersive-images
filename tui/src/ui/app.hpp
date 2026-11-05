#ifndef APP_HPP
#define APP_HPP

#include <ftxui/component/component.hpp>
#include <ftxui/component/screen_interactive.hpp>
#include <vector>
#include <string>
#include "../api_client.hpp"

namespace ui
{
    enum class Screen
    {
        List,
        Details,
        Error
    };

    class App
    {
    public:
        explicit App(api::ApiClient& client);
        ftxui::Component getComponent(ftxui::ScreenInteractive& screen);

    private:
        api::ApiClient& client_;
        Screen current_screen_ = Screen::List;
        std::string selected_uuid_ = "";
        std::string error_text_ = "";
        
        std::vector<models::ImageMeta> images_;
        std::vector<std::string> display_items_;
        int selected_index_ = 0;
        
        ftxui::Component menu_;
        ftxui::Component list_container_;
        ftxui::Component details_container_;

        void loadImages();
        void deleteImage();
        void initComponents(ftxui::ScreenInteractive& screen);

        ftxui::Element renderListScreen();
        ftxui::Element renderDetailsScreen();
        ftxui::Element renderErrorScreen();
    };
}

#endif 