#pragma once
#include <cstdint>

#include "state.h"

void initLED();

void handleLED(barmband::state::bandState state, uint32_t hexColor);