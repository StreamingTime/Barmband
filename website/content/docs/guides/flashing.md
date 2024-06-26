---
title: "Flashing"
date: 2023-12-18T10:00:00+02:00
lastmod: 2023-12-18T10:00:00+02:00
draft: false
weight: 25
toc: true
---
### Flashing the software

We use [PlatformIO](https://platformio.org/) to build and flash the software.
Please follow the instructions provided [here](https://platformio.org/install/integration) to install PlatformIO.
We will assume that a working `PlatformIO Core` installation is available on your machine during the next steps.

- Get the source code from the [git repository](https://github.com/StreamingTime/Barmband)
- Create a copy of the `include/config_example.h` and name it `include/config.h`
- Edit this file to fit your needs, especially `WIFI_SSID`  and `WIFI_PASSWORD`. You can change the MQTT Broker using `MQTT_HOST` and `MQTT_PORT` values.
- Build and flash the software using
```shell
pio run --target upload
```
- You can view the serial output using
```shell
pio device monitor
```
