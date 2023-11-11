#include <Arduino.h>
#include <SPI.h>
#include <MFRC522.h>

#define SS_PIN 5
#define RST_PIN 13

const int LED_PIN = 26;

MFRC522 rfid(SS_PIN, RST_PIN);

// Init array that will store new NUID
byte nuidPICC[4];

// Correct tag hex value
byte correctTag[4] = { 0x7A, 0x67, 0x06, 0xB1 };

void checkTarget();
bool compareByteArrays(byte *array1, byte *array2);
void printHex(byte *buffer, byte bufferSize);

void setup() {
  Serial.begin(9600);
  SPI.begin();      // Init SPI bus
  rfid.PCD_Init();  // Init MFRC522

  pinMode(LED_PIN, OUTPUT);
}

void loop() {
  // Reset the loop if no new card present on the sensor/reader. This saves the entire process when idle.
  if (!rfid.PICC_IsNewCardPresent()) {
    return;
  }

  // Verify if the NUID has been readed
  if (!rfid.PICC_ReadCardSerial()) {
    return;
  }

  Serial.println(F("Tag has been detected."));

  Serial.print(F("PICC type: "));
  MFRC522::PICC_Type piccType = rfid.PICC_GetType(rfid.uid.sak);
  Serial.println(rfid.PICC_GetTypeName(piccType));

  // Check is the PICC of Classic MIFARE type
  if (piccType != MFRC522::PICC_TYPE_MIFARE_MINI && piccType != MFRC522::PICC_TYPE_MIFARE_1K && piccType != MFRC522::PICC_TYPE_MIFARE_4K) {
    Serial.println(F("Your tag is not of type MIFARE Classic."));
    return;
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

bool compareByteArrays(byte *array1, byte *array2) {
  if (sizeof(array1) != sizeof(array2)) {
    return false;
  }

  bool equal = false;

  for (int i = 0; i < sizeof(array1) && !equal; i++) {
    if (array1[i] == array2[i]) {
      equal = true;
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
