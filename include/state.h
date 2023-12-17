#pragma once
namespace barmband::state {

enum bandState {
  startup = 0,
  idle,
  waiting,
  paired,
  registerId,
};

const char* bandStateNames[] {
  "STARTUP",
  "IDLE",
  "WAITING",
  "PAIRED",
  "REGISTERID",
};
}