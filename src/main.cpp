#include "Arduino.h"
#include "readChip.hpp"

void setup()
{
  Serial.begin(9600);
  init();
}

void loop()
{
  read();
}

