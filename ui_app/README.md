# Simple UI app for visual feedback on Barmband states

> Caution! This app is very much still in development and not yet feature complete.

## How it works

Simply listens to MQTT traffic and displays information therein. Currently, only `Hello <id>` and `New pair <id1> <id2> <color>` are implemented.

## How To

Open the `.pde` file in Processing IDE and run it. Depending on which MQTT messages are being sent, there's a visual representation.

Depending on you setup, you can also compile it via Processing-java, as seen [here](https://stackoverflow.com/questions/46833666/how-to-compilerun-another-pde-file-in-processing).