# -*- coding: utf-8 -*-

import sys
sys.path.append('/root/code')
from o import *
import time

#import pdb
#pdb.set_trace()

def getInfo(items):
    info = []
    if len(items.host_ok)>0:
        for i in items.host_ok:
            info.append(dict(type=str(i['task']), ip=str(i['ip']), result=str(i['result']._result['stdout']),status="success"))

    if len(items.host_unreachable) >0:
        for i in items.host_unreachable:
            info.append(dict(type=str(i['task']), ip=str(i['ip']), status="unreachable"))

    if len(items.host_failed) >0:
        for i in items.host_failed:
            info.append(dict(type=str(i['task']), ip=str(i['ip']), status="failed"))

    print info
#if "no such file" in result:
#    err = "aim files did not exist"

if __name__ == '__main__':
    items = run_playbook("/root/code/yml/info/process.yml")
    getInfo(items)
