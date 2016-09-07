# coding:utf-8
import socket
import threading

__author__ = 'bary'


def checkport(i, ip='192.168.2.136', timeout=10):
    tmp = True
    sk = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    sk.settimeout(timeout)
    try:
        sk.connect((ip, i))
    except Exception:
        tmp = False
    sk.close()
    return tmp


if __name__ == "__main__":
    print checkport(37017)
    print checkport(13306)
    print checkport(9000)
    print checkport(9000)
    print checkport(9002)
    print checkport(9333)
