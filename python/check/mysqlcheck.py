# coding:utf-8
from testport import checkport
from peewee import *
import argparse

__author__ = 'bary'

parser = argparse.ArgumentParser(description="Check whether the Mongo server port is open and check the tables")
parser.add_argument("--port", help="default: --port=13306", default=13306, type=int)
parser.add_argument("--ip", help="default: --ip=192.168.2.136", default='192.168.2.136')


def checkmysql(ip="192.168.2.136", port=13306):
    if not checkport(port, ip=ip):
        return "port not open"
        #return 1
    try:
        db = MySQLDatabase('cloud', autocommit=True, user='root', passwd='passwd', threadlocals=True,
                           connect_timeout=300, host=ip
                           , port=port)
        db.connect()
        if db.get_tables().count("ipc"):
            if db.get_tables().count("user"):
                pass
            else:
                return "table named user don't exist"
                #return 3
        else:
            return "table named ipc don't exist"
            #return 4
        db.close()
        return "success"
        #return 0
    except Exception as e:
        try:
            db.close()
        except Exception:
            pass
        print e
        return "can't connect to database cloud"
        #return 2


if __name__ == "__main__":
    args = parser.parse_args()
    print checkmysql(args.ip, args.port)
