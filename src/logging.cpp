#include "logging.h"
#include "config.h"

namespace barmband::log {

AsyncMqttClient* loggingMqttClient;

void setLoggingMqttclient(AsyncMqttClient* client) {
  loggingMqttClient = client;
}

void logln(String id, const char* msg) {
  Serial.println(msg);

  #ifdef MQTT_LOGGING_TOPIC
  if (loggingMqttClient != nullptr) {
    char msgBuff[80];
    sprintf(msgBuff, "%s: %s", id.c_str(), msg);
    loggingMqttClient->publish(MQTT_LOGGING_TOPIC, MQTT_QOS, true, msgBuff);
  }
  #endif
}

void logf(String id, const char* fmt, ...) {
  char msgBuff[80];

  va_list va;
  va_start(va, fmt);
  vsprintf(msgBuff, fmt, va);
  Serial.print(msgBuff);
  va_end(va);

  #ifdef MQTT_LOGGING_TOPIC
  if (loggingMqttClient != nullptr) {
    char msgWithId[91];
    sprintf(msgWithId, "%s: %s", id.c_str(), msgBuff);
    loggingMqttClient->publish(MQTT_LOGGING_TOPIC, MQTT_QOS, true, msgWithId);
  }
  #endif
}
}  // namespace barmband::log
