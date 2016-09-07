# coding:utf-8
from testport import checkport
import argparse

__author__ = 'bary'

parser = argparse.ArgumentParser(description="Check whether the web server port is open")
parser.add_argument("--port", help="default: --port=8081", default=8081, type=int)
parser.add_argument("--ip", help="default: --ip=192.168.2.136", default='192.168.2.136')

if __name__ == "__main__":
    args = parser.parse_args()
    if checkport(args.port, ip=args.ip):
        print 0
    else:
        print 1
