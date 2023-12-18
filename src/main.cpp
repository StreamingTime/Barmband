#include <AsyncMqttClient.h>
#include <Preferences.h>
#include <WiFi.h>
#include <rdm6300.h>

#include "Arduino.h"
#include "config.h"
#include "ledController.h"
#include "logging.h"
#include "messages.h"
#include "ota_update.h"
#include "state.h"

AsyncMqttClient mqttClient;
TimerHandle_t mqttReconnectTimer;
TimerHandle_t wifiReconnectTimer;

#define RDM6300_RX_PIN 5
#define BUTTON_PIN 4
String ownID = "";
String partnerID = "";
uint32_t color = 0;

Preferences preferences;

barmband::state::bandState currentState = barmband::state::startup;

int buttonLastState = HIGH;
int buttonCurrentState;  // the previous state from the input pin

// button debouncing
const unsigned long MIN_DEBOUNCE_TIME = 1500;  // in millis
unsigned long buttonLastActivationTime;
Rdm6300 rdm6300;

// packet id of the last registration message sent
uint16_t registrationPacketId = 0;

void setState(barmband::state::bandState newState) {
  barmband::log::logf(ownID, "New state: %s\n",
                      barmband::state::bandStateNames[newState]);
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
  registrationPacketId =
      mqttClient.publish(MQTT_SETUP_TOPIC, MQTT_QOS, true, message);

  barmband::log::setLoggingMqttclient(&mqttClient);
}

void onMqttDisconnect(AsyncMqttClientDisconnectReason reason) {
  barmband::log::logf(ownID, "Disconnected from MQTT: %d\n", reason);

  if (WiFi.isConnected()) {
    xTimerStart(mqttReconnectTimer, 0);
  }
}

void onMqttSubscribe(uint16_t packetId, uint8_t qos) {
  barmband::log::logln(ownID, "Subscribe acknowledged.");
}

void onMqttUnsubscribe(uint16_t packetId) {
  barmband::log::logln(ownID, "Unsubscribe acknowledged.");
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
          color = newPairMessage.color;
        } else if (newPairMessage.secondBandId == ownID) {
          Serial.println("It's for me!");
          partnerID = newPairMessage.firstBandId;
          color = newPairMessage.color;
        }
        barmband::log::logf(ownID, "New partner: %s\n", partnerID);
        setState(barmband::state::paired);
      }
    }

    // TODO: don't run other parsers when one succeeds
    auto abortMessage = barmband::messages::parseAbortMessage(msg);
    if (abortMessage.isOk) {
      Serial.println("got abort message");

      if (currentState == barmband::state::paired &&
          abortMessage.bandId == partnerID) {
        barmband::log::logln(ownID, "partner aborted challenge");
        setState(barmband::state::idle);
      }
    }

    auto pairFoundMessage = barmband::messages::parsePairFoundMessage(msg);
    if (pairFoundMessage.isOk) {
      Serial.println("got pair found message");

      bool isOwnID = pairFoundMessage.firstBandId == ownID ||
                     pairFoundMessage.secondBandId == ownID;

      bool isPartnerID = pairFoundMessage.firstBandId == partnerID ||
                         pairFoundMessage.secondBandId == partnerID;

      if (currentState == barmband::state::paired && isOwnID && isPartnerID) {
        barmband::log::logln(ownID, "partner found me");
        setState(barmband::state::idle);
      }
    }

    if (!newPairMessage.isOk && !abortMessage.isOk && !pairFoundMessage.isOk) {
      barmband::log::logf(ownID, "Unknown message '%s' in topic %s\n",
                          msg.c_str(), topic);
    }
  }
}

void onMqttPublish(uint16_t packetId) {
  if (packetId == registrationPacketId) {
    Serial.println("registration message sent");
    registrationPacketId = 0;
    setState(barmband::state::idle);
  }
}

void WiFiEvent(WiFiEvent_t event) {
  Serial.printf("[WiFi-event] event: %d\n", event);
  switch (event) {
    case SYSTEM_EVENT_STA_GOT_IP:
      barmband::log::logln(ownID, "WiFi connected");
      Serial.println("IP address: ");
      Serial.println(WiFi.localIP());
      connectToMqtt();
      break;
    case SYSTEM_EVENT_STA_DISCONNECTED:
      barmband::log::logln(ownID, "WiFi lost connection");
      xTimerStop(
          mqttReconnectTimer,
          0);  // ensure we don't reconnect to MQTT while reconnecting to Wi-Fi
      xTimerStart(wifiReconnectTimer, 0);
      break;
  }
}

void scanNewId() {
  uint32_t tagID = 0;

  while (tagID == 0) {
    Serial.println("Waiting for RFID tag...");
    tagID = rdm6300.get_new_tag_id();
    delay(1000);
  }

  char newID[9];
  sprintf(newID, "%08X", tagID);
  preferences.putString("ownID", newID);
}

void connectToWifi() { WiFi.begin(WIFI_SSID, WIFI_PASSWORD); }

void setup() {
  Serial.begin(9600);

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

  buttonCurrentState = digitalRead(BUTTON_PIN);

  preferences.begin("barmband", false);

  ownID = preferences.getString("ownID", "");

  if (ownID == "" || buttonCurrentState == HIGH) {
    Serial.println("Scanning new tag ID");
    scanNewId();
    ESP.restart();
  } else {
    Serial.println("Own ID found in preferences: " + ownID);
  }

  initOtaUpdate(ownID);


  setState(barmband::state::startup);
}

void loop() {

  server.handleClient();

  buttonCurrentState = digitalRead(BUTTON_PIN);

  handleLED(currentState, color);

  uint32_t id = rdm6300.get_new_tag_id();

  buttonCurrentState = digitalRead(BUTTON_PIN);
  if (id != 0) {
    barmband::log::logf(ownID, "Scanned tag: %08X", id);
    if (currentState == barmband::state::paired) {
      char message[29];
      sprintf(message, "Pair found %s %08X", ownID, id);
      mqttClient.publish(MQTT_CHALLENGE_TOPIC, 1, true, message);
    }
  }
  if (buttonLastState == LOW && buttonCurrentState == HIGH &&
      millis() - buttonLastActivationTime > MIN_DEBOUNCE_TIME) {
    buttonLastActivationTime = millis();
    barmband::log::logln(ownID, "Button input detected");
    switch (currentState) {
        // Request pardner 🤠
      case (barmband::state::idle):
        char messageIdle[25];
        sprintf(messageIdle, "Request partner %s", ownID);
        Serial.println(messageIdle);
        mqttClient.publish(MQTT_CHALLENGE_TOPIC, MQTT_QOS, true, messageIdle);
        setState(barmband::state::waiting);
        break;

      // Abort when waiting or paired
      case (barmband::state::paired):
      case (barmband::state::waiting):
        char messageAbort[15];
        sprintf(messageAbort, "Abort %s", ownID);
        Serial.println(messageAbort);
        mqttClient.publish(MQTT_CHALLENGE_TOPIC, MQTT_QOS, true, messageAbort);
        setState(barmband::state::idle);
        break;
    }
  }
  buttonLastState = buttonCurrentState;
}
