#include <AsyncMqttClient.h>
#include <FastLED.h>
#include <WiFi.h>
#include <rdm6300.h>
#include "messages.h"
#include "Arduino.h"
#include "config.h"
#include "readChip.hpp"
#include "state.h"

AsyncMqttClient mqttClient;
TimerHandle_t mqttReconnectTimer;
TimerHandle_t wifiReconnectTimer;

#define LED_PIN 12
#define RDM6300_RX_PIN 5

// #define NUM_LEDS 11
#define NUM_LEDS 8
#define BRIGHTNESS 64
#define LED_TYPE WS2812
CRGB leds[NUM_LEDS];

#define UPDATES_PER_SECOND 100

String ownID = "63D592A9";
String partnerID = "";

barmband::state::bandState currentState = barmband::state::startup;

Rdm6300 rdm6300;

// packet id of the last registration message sent
uint16_t registrationPacketId = 0;

void setState(barmband::state::bandState newState) {
  Serial.printf("New state: %s\n", barmband::state::bandStateNames[newState]);
  currentState = newState;
}

void connectToMqtt() {
  Serial.println("Connecting to MQTT...");
  mqttClient.connect();
}

void onMqttConnect(bool sessionPresent) {
  Serial.println("Connected to MQTT.");
  Serial.print("Session present: ");
  Serial.println(sessionPresent);
  uint16_t packetIdSub = mqttClient.subscribe(MQTT_TOPIC, 2);
  Serial.print("Subscribing at QoS 2, packetId: ");
  Serial.println(packetIdSub);

  // subscribe to topics to be able to receive messages
  mqttClient.subscribe("barmband/setup", 0);
  mqttClient.subscribe("barmband/challenge", 0);

  // Registration
  char message[16];
  sprintf(message, "Hello %s", ownID);
  Serial.println(message);
  registrationPacketId =  mqttClient.publish("barmband/setup", 1, true, message);
}

void onMqttDisconnect(AsyncMqttClientDisconnectReason reason) {
  Serial.printf("Disconnected from MQTT: %d\n", reason);

  if (WiFi.isConnected()) {
    xTimerStart(mqttReconnectTimer, 0);
  }
}

void onMqttSubscribe(uint16_t packetId, uint8_t qos) {
  Serial.println("Subscribe acknowledged.");
  Serial.print("  packetId: ");
  Serial.println(packetId);
  Serial.print("  qos: ");
  Serial.println(qos);
}

void onMqttUnsubscribe(uint16_t packetId) {
  Serial.println("Unsubscribe acknowledged.");
  Serial.print("  packetId: ");
  Serial.println(packetId);
}

void onMqttMessage(char *topic, char *payload,
                   AsyncMqttClientMessageProperties properties, size_t len,
                   size_t index, size_t total) {
  Serial.println("Publish received.");
  Serial.print("  topic: ");
  Serial.println(topic);
  Serial.print("  qos: ");
  Serial.println(properties.qos);
  Serial.print("  dup: ");
  Serial.println(properties.dup);
  Serial.print("  retain: ");
  Serial.println(properties.retain);
  Serial.print("  len: ");
  Serial.println(len);
  Serial.print("  index: ");
  Serial.println(index);
  Serial.print("  total: ");
  Serial.println(total);

  String msg(payload, len);

  if (strcmp(topic, "setup") == 0) {
    // msg should contain basic setup stuff (whatever that could be)
    Serial.println(msg);
  }

  if (strcmp(topic, "barmband/challenge") == 0) {
    // msg should contain two IDs who are then searching for each other
    Serial.println(msg);
    
    auto newPairMessage = barmband::messages::parseNewPairMessage(msg);
    if (newPairMessage.isOk) {
      Serial.println("got new pair message");

      if (currentState == barmband::state::waiting) {
        if (newPairMessage.firstBandId == ownID) {
          Serial.println("It's for me!");
          partnerID = newPairMessage.secondBandId;
        } else if (newPairMessage.secondBandId == ownID) {
          Serial.println("It's for me!");
          partnerID = newPairMessage.firstBandId;
        }
        Serial.printf("New partner: %s\n", partnerID);
        setState(barmband::state::paired);
      }
      
    }

    // TODO: don't run other parsers when one succeeds
    auto abortMessage = barmband::messages::parseAbortMessage(msg);
    if (abortMessage.isOk) {
      Serial.println("got abort message");

      if (currentState == barmband::state::paired && abortMessage.bandId == partnerID) {
        // TODO: notify user
        Serial.println("partner aborted challenge");
        setState(barmband::state::idle);
      }
    }

    auto pairFoundMessage = barmband::messages::parsePairFoundMessage(msg);
    if (pairFoundMessage.isOk) {
      Serial.println("got pair found message");

      if (currentState == barmband::state::paired && pairFoundMessage.firstBandId == ownID || pairFoundMessage.secondBandId == ownID ) {
        // TODO: notify user
        Serial.println("partner found me");
        setState(barmband::state::idle);
      }
    }
  }

  if (strcmp(topic, "scan") == 0) {
    // msg should contain two IDs who just matched
    // these IDs should not match again in the future
    Serial.println(msg);
  }
  // Serial.println(msg);
}

