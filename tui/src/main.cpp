#include "ui/app.hpp"
#include <ftxui/component/screen_interactive.hpp>
#include <iostream>

int main()
{
    std::string base_url = "http://localhost:8080";
    api::ApiClient client(base_url);

    std::cout << "Подключение к серверу...\n";
    if (!client.checkHealth())
    {
        std::cerr << "Ошибка сервера!\n";
        return 1;
    }

    auto screen = ftxui::ScreenInteractive::Fullscreen();
    
    ui::App app(client);
    auto component = app.getComponent(screen);

    screen.Loop(component);
    screen.Exit();

    return 0;
}