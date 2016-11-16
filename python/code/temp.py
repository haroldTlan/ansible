# -*- coding: utf-8 -*-
import json
from o import *
import os,time,sys
from errors import *
import argparse
import yaml


parser = argparse.ArgumentParser(description="rozofs")
parser.add_argument("--settingtype", help="default: --settingtype=mysql", default='storage')
parser.add_argument("--ip", help="default: --ip=192.168.2.149", default='192.168.2.149')
parser.add_argument("--expand", help="default: --expand=None", default='None')



def results_select(items, settingtype):
    if items.host_failed or items.host_unreachable:
        print "?False?%s"%errors_info(items, settingtype)
    elif items.host_ok:
        print "?True?%s"%items.host_ok[-1]['result']._result['stdout']
    else:
        print "?False?unknown!!!!!!!!!"
    
def rozofsset(settingType="storage", ip="192.168.2.190", expand="None"):
    print settingType,ip

def rozofsSet(settingType="storage", ip="192.168.2.190", expand="None"):
    aim = "~/code/yml/vars/rozofs.yml"

    if settingType == "exportInit" and ip:
        os.system("sed -i 's/export:.*/export: %s/g' %s"% (ip, aim))
    	items = run_playbook("yml/rozofs/exportInit.yml")
    	results_select(items, settingType)
		
    elif settingType == "storageInit" and ip:
    	os.system("sed -i 's/storage:.*/storage: %s/g' %s"% (ip, aim))
    	items = run_playbook("yml/rozofs/storageInit.yml")
    	results_select(items, settingType)

    elif settingType == "export" and ip:
        with open('/root/code/yml/vars/rozofs.yml', 'r') as conf:
            expands = yaml.load(conf)['storage']
            expand = ' '.join(expands.split(","))

        if len(expand.split()) %4 ==0:
            os.system("sed -i 's/export:.*/export: %s/g' %s"% (ip, aim))
    	    os.system("sed -i 's/expand:.*/expand: %s/g' %s"% (expand, aim))
    	    items = run_playbook("yml/rozofs/export.yml")
            results_select(items, settingType)

        else:
            print "?False?need multiples of 4"
    	
    elif settingType == "prestorage" and ip:
        
        check = ip.split(",")
        if True:
            print "?False?success"

        else:
            if len(check) %4 ==0:
                os.system("sed -i 's/prestorage:.*/prestorage: %s/g' %s"% (ip, aim))
    	        items = run_playbook("yml/rozofs/prestorage.yml")
    	        results_select(items, settingType)
            else:
                print "?False?need multiples of 4"

    elif settingType == "storage" and ip:
    	os.system("sed -i 's/storage:.*/storage: %s/g' %s"% (ip, aim))
    	items = run_playbook("yml/rozofs/storage.yml")
    	results_select(items, settingType)

    elif settingType == "client" and ip:
    	os.system("sed -i 's/client:.*/client: %s/g' %s"% (ip, aim))
    	items = run_playbook("yml/rozofs/client.yml")
    	results_select(items, settingType)

    elif settingType == "temp" and ip:
        check = ip.split(",")
        print check
        if len(check) %4 ==0:
            for i in range(len(check)):
       	        os.system("sed -i 's/single:.*/single: %s/g' %s"% (check[i], aim))
    	        os.system("sed -i 's/temp:.*/temp: %s/g' %s"% (i, aim))
    	        items = run_playbook("yml/rozofs/temp.yml")
    	        results_select(items, settingType)
        else:
            print "?False?need multiples of 4"




if __name__ == '__main__':
    args = parser.parse_args()
    rozofsset(args.settingtype,args.ip)


