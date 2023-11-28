# Broker

The `/broker` directory contains a `compose.yml` to run a local MQTT broker. The broker is responsible for receiving messages from the clients and forwarding them to the appropriate subscribers. We use the [Eclipse Mosquitto](https://mosquitto.org/) broker.