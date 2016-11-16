# coding:utf-8
import psutil
import threading
import time
import argparse

import subprocess

parser = argparse.ArgumentParser(description="")
parser.add_argument("--sort", help="default: --sort= [mem, cpu]", default='mem')
parser.add_argument("--interval", help="default: --interval=0.5", type=float, default=0.5)

__author__ = 'bary'

lock = threading.RLock()

TMP = {}

threads = []

cpu_total = 0.0

mem_total = 0.0

cpu_count = float(psutil.cpu_count())


class ProcessInfo(psutil.Process):
    def run(self, inter):
        global TMP
        global cpu_total
        global mem_total
        tmp = float(self.memory_percent())
        temp = float(self.cpu_percent(inter))
        lock.acquire()
        cpu_total += temp
        mem_total += tmp
        TMP[self.name()] = {}
        TMP[self.name()]["%cpu"] = temp
        TMP[self.name()]["%mem"] = tmp
        lock.release()


def getallprocess_thread(inter):
    global TMP
    for i in psutil.pids():
        try:
            p = ProcessInfo(i)
        except psutil.NoSuchProcess:
            continue
        else:
            t = threading.Thread(target=p.run, args=(inter,))
            threads.append(t)
            t.start()

    for j in threads:
        if j.is_alive():
            j.join()

    dict = sorted(TMP.iteritems(), key=lambda d: d[1]["%" + args.sort], reverse=True)
    print "[total\n\t%0.4f" % (cpu_total / cpu_count) 
    print "\t%0.4f" % mem_total
    for i in range(len(dict)):
        print "[{0}".format(dict[i][0])
        print "\t%.4f" % (float(dict[i][1]["%cpu"]) / cpu_count)
        print "\t%.4f" % float(dict[i][1]["%mem"])


def getnetio_t():
    tmp = psutil.net_io_counters()
    interval = 0.1
    time.sleep(interval)
    temp = psutil.net_io_counters()
    sum = (temp[0] - tmp[0]) / interval
    print "?flow\n\t%.6f" % (float(sum / 1024.0/ 1024.0 ))
    rec = (temp[1] - tmp[1]) / interval
    print "\t%.6f" % (float(rec / 1024.0/ 1024.0))


def monitor_process():
    p1 = subprocess.Popen(['df', '-h','/tmp'], stdout=subprocess.PIPE)

    lines = p1.stdout.readlines()
    if len(lines) == 0:
        print "/tmp isn't mount"
    else:
        a = [a for a in lines[1].split(" ") if a]

        total = float(a[1][:-1])*1024.0
        used = float(a[1][:-1])/1024.0
        print "?cache\n\t", total
        print "\t", used
        print "\t", a[4][:-1]

if __name__ == "__main__":
    args = parser.parse_args()
    sort_pra = args.sort
    interval = args.interval
    getallprocess_thread(args.interval)

    getnetio_t()
    monitor_process()
