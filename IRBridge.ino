#include <IRremote.h>


IRrecv irrecv(A2);
IRsend irsend;

struct Entry {
  String key;
  decode_type_t type;
  int bits;
  unsigned long long value;
  String trigger;
};

const Entry TABLE[] = {
  {"dac-on",       NEC,  32, 0x60ff11ee},
  {"dac-off",      NEC,  32, 0x60ff916e},
  {"dac-dim",      NEC,  32, 0x60FF21DE},
  {"dac-vol-up",   NEC,  32, 0x60ff936c},
  {"dac-vol-down", NEC,  32, 0x60ffa35c},
  {"dac-mute",     NEC,  32, 0x60ffa15e},
  {"dac-opt1",     NEC,  32, 0x60ff51ae},
  {"dac-opt2",     NEC,  32, 0x60ffd12e},
  {"dac-usb",      NEC,  32, 0x60ff13ec},
  {"dac-aes",      NEC,  32, 0x60ff619e},
  {"dac-coax1",    NEC,  32, 0x60ff23dc},
  {"dac-coax2",    NEC,  32, 0x60ffe11e},
  {"dac-bypass",   NEC,  32, 0x60ff639c},
  {"tv-green",     SONY, 15, 0x32e9, "dac-on"},
  {"tv-red",       SONY, 15, 0x52e9, "dac-off"},
  {"tv-yellow",    SONY, 15, 0x72e9, "dac-opt1"},
  {"tv-blue",      SONY, 15, 0x12e9, "dac-opt2"},
  {"tv-chan-up",   SONY, 12, 0x090,  "dac-vol-up"},
  {"tv-chan-down", SONY, 12, 0x890,  "dac-vol-down"},
  {"tv-vol-up",    SONY, 12, 0x490},
  {"tv-vol-down",  SONY, 12, 0xC90},
};


void
setup()
{
  Serial.begin(9600);
  Serial.println("setup");

  irrecv.enableIRIn();
}


void
send(const struct Entry &entry) {
  Serial.print("sending ");
  Serial.print(entry.key);
  Serial.print(" 0x");
  Serial.println((unsigned long) entry.value, HEX);

  switch (entry.type) {
    case NEC:
      irsend.sendNEC(entry.value, entry.bits);
      break;
    case SONY:
      irsend.sendSony(entry.value, entry.bits);
      break;
    default:
      Serial.print("unsupported target type: ");
      Serial.println(entry.type);
      break;
  }
}


const Entry *
find_entry_by_name(String name) {
  for (int i = 0; i < sizeof(TABLE) / sizeof(Entry); i++) {
    if (TABLE[i].key == name) {
      return &TABLE[i];
    }
  }
  return 0;
}


const Entry *
find_entry_by_value(decode_type_t type, int bits, unsigned long code) {
  for (int i = 0; i < sizeof(TABLE) / sizeof(Entry); i++) {
    if (TABLE[i].type == type && TABLE[i].bits == bits && TABLE[i].value == code) {
      return &TABLE[i];
    }
  }
  return 0;
}


void
loop() {
  // read commands from serial port
  static String buf = String("");
  while (Serial.available() > 0) {
    char data = Serial.read();
    if (data != '\n' && data != '\0') {
      buf += data;
      continue;
    }

    Serial.print("received serial command: ");
    Serial.println(buf);

    const Entry *code = find_entry_by_name(buf);
    if (code != 0) {
      send(*code);
      setup();
    }

    buf = String("");
  }

  // read ir signals
  decode_results results;
  if (!irrecv.decode(&results)) {
    return;
  }

  if (results.decode_type > 0 && results.bits > 0) {
    Serial.print("received ir code type: ");
    Serial.print(results.decode_type);
    Serial.print(" bits: ");
    Serial.print(results.bits);
    Serial.print(" value: 0x");
    Serial.println(results.value, HEX);

    // transmit trigger code if one exists
    const Entry *entry = find_entry_by_value(results.decode_type, results.bits, results.value);
    if (entry && entry->trigger) {
      const Entry *code = find_entry_by_name(entry->trigger);
      if (code) {
        send(*code);
        delay(40);
        send(*code);
        delay(40);
        send(*code);
        
        setup();
      }
    }
  }

  irrecv.resume();
}
