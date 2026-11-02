#include "api_client.hpp"
#include <iostream>

int main() {
    // 1. Инициализируем клиент
    std::string base_url = "http://localhost:8080";
    api::ApiClient client(base_url);

    std::cout << "--- Testing API Client ---" << std::endl;

    // 2. Проверка здоровья сервера
    std::cout << "1. Checking health..." << std::endl;
    if (!client.checkHealth()) {
        std::cerr << "ERROR: Server is DOWN. Start mock_server.py first!" << std::endl;
        return 1;
    }
    std::cout << "Server is UP!" << std::endl << std::endl;

    // 3. Получение списка картинок
    std::cout << "2. Fetching images..." << std::endl;
    try {
        auto list = client.listImages();
        std::cout << "Total images: " << list.total << std::endl;
        for (const auto& img : list.items) {
            std::cout << "- UUID: " << img.uuid << " | Name: " << img.name << std::endl;
        }
    } catch (const std::exception& e) {
        std::cerr << "ERROR fetching list: " << e.what() << std::endl;
    }
    std::cout << std::endl;

    // 4. Создание новой картинки (проверка POST-запроса)
    std::cout << "3. Creating a new image..." << std::endl;
    try {
        models::CreateImageRequest req;
        req.seed = 12345;
        req.name = "My Test Image from C++";
        // Оставляем settings пустыми, сервер сам подставит дефолтные
        
        auto new_img = client.createImage(req);
        std::cout << "Created successfully! UUID: " << new_img.uuid << std::endl;
    } catch (const std::exception& e) {
        std::cerr << "ERROR creating image: " << e.what() << std::endl;
    }

    return 0;
}