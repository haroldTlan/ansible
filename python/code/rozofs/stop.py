# -*- coding: utf-8 -*-
import json
from o import *
import os,time,sys
from errors import *
import argparse
import yaml


parser = argparse.ArgumentParser(description="rozofs")
parser.add_argument("--stoptype", help="default: --stoptype=mysql", default='storage')
parser.add_argument("--ip", help="default: --ip=192.168.2.149", default='192.168.2.149')



def results_select(items, settingtype):
    if items.host_failed or items.host_unreachable:
        print "?False?%s"%errors_info(items, settingtype)
    elif items.host_ok:
        print "?True?%s"%items.host_ok[-1]['result']._result['stdout']
    else:
        print "?False?unknown!!!!!!!!!"
    

def rozofsSet(settingType="storage", ip="192.168.2.190"):
    aim = "~/code/rozofs/yml/vars/rozofs.yml"
    yml = "/root/code/rozofs/yml/stop"

    if settingType == "storage" and ip:
        #each = ip.split(",")
        os.system("sed -i 's/serviceStop:.*/serviceStop: %s/g' %s"% (ip, aim))
    	items = run_playbook("%s/storage.yml"%yml)
    	results_select(items, settingType)

    elif settingType == "client" and ip:
        os.system("sed -i 's/serviceStop:.*/serviceStop: %s/g' %s"% (ip, aim))
    	items = run_playbook("%s/client.yml"%yml)
    	results_select(items, settingType)

    elif settingType == "export" and ip:
        os.system("sed -i 's/serviceStop:.*/serviceStop: %s/g' %s"% (ip, aim))
        items = run_playbook("%s/export.yml"%yml)
        results_select(items, settingType)
    	
if __name__ == '__main__':
    args = parser.parse_args()
    rozofsSet(args.stoptype, args.ip)


