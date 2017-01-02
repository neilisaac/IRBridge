#!/usr/bin/env python2

import logging
import serial
import sys
import tornado.log
import tornado.ioloop
import tornado.web


class FileHandler(tornado.web.RequestHandler):

    def initialize(self, filename, content_type='text/html'):
        self.logger = logging.getLogger(self.__class__.__name__)
        self.content_type = content_type
        with open(filename) as f:
            self.data = f.read()

    def get(self):
        self.set_header("Content-type", self.content_type)
        self.write(self.data)


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
        (r"/", FileHandler, dict(filename='remote.html')),
        (r"/icon.png", FileHandler, dict(filename='icon.png', content_type='image/png')),
        (r"/favicon.png", FileHandler, dict(filename='favicon.png', content_type='image/png')),
        (r"/send/(.*)", SendHandler, dict(device='/dev/ttyACM0', baud=9600)),
    ])
    app.listen(80)
    tornado.ioloop.IOLoop.current().start()


if __name__ == "__main__":
    main()
