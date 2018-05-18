A pretty basic painting app for the Raspberry Pi Unicorn HD hat. 
Allows multiple people to paint at once & updates in real time.

To run: 
`curl "https://raw.githubusercontent.com/MFAshby/unicornpaint/master/download_and_run.sh" | sh`
Open your browser to http://<your raspberry IP address>:3001/

Alternatively, run with docker:
```
docker pull mfashby/unicornpaint
docker run --device /dev/spidev0.0:/dev/spidev0.0 --publish 3001:3001 mfashby/unicornpaint
```
Or use the example docker-compose.yml file 