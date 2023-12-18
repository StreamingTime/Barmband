---
title: "Assembling"
date: 2023-12-18T10:00:00+02:00
lastmod: 2023-12-18T10:00:00+02:00
draft: false
weight: 20
toc: true
---
## Assembling

### Connecting the parts


|      Device pin     | ESP32 board pin |
|:-------------------:|:---------------:|
| **Tag reader (RDM3600)** | |
|         5V         |        5V       |
| GND | GND |
| TX | 5 |
| **LED (WS2812)** | |
| 5V | 5V   |
| GND | GND |
| Din |        12       |
| **Button**  | |
|        | 4 |
|        | GND      |

{{< figure
  src="images/schema.drawio.png"
  alt="Schema"
>}}

### Flashing the software

We use [PlatformIO](https://platformio.org/) to build and flash the software.
Please follow the instructions provided [here](https://platformio.org/install/integration) to install PlatformIO.
We will assume that a working `PlatformIO Core` installation is available on your machine during the next steps.

- Get the source code from the [git repository](https://gitlab.hs-flensburg.de/teaching/microcontroller-programmierung-wise-23-24/barmband)
- Create a copy of the `include/config_example.h` and name it `include/config.h`
- Edit this file to fit your needs, especially `WIFI_SSID`  and `WIFI_PASSWORD`
- Build and flash the software using
```shell
pio run --target upload
```
- You can view the serial output using
```shell
pio device monitor
```
