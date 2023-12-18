#pragma once
#include <WiFi.h>

#define WIFI_SSID "SSID_PLACEHOLDER"
#define WIFI_PASSWORD "PASSWORD"

#define MQTT_HOST IPAddress(192, 168, 2, 133)
#define MQTT_PORT 8083

#define MQTT_QOS 2

#define MQTT_SETUP_TOPIC "barmband/setup"
#define MQTT_CHALLENGE_TOPIC "barmband/challenge"
// this is optional
#define MQTT_LOGGING_TOPIC "barmband/logging"
