#pragma once
#include <Arduino.h>

namespace barmband::messages {

struct NewPairMessage {
  String firstBandId;
  String secondBandId;
  bool isOk;
};

struct AbortMessage {
    String bandId;
    bool isOk;
};

NewPairMessage parseNewPairMessage(String message);
AbortMessage parseAbortMessage(String message);

}