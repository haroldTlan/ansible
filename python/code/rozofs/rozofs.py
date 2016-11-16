# -*- coding: utf-8 -*-
import json
from o import *
import os,time,sys
from errors import *
import argparse
import yaml
import commands


parser = argparse.ArgumentParser(description="rozofs")
parser.add_argument("--settingtype", help="default: --settingtype=mysql", default='storage')
parser.add_argument("--ip", help="default: --ip=192.168.2.149", default='192.168.2.149')
parser.add_argument("--slot", help="default: --slot=None", default='None')
parser.add_argument("--expand", help="default: --expand=None", default='None')


def results_select(items, settingtype):
    if items.host_failed or items.host_unreachable:
        print "?False?%s"%errors_info(items, settingtype)
    elif items.host_ok:
        print "?True?%s"%items.host_ok[-1]['result']._result['stdout']
    else:
        print "?False?unknown!!!!!!!!!"
    

def rozofsSet(settingType="storage", ip="192.168.2.190", slot="None", expand="None"):
    aim = "/root/code/rozofs/yml/vars/rozofs.yml"
    yml = "/root/code/rozofs/yml/"
    var = "/root/code/rozofs/yml/vars/rozofs.yml"

    if settingType == "exportInit" and ip:
        os.system("sed -i 's/export:.*/export: %s/g' %s"% (ip, aim))
    	items = run_playbook("%s/exportInit.yml"%yml)
    	results_select(items, settingType)
		
    elif settingType == "storageInit" and ip:
    	os.system("sed -i 's/storage:.*/storage: %s/g' %s"% (ip, aim))
    	items = run_playbook("%s/storageInit.yml"%yml)
    	results_select(items, settingType)

    elif settingType == "storage" and ip:
        storage = ' '.join(ip.split(","))
        storages = ip.split(",")
        slots = slot.split(",")

    	os.system("sed -i 's/storage:.*/storage: %s/g' %s"% (storage, aim))
    	os.system("sed -i 's/storages:.*/storages: %s/g' %s"% (storages, aim))
    	os.system("sed -i 's/slot:.*/slot: %s/g' %s"% (str(slots).replace("'","\""), aim))

    	items = run_playbook("%s/storage.yml"%yml)
    	results_select(items, settingType)

    elif settingType == "client" and ip:
        os.system("sed -i 's/client:.*/client: %s/g' %s"% (ip, aim))
        items = run_playbook("%s/client.yml"%yml)
        results_select(items, settingType)

    elif settingType == "export" and ip:
        expands = ' '.join(expand.split(","))
        results = commands.getoutput("rozo volume expand %s --vid 1 --layout 0  --exportd %s"%(expands,ip))
        results = commands.getoutput("rozo node start -E %s"%ip)
        if "FAILED" in  results :
            result = results
            print "?False?%s"%result
        else:
            print "?True?"

    	
    elif settingType == "prestorage" and ip:
        
        check = ip.split(",")
        if True:
            print "?False?success"

        else:
            if len(check) %4 ==0:
                os.system("sed -i 's/prestorage:.*/prestorage: %s/g' %s"% (ip, aim))
    	        items = run_playbook("yml/prestorage.yml")
    	        results_select(items, settingType)
            else:
                print "?False?need multiples of 4"

if __name__ == '__main__':
    args = parser.parse_args()
    rozofsSet(args.settingtype, args.ip, args.slot,args.expand)


