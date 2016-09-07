import sys
import argparse
import commands
from o import *

parser = argparse.ArgumentParser(description="use bary's check module")
parser.add_argument("--stoptype", help="default: --stoptype=mysql", default='mysql')
parser.add_argument("--ip", help="default: --ip=192.168.2.149", default='192.168.2.149')

aim = "~/code/yml/vars/server.yml"

def stopmodule(stoptype="mysql", ip="192.168.2.149"):
    irrelevant_service = ["cloudstor", "service", "store"]  
    if stoptype == "mysql":
        os.system("sed -i 's/mysql:.*/mysql: %s/g' %s"% (ip, aim))
        items = run_playbook("yml/stop/mysql.yml")
        if items.host_ok:
            print "True"
        else:
            print items.host_failed.values()[0]._result['stderr']


    elif stoptype == "mongo":
        os.system("sed -i 's/mongo:.*/mongo: %s/g' %s"% (ip, aim))
        items = run_playbook("yml/stop/mongo.yml")
        if items.host_ok:
            print "True"
        else:
            print items.host_failed.values()[0]._result['stderr']

    elif stoptype == "master":
        os.system("sed -i 's/master:.*/master: %s/g' %s"% (ip, aim))
        items = run_playbook("yml/stop/master.yml")
        if items.host_ok:
            print "True"
        else:
            print items.host_failed.values()[0]._result['stderr']


    elif stoptype == "gateway":
        os.system("sed -i 's/master:.*/master: %s/g' %s"% (ip, aim))
        items = run_playbook("yml/stop/gateway.yml")
        if items.host_ok:
            print "True"
        else:
            print items.host_failed.values()[0]._result['stderr']


    elif stoptype == "fileserver":
        os.system("sed -i 's/master:.*/master: %s/g' %s"% (ip, aim))
        items = run_playbook("yml/stop/fileserver.yml")
        if items.host_ok:
            print "True"
        else:
            print items.host_failed.values()[0]._result['stderr']

    elif stoptype == "web":
        os.system("sed -i 's/master:.*/master: %s/g' %s"% (ip, aim))
        items = run_playbook("yml/stop/web.yml")
        if items.host_ok:
            print "True"
        else:
            print items.host_failed.values()[0]._result['stderr']

    elif stoptype in irrelevant_service:
        print "True"
    else:
        print "Please enter the right value"

if __name__ == '__main__':
    args = parser.parse_args()
    stopmodule(args.stoptype, args.ip)
