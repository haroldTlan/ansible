import sys
import argparse
import commands
sys.path.append('/root/check')
from testport import checkport
from mongocheck import *
from mysqlcheck import *
from gatewaycheck import *
from fileservercheck import *
#from webcheck import *

parser = argparse.ArgumentParser(description="use bary's check module")
parser.add_argument("--ip", help="default: --ip=192.168.2.149", default='192.168.2.149')
parser.add_argument("--checktype", help="default: --checktype=mysql", default='mysql')


def checkmodule(ip="192.168.2.149", checktype="mysql"):
    
    if checktype == "mysql":
        if checkmysql(ip) :
            return "failed:%d"%checkmysql(ip)
        else:
            return "True"

    elif checktype == "mongo":
        if checkmongo(ip) :
            return "failed:%d"%checkmongo(ip)
        else:
            return "True"


    elif checktype == "gateway":
        result =  commands.getoutput("python /root/check/gatewaycheck.py --ip=%s"%ip)
        if result:
            return "True"
        else:
            return "failed:%d"%result
    

    elif checktype == "fileserver":
        result = commands.getoutput("python /root/check/fileservercheck.py --ip=%s"%ip)
        if result:
            return "True"
        else:
            return "failed:%s"%result


    elif checktype == "web":
        if checkport(8081, ip=args.ip):
            return "True"
        else:
            return "failed"

if __name__ == '__main__':
    args = parser.parse_args()
    print checkmodule(args.ip, args.checktype)
