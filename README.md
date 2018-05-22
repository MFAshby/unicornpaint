# Unicorn Paint
A painting webapp, controlling a [16x16 LED matrix](https://shop.pimoroni.com/products/unicorn-hat-hd). It allows multiple people to paint at once & updates in real time.

This project includes the following components:
* A [Go library](https://github.com/MFAshby/unicornpaint/tree/master/unicorn) for controlling the physical device 
* A [Go web server](https://github.com/MFAshby/unicornpaint/blob/master/Server.go) utilising websocket communication
* A [React JS client application](https://github.com/MFAshby/unicornpaint/tree/master/src)
* Distribution & deployment with [Docker](https://github.com/MFAshby/unicornpaint/blob/master/Dockerfile)

Live example hosted at:
https://unicorn.mfashby.net

![screenshot](https://github.com/MFAshby/unicornpaint/raw/master/Screenshot.png "Screenshot")

## Running
Ensure SPI is enabled on your raspberry pi:
```
sudo raspi-config nonint do_spi 0
sudo reboot
```

Run with docker:
```
docker pull mfashby/unicornpaint
docker run --device /dev/spidev0.0:/dev/spidev0.0 --publish 3001:3001 mfashby/unicornpaint
```

Or download the [docker-compose](https://raw.githubusercontent.com/MFAshby/unicornpaint/master/docker-compose.yml) file and run with 
```
docker-compose up -d
```

Or install Go, clone the repository & run directly with 
```
go get github.com/ecc1/spi github.com/gorilla/websocket github.com/MFAshby/unicornpaint/unicorn
go run Server.go
```
