#include <Arduino.h>
#include <MFRC522.h>
#include <SPI.h>

#define SS_PIN 5
#define RST_PIN 13

const int LED_PIN = 26;

MFRC522 rfid(SS_PIN, RST_PIN);

const size_t TAG_ID_SIZE = 4;
// Init array that will store new NUID
byte nuidPICC[4];

// Correct tag hex value
byte correctTag[4] = {0x7A, 0x67, 0x06, 0xB1};

void checkTarget();
bool compareByteArrays(byte *array1, byte *array2);
void printHex(byte *buffer, byte bufferSize);

void init() {
  SPI.begin();      // Init SPI bus
  rfid.PCD_Init();  // Init MFRC522

  pinMode(LED_PIN, OUTPUT);
}

String read() {
  // Reset the loop if no new card present on the sensor/reader. This saves the
  // entire process when idle.
  if (!rfid.PICC_IsNewCardPresent()) {
    return "";
  }

  // Verify if the NUID has been readed
  if (!rfid.PICC_ReadCardSerial()) {
    return "";
  }

  Serial.println(F("Tag has been detected."));

  Serial.print(F("PICC type: "));
  MFRC522::PICC_Type piccType = rfid.PICC_GetType(rfid.uid.sak);
  Serial.println(rfid.PICC_GetTypeName(piccType));

  // Check is the PICC of Classic MIFARE type
  if (piccType != MFRC522::PICC_TYPE_MIFARE_MINI &&
      piccType != MFRC522::PICC_TYPE_MIFARE_1K &&
      piccType != MFRC522::PICC_TYPE_MIFARE_4K) {
    Serial.println(F("Your tag is not of type MIFARE Classic."));
    return "";
  }

  // Store NUID into nuidPICC array
  for (byte i = 0; i < 4; i++) {
    nuidPICC[i] = rfid.uid.uidByte[i];
  }

  Serial.println(F("The NUID tag is (hex):"));
  printHex(rfid.uid.uidByte, rfid.uid.size);
  checkTarget();
  Serial.println();
  Serial.println();

  // Halt PICC
  rfid.PICC_HaltA();

  // Stop encryption on PCD
  rfid.PCD_StopCrypto1();

  char s[8];

  sprintf(s, "%02X%02X%02X%02X", nuidPICC[0], nuidPICC[1], nuidPICC[2], nuidPICC[3]);
  return String(nuidPICC, TAG_ID_SIZE);
}

void checkTarget() {
  if (compareByteArrays(correctTag, nuidPICC)) {
    digitalWrite(LED_PIN, HIGH);
    Serial.println("Correct tag!");
  } else {
    digitalWrite(LED_PIN, LOW);
    Serial.println("Wrong tag!");
  }
}

/**
 * Compare the first TAG_ID_SIZE bytes
 */
bool compareByteArrays(byte *array1, byte *array2) {
  bool equal = true;

  for (int i = 0; i < TAG_ID_SIZE && equal; i++) {
    if (array1[i] != array2[i]) {
      equal = false;
    }
  }

  return equal;
}

/**
 * Helper routine to dump a byte array as hex values to Serial.
 */
void printHex(byte *buffer, byte bufferSize) {
  for (byte i = 0; i < bufferSize; i++) {
    Serial.print(buffer[i] < 0x10 ? " 0" : " ");
    Serial.print(buffer[i], HEX);
  }
  Serial.println();
}
