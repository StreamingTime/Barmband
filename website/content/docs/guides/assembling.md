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
