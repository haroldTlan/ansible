# coding:utf-8
import subprocess
import signal
import re
import argparse
import time

__author__ = 'bary'

parser = argparse.ArgumentParser(description="Check whether the gw port is open")
parser.add_argument("--port", help="default: --port=9000", default=9000, type=int)
parser.add_argument("--ip", help="default: --ip=192.168.2.216", default='192.168.2.136')


def run_command(cmd, shell=False, stdout=subprocess.PIPE,
                stderr=subprocess.PIPE, stdin=subprocess.PIPE, throw=True,
                log=False, input=None, needlog=True):
    cmd = map(str, cmd)
    p = subprocess.Popen(cmd, shell=shell, stdout=stdout, stderr=stderr,
                         stdin=stdin)
    out, err = p.communicate(input=input)
    out = out.split('\n')
    err = err.split('\n')
    rc = p.returncode
    return out, err, rc


def getaclinfo():
    out, err, rc = run_command(
            ['curl', '-d', " ", ''.join(["http://", ''.join([options.ip, ":", str(options.port)]), "/"])])
    try:
        return re.findall(r"<AccessKeyID>.+</AccessKeyID>", out[1])[0].lstrip(r"<AccessKeyID>").rstrip(
            r"</AccessKeyID>"), \
               re.findall(r"<SecretAccessKey>.+</SecretAccessKey>",
                          out[1])[0].lstrip(r"<SecretAccessKey>").rstrip(r"</SecretAccessKey>")
    except Exception:
        print "error"
        time.sleep(4)


def handler(signum, frame):
    raise AssertionError


if __name__ == '__main__':
    options = parser.parse_args()
    try:
        signal.signal(signal.SIGALRM, handler)
        signal.alarm(3)
        print getaclinfo()
        signal.alarm(0)
        #print 0
        print "success"
    except AssertionError:
        #print 1
        print "port not open"
