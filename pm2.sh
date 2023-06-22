go build -a -installsuffix cgo -ldflags="-w -s" main.go
sudo chmod +x main
sudo pm2 delete LARA
sudo pm2 start ./main --name LARA