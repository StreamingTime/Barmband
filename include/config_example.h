#include <WiFi.h>

#define WIFI_SSID "SSID_PLACEHOLDER"
#define WIFI_PASSWORD "PASSWORD"

#define MQTT_HOST IPAddress(192, 168, 2, 133)
#define MQTT_PORT 8083

const char* MQTT_SETUP_TOPIC = "barmband/setup";
const char* MQTT_CHALLENGE_TOPIC = "barmband/challenge";
