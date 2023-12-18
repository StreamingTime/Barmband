#pragma once

#include <AsyncMqttClient.h>

namespace barmband::log {

void setLoggingMqttclient(AsyncMqttClient* client);
void logln(String id, const char*);
void logf(String id, const char* fmt, ...);

}  // namespace barmband::log