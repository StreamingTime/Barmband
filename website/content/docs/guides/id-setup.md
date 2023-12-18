---
title: "Initial Setup"
date: 2023-12-18T10:00:00+02:00
lastmod: 2023-12-18T10:00:00+02:00
draft: false
weight: 27
toc: true
---


Each Barmband has to know the identifier of its own RFID tag, which will be used to identify the Barmband both in the physical world (by scanning the tag) and in the MQTT communication.

During setup mode, the Barmband will scan a presented RFID tag and use the ID as its own.
The Barmband stores the ID in its own non-volatile storage, which is persistent between reboots.

## Setting the Barmband ID

- Enter the setup mode by restarting the barmband while the button is pressed. Setup mode is entered automatically when the ID is not configured yet
- Scan an RFID tag
- The barmband will restart automatically into normal mode

