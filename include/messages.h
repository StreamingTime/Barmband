#pragma once
#include <Arduino.h>

namespace barmband::messages {

struct NewPairMessage {
  String firstBandId;
  String secondBandId;
  String color;
  bool isOk;
};

struct AbortMessage {
    String bandId;
    bool isOk;
};

struct PairFoundMessage {
  String firstBandId;
  String secondBandId;
  bool isOk;
};

NewPairMessage parseNewPairMessage(String message);
AbortMessage parseAbortMessage(String message);
PairFoundMessage parsePairFoundMessage(String message);

}