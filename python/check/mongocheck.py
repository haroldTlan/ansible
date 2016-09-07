# coding:utf-8
from testport import checkport
import pymongo
import argparse

__author__ = 'bary'

parser = argparse.ArgumentParser(description="Check whether the Mongo server port is open")
parser.add_argument("--port", help="default: --port=37017", default=37017, type=int)
parser.add_argument("--ip", help="default: --ip=192.168.2.136", default='192.168.2.136')


def checkmongo(ip="192.168.2.136", port=37017):
    if not checkport(port, ip=ip):
        return 1
    try:
        client = pymongo.MongoClient(ip, port)
        client.database_names()
        return 0
    except Exception:
        return 2


if __name__ == "__main__":
    args = parser.parse_args()
    print checkmongo(args.ip, args.port)
