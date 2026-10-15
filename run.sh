cd backend
make server

cd build
rm -rf ~/.local/share/token-jack
JACKER_USER=admin JACKER_PASSWD=admin ./server setup
