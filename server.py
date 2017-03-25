#!/usr/bin/env python2

import argparse
import logging
import requests
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

    def post(self, code, count=1):
        self.logger.info("sending {} x{}".format(code, count))
        for _ in range(int(count)):
            self.serial.write('{}\n'.format(code))
        self.write("sent " + code)


class IFTTTHandler(tornado.web.RequestHandler):

    def initialize(self, apikey):
        self.logger = logging.getLogger(self.__class__.__name__)
        self.apikey = apikey

    def post(self, event):
        self.logger.info("sending {}".format(event))
        if self.apikey:
            requests.post('https://maker.ifttt.com/trigger/{}/with/key/{}'.format(event, self.apikey))
            self.write("triggered " + event)
        else:
            self.write("ignored event because api key was not provided " + event)


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument('--apikey', default=None)
    parser.add_argument('--device', default='/dev/ttyACM0')
    parser.add_argument('--baud', default=9600)
    parser.add_argument('--port', default=80)
    args = parser.parse_args()

    tornado.log.enable_pretty_logging()
    app = tornado.web.Application([
        (r"/", FileHandler, dict(filename='remote.html')),
        (r"/icon.png", FileHandler, dict(filename='icon.png', content_type='image/png')),
        (r"/favicon.png", FileHandler, dict(filename='favicon.png', content_type='image/png')),
        (r"/send/(.*)/(.*)", SendHandler, dict(device=args.device, baud=args.baud)),
        (r"/send/(.*)", SendHandler, dict(device=args.device, baud=args.baud)),
        (r"/ifttt/(.*)", IFTTTHandler, dict(apikey=args.apikey)),
    ])

    app.listen(args.port)
    tornado.ioloop.IOLoop.current().start()


if __name__ == "__main__":
    main()
