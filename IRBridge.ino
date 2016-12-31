#include <IRremote.h>


IRrecv irrecv(A2);
IRsend irsend;

String buf = String("");

struct Entry {
  String key;
  decode_type_t type;
  int bits;
  const unsigned long long value;
};

const Entry TABLE[] = {
  {"dac-on",       NEC,  32, 0x60ff11ee},
  {"dac-off",      NEC,  32, 0x60ff916e},
  {"dac-vol-up",   NEC,  32, 0x60ff936c},
  {"dac-vol-down", NEC,  32, 0x60ffa35c},
  {"dac-opt1",     NEC,  32, 0x60ff51ae},
  {"dac-opt2",     NEC,  32, 0x60ffd12e},
  {"dac-usb",      NEC,  32, 0x60ff13ec},
  {"dac-aes",      NEC,  32, 0x60ff619e},
  {"dac-coax1",    NEC,  32, 0x60ff23dc},
  {"dac-coax2",    NEC,  32, 0x60ffe11e},
  {"dac-bypass",   NEC,  32, 0x60ff639c},
  {"tv-green",     SONY, 15, 0x32e9},
  {"tv-red",       SONY, 15, 0x52e9},
  {"tv-blue",      SONY, 15, 0x12e9},
  {"tv-yellow",    SONY, 15, 0x72e9},
  {"tv-chan-up",   SONY, 12, 0x090},
  {"tv-chan-down", SONY, 12, 0x890},
  {"tv-vol-up",    SONY, 12, 0x490},
  {"tv-vol-down",  SONY, 12, 0xC90},
};


void setup()
{
  Serial.begin(9600);
  Serial.println("setup");
  irrecv.enableIRIn();
}


void send(decode_type_t type, int bits, unsigned long code) {
  Serial.print("sending ");
  Serial.println(code, HEX);

  for (int i = 0; i < 3; i++) {
    irsend.sendNEC(code, bits);
    delay(40);
  }
}


void loop() {
  while (Serial.available() > 0) {
    char data = Serial.read();
    if (data != '\n' && data != '\0') {
      buf += data;
      continue;
    }
    
    Serial.print("received ");
    Serial.println(buf);
    
    const Entry *code = 0;
    for (int i = 0; i < sizeof(TABLE) / sizeof(Entry); i++) {
      if (TABLE[i].key == buf) {
        code = &TABLE[i];
        break;
      }
    }

    if (code != 0)
      send(code->type, code->bits, code->value);
    
    buf = String("");
    setup();
  }
  
  decode_results results;
  if (!irrecv.decode(&results)) {
    return;
  }
  
  if (results.decode_type <= 0 || results.bits == 0) {
    irrecv.resume();
    return;
  }

  Serial.print(results.decode_type);
  Serial.print(" (");
  Serial.print(results.bits);
  Serial.print(" bits) 0x");
  Serial.println(results.value, HEX);

  switch (results.value) {
    // green: on
    case 0x32e9:
      send(NEC, 32, 0x60ff11ee);
      break;

    // red: off
    case 0x52e9:
      send(NEC, 32, 0x60ff916e);
      break;
    
    // yellow: optical 1
    case 0x72e9:
      send(NEC, 32, 0x60ff51ae);
      break;
    
    // blue: optical 2
    case 0x12e9:
      send(NEC, 32, 0x60ffd12e);
      break;
    
    // channel up: volume up
    case 0x090:
      send(NEC, 32, 0x60ff936c);
      break;
    
    // channel down: volume down
    case 0x890:
      send(NEC, 32, 0x60ffa35c);
      break;
    
    default:
      break;
  }

  // receiver seems to get into a bad state after sending
  setup();
}

