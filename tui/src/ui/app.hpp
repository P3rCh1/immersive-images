#ifndef APP_HPP
#define APP_HPP

#include <ftxui/component/component.hpp>
#include <ftxui/component/screen_interactive.hpp>
#include <string>
#include <vector>
#include "../api_client.hpp"
#include "../utils/mouse_tracker.hpp"

namespace ui
{
    enum class Screen
    {
        List = 0,
        Details = 1,
        Create = 2,
        Error = 3
    };

    class App
    {
    public:
        explicit App(api::ApiClient& client);
        ftxui::Component getComponent(ftxui::ScreenInteractive& screen);

    private:
        api::ApiClient& client_;
        int screen_index_ = static_cast<int>(Screen::List); // индекс для Container::Tab
        std::string selected_uuid_ = "";
        std::string error_text_ = "";
        std::string status_text_ = "";

        std::vector<models::ImageMeta> images_;
        std::vector<std::string> display_items_;
        int selected_index_ = 0;

        utils::MouseTracker mouse_;

        ftxui::Component menu_;
        ftxui::Component list_container_;
        ftxui::Component details_container_;
        ftxui::Component create_container_;
        ftxui::Component error_container_;

        std::string new_name_ = "";
        ftxui::Component input_name_;
        ftxui::Component btn_generate_;

        Screen screen() const { return static_cast<Screen>(screen_index_); }
        void setScreen(Screen s) { screen_index_ = static_cast<int>(s); }

        void loadImages();
        void deleteImage();
        void createImage();
        void openImage(const std::string& uuid);
        void openSelected();
        void openCreate();
        void initComponents(ftxui::ScreenInteractive& screen);

        ftxui::Element renderListScreen();
        ftxui::Element renderDetailsScreen();
        ftxui::Element renderCreateScreen();
        ftxui::Element renderErrorScreen();
    };
}

#endif  