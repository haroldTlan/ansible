import sys
import argparse
sys.path.append('/root/check')
from testport import checkport
from mongocheck import *
from mysqlcheck import *
from gatewaycheck import *
from fileservercheck import *
from o import *
#from webcheck import *

parser = argparse.ArgumentParser(description="use bary's check module")
parser.add_argument("--ip", help="default: --ip=192.168.2.149", default='192.168.2.149')
parser.add_argument("--checktype", help="default: --checktype=mysql", default='mysql')

aim = "~/code/yml/vars/server.yml"

def checkmodule(ip="192.168.2.149", checktype="mysql"):

    if checktype == "mysql":
        if checkmysql(ip) == "success" :
            results = "?True?%s"%checkmysql(ip)
        else:
            results = "?False?%s"%checkmysql(ip)

    elif checktype == "mongo":
        if checkmongo(ip) == "success":
            results = "?True?%s"%checkmongo(ip)
        else:
            results = "?False?%s"%checkmongo(ip)


    elif checktype == "gateway":
        if checkport(9000, ip=args.ip):
            results = "?True?success"
        else:
            results = "?False?port not open"
    
    elif checktype == "beanstalkd":
        if checkport(11300, ip=args.ip):

            results = "?True?success"
        else:
            results = "?False?port not open"

    elif checktype == "fileserver":
        if checkport(9002, ip=args.ip):
            results = "?True?success"
        else:
            results = "?False?port not open"
        
    elif checktype == "web":
        if checkport(8081, ip=args.ip):
            results = "?True?success"
        else:
            results = "?False?port not open"

    elif checktype == "node":
        os.system("sed -i 's/master:.*/master: %s/g' %s"% (ip, aim))
        items = run_playbook("yml/check/node.yml")
        if items.host_ok:
            for i in items.host_ok:
                if i['task'] == 'checknode':
                    result = i['result']._result['stdout_lines'][0]
            
            if "." in result:
                results = "?True?" + result
            else:
                results = "?False?" + result
        else:
            results = "?False?port not open"

    else:
        results = "?False?No such type"

    
    return results
    
if __name__ == '__main__':
    args = parser.parse_args()
    print checkmodule(args.ip, args.checktype)
