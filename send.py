#!/usr/bin/env python2

import serial
import sys

s = serial.Serial('/dev/ttyACM0', 9600)

for a in sys.argv:
    s.write('{}\n'.format(a))

