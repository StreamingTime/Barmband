#include <Arduino.h>
#include "messages.h"

namespace barmband::messages {

NewPairMessage parseNewPairMessage(String message) {

    NewPairMessage msg;
    msg.isOk = false;

    char bandACstr[9];
    char bandBCstr[9];

    size_t n = sscanf(message.c_str(), "New pair %s %s", bandACstr, bandBCstr);

    if (n != 2) {
        return msg;
    }

    msg.firstBandId = String(bandACstr);
    msg.secondBandId = String(bandBCstr);
    msg.isOk = true;

    Serial.println(msg.firstBandId);
    Serial.println(msg.secondBandId);
    return msg;
}

}