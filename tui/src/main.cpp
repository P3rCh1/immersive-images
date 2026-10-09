#include <ftxui/component/component.hpp>
#include <ftxui/component/screen_interactive.hpp>
#include <ftxui/dom/elements.hpp>

using namespace ftxui;

int main() 
{
    auto screen = ScreenInteractive::TerminalOutput();
    
    auto component = Renderer([&] {
        return vbox({
            text("TUI Client — Immersive Images") | bold | center,
            separator(),
            text("Если ты это видишь — FTXUI работает!") | color(Color::Green),
            separator(),
            text("Нажми Q для выхода") | dim | center,
        }) | border;
    });
    
    component |= CatchEvent([&](Event event) {
        if (event == Event::Character('q') || event == Event::Character('Q')) {
            screen.Exit();
            return true;
        }
        return false;
    });
    
    screen.Loop(component);
    return 0;
}