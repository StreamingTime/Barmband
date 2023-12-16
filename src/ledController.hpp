#include <FastLED.h>

#include "state.h"

#define LED_PIN 12
#define NUM_LEDS 8
#define BRIGHTNESS 20
#define LED_TYPE WS2812
CRGB leds[NUM_LEDS];

uint8_t brightness = BRIGHTNESS;

void solidColor(CRGB color) {
  for (int i = 0; i < NUM_LEDS; i++) {
    leds[i] = color;
  }
  FastLED.show();
}

void breathingAnimation() {
  int breathingValue = (exp(sin(millis() / 2000.0 * PI)) - 0.36787944) * 108.0;
  brightness = map(breathingValue, 0, 255, 0, 255);

  for (int i = 0; i < NUM_LEDS; i++) {
    leds[i] = CRGB(brightness, brightness, brightness);
  }

  FastLED.show();
}

void initLED() {
  FastLED.addLeds<WS2812, LED_PIN, RGB>(leds, NUM_LEDS);
  FastLED.setBrightness(BRIGHTNESS);

  solidColor(CRGB::Black);
}

void handleLED(barmband::state::bandState state) {
  switch (state) {
    case barmband::state::startup:
      solidColor(CRGB::Blue);
      delay(1500);
      solidColor(CRGB::Black);
      delay(1500);
      solidColor(CRGB::Blue);
      delay(1500);
      solidColor(CRGB::Black);
      delay(1500);
      solidColor(CRGB::Red);
      delay(1500);
      solidColor(CRGB::Black);
      delay(1500);
      break;

    case barmband::state::idle:
      solidColor(CRGB::Black);
      break;

    case barmband::state::waiting:
      breathingAnimation();
      break;

    case barmband::state::paired:
      solidColor(CRGB::Blue);  // TODO: change color depending on match
      break;

    default:
      // TODO: show some kind of error animation?
      solidColor(CRGB::Black);
      break;
  }
}