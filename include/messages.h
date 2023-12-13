#pragma once
#include <Arduino.h>

namespace barmband::messages {

struct NewPairMessage {
  String firstBandId;
  String secondBandId;
  bool isOk;
};

NewPairMessage parseNewPairMessage(String message);

}