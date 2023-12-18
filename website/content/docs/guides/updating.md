---
title: "OTA Updating"
date: 2023-12-18T10:00:00+02:00
lastmod: 2023-12-18T10:00:00+02:00
draft: false
weight: 100
toc: true
---

The Barmband supports wireless updating.
The following steps assume there is an existing Barmband with the IP Address <BARMBAND_IP> as well as a `firmware.bin` file produced by PlatformIO
under `.pio/build/esp32dev/firmware.bin`.

## Using a webbroswser

Navigate your browser to `http://<BARMBAND_IP>`, select the `firmware.bin` file and upload it.

## Using cURL

To update multiple barmbands, it might be more convenient to use curl:

```shell
curl -F update=@.pio/build/esp32dev/firmware.bin http://<BARMBAND_IP>/update
```