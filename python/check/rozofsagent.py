# coding:utf-8
from testport import checkport
import argparse

__author__ = 'bary'

parser = argparse.ArgumentParser(description="Check whether the agent is running")
parser.add_argument("--port", help="default: --port=10050", default=10050, type=int)
parser.add_argument("--ip", help="default: --ip=192.168.2.180", default='192.168.2.180')

if __name__ == "__main__":
    args = parser.parse_args()
    if checkport(args.port, ip=args.ip):
        print 0
    else:
        print 1