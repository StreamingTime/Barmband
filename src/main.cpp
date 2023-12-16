#include <AsyncMqttClient.h>
#include <WiFi.h>
#include <rdm6300.h>

#include "Arduino.h"
#include "config.h"
#include "ledController.hpp"
#include "messages.h"
#include "state.h"

AsyncMqttClient mqttClient;
TimerHandle_t mqttReconnectTimer;
TimerHandle_t wifiReconnectTimer;

#define RDM6300_RX_PIN 5
#define BUTTON_PIN 4

String ownID = "63D592A9";
String partnerID = "";

barmband::state::bandState currentState = barmband::state::startup;

int buttonLastState = HIGH;
int buttonCurrentState;      // the previous state from the input pin
bool buttonPressed = false;  // prevent button event from being triggered twice
                             // (on press & on release)

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

  // subscribe to topics to be able to receive messages
  mqttClient.subscribe(MQTT_SETUP_TOPIC, 0);
  mqttClient.subscribe(MQTT_CHALLENGE_TOPIC, 0);

  // Registration
  char message[16];
  sprintf(message, "Hello %s", ownID);
  Serial.println(message);
  registrationPacketId = mqttClient.publish(MQTT_SETUP_TOPIC, 1, true, message);
}

void onMqttDisconnect(AsyncMqttClientDisconnectReason reason) {
  Serial.printf("Disconnected from MQTT: %d\n", reason);

  if (WiFi.isConnected()) {
    xTimerStart(mqttReconnectTimer, 0);
  }
}

void onMqttSubscribe(uint16_t packetId, uint8_t qos) {
  Serial.println("Subscribe acknowledged.");
}

void onMqttUnsubscribe(uint16_t packetId) {
  Serial.println("Unsubscribe acknowledged.");
}

void onMqttMessage(char *topic, char *payload,
                   AsyncMqttClientMessageProperties properties, size_t len,
                   size_t index, size_t total) {
  Serial.printf("Publish received on topic %s\n", topic);

  String msg(payload, len);

  if (strcmp(topic, MQTT_CHALLENGE_TOPIC) == 0) {
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

      if (currentState == barmband::state::paired &&
          abortMessage.bandId == partnerID) {
        // TODO: notify user
        Serial.println("partner aborted challenge");
        setState(barmband::state::idle);
      }
    }

    auto pairFoundMessage = barmband::messages::parsePairFoundMessage(msg);
    if (pairFoundMessage.isOk) {
      Serial.println("got pair found message");

      if (currentState == barmband::state::paired &&
              pairFoundMessage.firstBandId == ownID ||
          pairFoundMessage.secondBandId == ownID) {
        // TODO: notify user
        Serial.println("partner found me");
        setState(barmband::state::idle);
      }
    }
  }
}

void onMqttPublish(uint16_t packetId) {
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
  initLED();

  rdm6300.begin(RDM6300_RX_PIN);

  Serial.println("\nrdm6300 started...\n");

  pinMode(BUTTON_PIN, INPUT_PULLDOWN);
}

void loop() {
  // todo: blink/wa
  handleLED(currentState);

  byte id = rdm6300.get_tag_id();

  buttonCurrentState = digitalRead(BUTTON_PIN);

  if (buttonLastState == LOW && buttonCurrentState == HIGH && !buttonPressed) {
    switch (currentState) {
        // Request pardner 🤠
      case (barmband::state::idle):
        char messageIdle[25];
        sprintf(messageIdle, "Request partner %s", ownID);
        Serial.println(messageIdle);
        mqttClient.publish("barmband/challenge", 1, true, messageIdle);
        buttonPressed = true;
        break;

      // Abort when waiting or paired
      case (barmband::state::paired):
      case (barmband::state::waiting):
        char messageAbort[15];
        sprintf(messageAbort, "Abort %s", ownID);
        Serial.println(messageAbort);
        mqttClient.publish("barmband/challenge", 1, true, messageAbort);
        buttonPressed = true;
        break;
    }
  }
  buttonLastState = buttonCurrentState;
}
