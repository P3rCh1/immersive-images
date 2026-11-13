sudo apt update
sudo apt install -y build-essential cmake libcurl4-openssl-dev libssl-dev

cd immersive-images/tui
cmake -B build
cmake --build build -j$(nproc)

API_URL="http://localhost:8080" ./build/tui_client
