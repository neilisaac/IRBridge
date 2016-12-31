#!/usr/bin/env python2

import logging
import serial
import sys
import tornado.log
import tornado.ioloop
import tornado.web


class MainHandler(tornado.web.RequestHandler):

    def initialize(self, filename):
        self.logger = logging.getLogger(self.__class__.__name__)
        with open(filename) as f:
            self.html = f.read()

    def get(self):
        self.write(self.html)


class SendHandler(tornado.web.RequestHandler):

    def initialize(self, device, baud):
        self.logger = logging.getLogger(self.__class__.__name__)
        self.serial = serial.Serial(device, baud)

    def post(self, code):
        self.logger.info("sending {}".format(code))
        self.serial.write('{}\n'.format(code))
        self.write("sent " + code)


def main():
    tornado.log.enable_pretty_logging()
    app = tornado.web.Application([
        (r"/", MainHandler, dict(filename='remote.html')),
        (r"/send/(.*)", SendHandler, dict(device='/dev/ttyACM0', baud=9600)),
    ])
    app.listen(80)
    tornado.ioloop.IOLoop.current().start()


if __name__ == "__main__":
    main()
