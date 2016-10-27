# coding:utf-8
import subprocess
import sys

__author__ = 'bary'


def monitor_process(key_word):
    p1 = subprocess.Popen(['ps', '-ef'], stdout=subprocess.PIPE)
    p2 = subprocess.Popen(['grep', key_word], stdin=p1.stdout, stdout=subprocess.PIPE)
    p3 = subprocess.Popen(['grep', '-v', 'grep'], stdin=p2.stdout, stdout=subprocess.PIPE)
    p4 = subprocess.Popen(['grep', '-v', 'python'], stdin=p3.stdout, stdout=subprocess.PIPE)

    lines = p4.stdout.readlines()
    add = ""
    for i in lines:
        if "master" in i:
            add = "master "
        elif "volume" in i:
            add = "volume"
        elif "master" and "volume" in i:
            add = "master and volume"

    if len(lines) == 0:
        print "the", key_word, "isn't running"
        return 1
    print "the", key_word, "(", add, ")" "is running"
    return 0


if __name__ == "__main__":
    monitor_process("weed")
