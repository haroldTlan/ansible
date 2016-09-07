# -*- coding: utf-8 -*-
import json
import yaml
from o import *
import os,time,sys


def docker_init(master):
    aim = "~/code/yml/vars/server.yml"
    items = run_playbook("yml/docker_master.yml")
    
    if items.host_ok:

        for i in items.host_ok:
            if i['task'] == 'refresh key':
                keys = i['result']._result['stdout'].split(':\n')[-1]
        key = str(keys).replace('\\\n','').strip()
        
        os.system("sed -i 's/key:.*/key: %s/g' %s"% (key, aim))
#    elif items.host_failed:
#        print items.host_failed.values()[0]._result['stderr']

def server_start(config):
    aim = "~/code/yml/vars/server.yml"

    for i in config:
        if i['order'] == "mysql" and len(i['ip']) > 0:
            os.system("sed -i 's/mysql:.*/mysql: %s/g' %s"% (i['ip'], aim))
            items = run_playbook("yml/mysql.yml")

            if items.host_ok:
                print "True"
            elif items.host_failed:
                print items.host_failed.values()[0]._result['stderr']
            
        elif i['order'] == "mongo" and len(i['ip']) > 0:
            os.system("sed -i 's/mongo:.*/mongo: %s/g' %s"% (i['ip'], aim))
            items = run_playbook("yml/mongo.yml")

            if items.host_ok:
                print "True"
            elif items.host_failed:
                print items.host_failed.values()[0]._result['stderr']

        elif i['order'] == "master" and i['ip']:
            os.system("sed -i 's/master:.*/master: %s/g' %s"% (i['ip'], aim))
            docker_init(i['ip'])
            print "True"

        elif i['order'] == "worker" and i['ip']:
            os.system("sed -i 's/master:.*/master: %s/g' %s"% (i['ip'], aim))
            print i['ip']

        elif i['order'] == "service" and i['ip']:
            os.system("sed -i 's/service:.*/service: %s/g' %s"% (i['ip'], aim))
            print "True"
        #   items = run_playbook("yml/docker_join.yml")
            #key = "docker swarm join     --token SWMTKN-1-51qnav9l692xjojl07pdmt2nfx324ivy8r1c9yfb2vlkkryr3g-5ednjclettpdr7qb25wpkzhb7     192.168.2.225:2377"
            #items = run_adhoc("192.168.2.226", key)
        #    if items.host_ok:
        #        print "mysql success"
        #    elif items.host_failed:
        #        print items.host_failed.values()[0]._result['stderr']

        elif i['order'] == "cloudstor" and i['ip']:
            os.system("sed -i 's/cloudstor:.*/cloudstor: %s/g' %s"% (i['ip'], aim))
            print "True"

        elif i['order'] == "store" and i['ip'] :
            os.system("sed -i 's/store:.*/store: %s/g' %s"% (i['ip'], aim))
            items = run_playbook("yml/store.yml")
            if items.host_ok:
                print "True"
            else:
                print items.host_failed.values()[0]._result['stderr']

        elif i['order'] == "web" and i['ip'] :
            items = run_playbook("yml/web.yml")
            if items.host_ok:
                print "True"
            else:
                print items.host_failed.values()[0]._result['stderr']

        elif i['order'] == "fileserver" and i['ip'] :
            items = run_playbook("yml/fileserver.yml")
            if items.host_ok:
                print "True"
            else:
                print items.host_failed.values()[0]._result['stderr']

        elif i['order'] == "gateway" and i['ip'] :
            items = run_playbook("yml/gateway.yml")
            if items.host_ok:
                print "True"
            else:
                print items.host_failed.values()[0]._result['stderr']


if __name__ == '__main__':
    config=[]
    if len(sys.argv) > 2:
        print "enter too much "
    elif len(sys.argv) == 2:
        strings = sys.argv[1].split(',')
        for singles in strings:
            single = singles.split('=')
            config.append(dict(order=single[0], ip=single[1]))
        server_start(config)
    else:
        print "please add ip address"

