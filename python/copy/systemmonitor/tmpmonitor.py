# coding:utf-8
import subprocess
import sys
import time as timme

__author__ = 'bary'


def turnmodtime(modtime):
    tts = timme.localtime(modtime)

    def sttr(x):
        if x < 10:
            return "0" + str(x)
        else:
            return str(x)

    return str(tts.tm_year) + sttr(tts.tm_mon) + sttr(tts.tm_mday) + " " + sttr(tts.tm_hour) + ":" + sttr(
            tts.tm_min) + ":" + sttr(tts.tm_sec)


log = open("/root/tmp.log", "aw")


def monitor_process():
    p1 = subprocess.Popen(['df', '-h', '/tmp'], stdout=subprocess.PIPE)

    lines = p1.stdout.readlines()
    if len(lines) == 0:
        print "/tmp isn't mount"
        return 1
    a = [a for a in lines[1].split(" ") if a]
    print >> log, turnmodtime(timme.time()),
    print turnmodtime(timme.time()),
    print "total size:", a[1],
    print >> log, "total size:", a[1],
    print "used size:", a[2],
    print >> log, "used size:", a[2],
    print "used percentage:", a[4]
    print >> log, "used percentage:", a[4]
    return 0


if __name__ == "__main__":
    monitor_process()
