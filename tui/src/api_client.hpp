#ifndef API_CLIENT_HPP
#define API_CLIENT_HPP
#include <string>
#include "models/models.hpp"

namespace api
{
    class ApiClient
    {
    public:
        explicit ApiClient(const std::string& base_url);

        bool checkHealth();

        models::ImageMeta createImage(const models::CreateImageRequest&);
        models::ListImagesResponse listImages(int limit, const std::string& start_uuid);
        models::ImageMeta getImage(const std::string& uuid);
        models::ImageMeta updateImage(const std::string& uuid, const models::UpdateImageRequest&);
        void deleteImage(const std::string& uuid);
        std::string downloadImage(const std::string& uuid);

    private:
        std::string base_url_;
        std::string extractError(const std::string& response_body, int status_code);
    };
}

#endif