void onMqttPublish(uint16_t packetId) {
  Serial.println("Publish acknowledged.");
  Serial.print("  packetId: ");
  Serial.println(packetId);
  if (packetId == registrationPacketId) {
    Serial.println("registration message send");
    registrationPacketId = 0;
  }
}

void WiFiEvent(WiFiEvent_t event) {
  Serial.printf("[WiFi-event] event: %d\n", event);
  switch (event) {
    case SYSTEM_EVENT_STA_GOT_IP:
      Serial.println("WiFi connected");
      Serial.println("IP address: ");
      Serial.println(WiFi.localIP());
      connectToMqtt();
      break;
    case SYSTEM_EVENT_STA_DISCONNECTED:
      Serial.println("WiFi lost connection");
      xTimerStop(
          mqttReconnectTimer,
          0);  // ensure we don't reconnect to MQTT while reconnecting to Wi-Fi
      xTimerStart(wifiReconnectTimer, 0);
      break;
  }
}

void connectToWifi() { WiFi.begin(WIFI_SSID, WIFI_PASSWORD); }

void setup() {
  Serial.begin(9600);
  setState(barmband::state::startup);

  mqttReconnectTimer =
      xTimerCreate("mqttTimer", pdMS_TO_TICKS(2000), pdFALSE, (void *)0,
                   reinterpret_cast<TimerCallbackFunction_t>(connectToMqtt));
  wifiReconnectTimer =
      xTimerCreate("wifiTimer", pdMS_TO_TICKS(2000), pdFALSE, (void *)0,
                   reinterpret_cast<TimerCallbackFunction_t>(connectToWifi));

  WiFi.onEvent(WiFiEvent);

  mqttClient.onConnect(onMqttConnect);
  mqttClient.onDisconnect(onMqttDisconnect);
  mqttClient.onSubscribe(onMqttSubscribe);
  mqttClient.onUnsubscribe(onMqttUnsubscribe);
  mqttClient.onMessage(onMqttMessage);
  mqttClient.onPublish(onMqttPublish);
  mqttClient.setServer(MQTT_HOST, MQTT_PORT);

  connectToWifi();

  init();

  FastLED.addLeds<WS2812, LED_PIN, RGB>(leds,
                                        NUM_LEDS);  // GRB ordering is typical
  FastLED.setBrightness(BRIGHTNESS);

  rdm6300.begin(RDM6300_RX_PIN);

  Serial.println("\nrdm6300 started...\n");
}

void loop() {
  // todo: blink/wa

  byte id = rdm6300.get_tag_id();

  if (id != 0) {
    Serial.println(id);
    /*
    char *buff = (char *)malloc(30);
    sprintf(buff, "{ownID: %s, message: Scanned card with ID %s",
    mqttClient.getClientId(), id);

    mqttClient.publish(MQTT_TOPIC, 1, true, buff);

    // check if target id is correct
    char *buff = (char *)malloc(25);

    sprintf(buff, "scanned tag %s", id);

    mqttClient.publish(MQTT_TOPIC, 1, true, buff);
    */
  } else {
    // solid color
    for (int i = 0; i < NUM_LEDS; i++) {
      leds[i] = CRGB::Cyan;
    }
  }
  FastLED.show();
}
