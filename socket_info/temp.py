# -*- coding: utf-8 -*-

import sys
sys.path.append('/root/code')
from o import *

import json
import os
#import pdb
#pdb.set_trace()

def getInfo(items):
    info = []
    if len(items.host_ok)>0:
        for i in items.host_ok:
            a = i['result']._result['stdout']
            try:
                infos = infoTrans(i['task'], str(i['result']._result['stdout']))
                info.append(dict(type=str(i['task']), ip=str(i['ip']), result=infos, status="success"))
            except:
                pass

    if len(items.host_unreachable) >0:
        for i in items.host_unreachable:
            info.append(dict(type=str(i['task']), ip=str(i['ip']), status="unreachable"))

    if len(items.host_failed) >0:
        for i in items.host_failed:
            info.append(dict(type=str(i['task']), ip=str(i['ip']), status="failed"))
    infos = json.dumps(info)
    os.system("echo '%s' > static"%infos)
    #print info
#if "no such file" in result:
#    err = "aim files did not exist"

def infoTrans(dev, infos):
#    if "master" in dev:
#        return json.loads(infos)
#    elif "store" in dev:
        a = json.loads(infos)
        return a["samples"]


if __name__ == '__main__':
    items = run_playbook("/root/code/yml/info/store.yml")
    getInfo(items)
