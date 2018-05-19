A pretty basic painting app for the Raspberry Pi Unicorn HD hat. 
Allows multiple people to paint at once & updates in real time.

Run with docker:
```
docker pull mfashby/unicornpaint
docker run --device /dev/spidev0.0:/dev/spidev0.0 --publish 3001:3001 mfashby/unicornpaint
```
Or clone the repository and use the example docker-compose.yml file to build & run

![alt text](https://github.com/MFAshby/unicornpaint/raw/master/Screenshot.png "Screenshot")
