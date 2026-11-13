#include <ftxui/component/screen_interactive.hpp>
#include <cstdlib>
#include <iostream>
#include "ui/app.hpp"

int main()
{
    const char* env_url = std::getenv("API_URL");
    std::string base_url = env_url ? env_url : "http://localhost:8080";

    api::ApiClient client(base_url);

    std::cout << "Подключение к серверу: " << base_url << "\n";
    if (!client.checkHealth())
    {
        std::cerr << "Ошибка сервера! Проверьте, что бекенд запущен по адресу " << base_url << "\n";
        return 1;
    }

    auto screen = ftxui::ScreenInteractive::TerminalOutput();
    
    ui::App app(client);
    auto component = app.getComponent(screen);

    screen.Loop(component);

    return 0;
}
