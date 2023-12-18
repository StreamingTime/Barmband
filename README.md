# Barmband

Link: https://barmband.fb3pool.hs-flensburg.de/

Barmband is a wearable wristband designed to gamify social interactions. The concept is driven by the often daunting task of starting conversations in bars and pubs, where socialising is the lifeblood. Barmband offers an innovative solution - a fun, engaging tool that acts as an icebreaker in these social environments.

## Useful Links

- Project description: https://barmband.fb3pool.hs-flensburg.de/project/
- Get started (manual): https://barmband.fb3pool.hs-flensburg.de/get_started/
- Lookbook: https://barmband.fb3pool.hs-flensburg.de/lookbook/

[Graphical messaging and communication overview](https://www.tldraw.com/v/N9df8NquTPFi5-Oo25JAq?viewport=-170,48,1920,963&page=page:page)

## Prerequisites

### MQTT broker

In order to run the application, you need to have a MQTT broker running. You can use the `compose.yml` file in the `broker` directory to run a local MQTT broker. The broker is responsible for receiving messages from the clients and forwarding them to the appropriate subscribers. After running the broker, you have to configure the application to use the correct IP address and port of the broker. Adjust the `MQTT_HOST` and `MQTT_PORT` in `include/config.h` accordingly. The default port is `1883`. The IP address of the broker depends on your network configuration. If you are running the broker on the same machine as the application, you can use `localhost`.

Of course, you can also use a remote MQTT broker.

### Network configuration

The application uses the ESP32 WiFi module. You have to configure the WiFi module to connect to your WiFi network. Adjust the `WIFI_SSID` and `WIFI_PASSWORD` in `include/config.h` accordingly.

### Configuration

Configuration is done in `include/config.h`.
See `include/config_example.h`.

## Build and flash

```bash
pio run --target upload
```

### Serial monitor
```bash
pio device monitor
```
## Sources

- https://github.com/marvinroger/async-mqtt-client examples