# coding:utf-8
import psutil
import threading
import time
import argparse

parser = argparse.ArgumentParser(description="")
parser.add_argument("--sort", help="default: --sort= [mem, cpu]", default='mem')
parser.add_argument("--interval", help="default: --interval=0.5", type=float, default=0.5)

__author__ = 'bary'

lock = threading.RLock()

TMP = {}

threads = []

cpu_total = 0.0

mem_total = 0.0


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
    print "total %%cpu: %0.4f" % cpu_total
    print "total %%mem: %0.4f" % mem_total
    for i in range(len(dict)):
        print "[{0}]".format(dict[i][0])
        print "\t%.4f" % float(dict[i][1]["%cpu"])
        print "\t%.4f" % float(dict[i][1]["%mem"])


def getnetio_t():
    tmp = psutil.net_io_counters()
    interval = 0.1
    time.sleep(interval)
    temp = psutil.net_io_counters()
    sum = (temp[0] - tmp[0]) / interval
    print "sent: %.3f KB" % (float(sum / 1024.0))
    # print "sent: %.3f MB" % (float(sum / 1024.0 / 1024.0))
    rec = (temp[1] - tmp[1]) / interval
    print "rec: %.3f KB" % (float(rec / 1024.0))
    # print "rec: %.3f MB" % (float(rec / 1024.0 / 1024.0))


if __name__ == "__main__":
    args = parser.parse_args()
    sort_pra = args.sort
    interval = args.interval
    getallprocess_thread(args.interval)

    #getnetio_t()