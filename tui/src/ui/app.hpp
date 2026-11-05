#ifndef APP_HPP
#define APP_HPP
#include <ftxui/component/screen_interactive.hpp>
#include <ftxui/component/component.hpp>
#include "../api_client.hpp"

namespace ui
{
    class App
    {
    public:
        explicit App(api::ApiClient& client);
        ftxui::Component getComponent(ftxui::ScreenInteractive& screen);

    private:
        api::ApiClient& client_;

        std::vector< models::ImageMeta > images_;
        std::vector< std::string > display_items_;
        int selected_index_ = 0;

        ftxui::Component menu_;
        ftxui::Component btn_refresh_;
        ftxui::Component btn_create_;
        ftxui::Component btn_quit_;
        ftxui::Component main_container_;

        void loadImages();
    };
}

#endif 