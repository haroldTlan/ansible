# coding:utf-8
import commands

__author__ = 'bary'


def getinfo():
    rc, out = commands.getstatusoutput("netstat -anp | grep weed | grep -v 172")
    listen = []
    established = []
    if rc != 0:
        raise Exception(out)
    else:
        tmp = out.splitlines()
        ttmp = []
        for i in tmp:
            tttmp = []
            temp = i.split(" ")
            for j in temp:
                if j != '':
                    tttmp.append(j)
            ttmp.append(tttmp)
        for i in ttmp:
            if i[5] == "LISTEN":
                listen.append(i[3].split(":")[-1])
            if i[5] == "ESTABLISHED":
                established.append(i[4].split(":")[0])
        print "listenning in port", ''.join(listen)
        print "IP had been established is", ','.join(established)


if __name__ == "__main__":
    getinfo()
