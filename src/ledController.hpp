#include <FastLED.h>

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

void breathingAnimation(bool state) {
  if (state) {
    int breathingValue =
        (exp(sin(millis() / 2000.0 * PI)) - 0.36787944) * 108.0;
    brightness = map(breathingValue, 0, 255, 0, 255);

    for (int i = 0; i < NUM_LEDS; i++) {
      leds[i] = CRGB(brightness, brightness, brightness);
    }
  }
  FastLED.show();
}

void initLED() {
  FastLED.addLeds<WS2812, LED_PIN, RGB>(leds,
                                        NUM_LEDS);  // GRB ordering is typical
  FastLED.setBrightness(BRIGHTNESS);

  solidColor(CRGB::Black);
}