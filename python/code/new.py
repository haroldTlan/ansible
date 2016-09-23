# -*- coding: utf-8 -*-
import json
from o import *
import os,time,sys
from errors import *


def results_select(items, settingtype):
    import pdb
    pdb.set_trace()
    if items.host_failed or items.host_unreachable:
        print errors_info(items, settingtype)
    elif items.host_ok:
        print "?True?%s"%items.host_ok[-1]['result']._result['stdout']
    else:
        print "?False?unknown!!!!!!!!!"
    

def docker_init(master):
    aim = "~/code/yml/vars/server.yml"
    items = run_playbook("yml/docker_master.yml")
    if items.host_failed or items.host_unreachable:
        results_select(items, "master")
    
    elif items.host_ok:

        for i in items.host_ok:
            if i['task'] == 'refresh key':
                keys = i['result']._result['stdout'].split(':\n')[-1]
            elif i['task'] == 'docker swarm init':
                print "?True?%s"%items.host_ok[-1]['result']._result['stdout']
        key = str(keys).replace('\\\n','').strip()
        
        os.system("sed -i 's/key:.*/key: %s/g' %s"% (key, aim))


def server_start(config):
    aim = "~/code/yml/vars/server.yml"

    for i in config:
        if i['order'] == "mysql" and i['ip']:
            os.system("sed -i 's/mysql:.*/mysql: %s/g' %s"% (i['ip'], aim))
            items = run_playbook("yml/mysql.yml")
            results_select(items, i['order'])
            
        elif i['order'] == "mongo" and i['ip']:
            os.system("sed -i 's/mongo:.*/mongo: %s/g' %s"% (i['ip'], aim))
            items = run_playbook("yml/mongo.yml")
            results_select(items, i['order'])

        elif i['order'] == "master" and i['ip']:
            os.system("sed -i 's/master:.*/master: %s/g' %s"% (i['ip'], aim))
            docker_init(i['ip'])

        elif i['order'] == "worker" and i['ip']:
            os.system("sed -i 's/master:.*/master: %s/g' %s"% (i['ip'], aim))
            items = run_playbook("yml/docker_join.yml")
            results_select(items, i['order'])

        elif i['order'] == "service" and i['ip']:
            os.system("sed -i 's/service:.*/service: %s/g' %s"% (i['ip'], aim))
            print "?True?success"

        elif i['order'] == "cloudstor" and i['ip']:
            os.system("sed -i 's/cloudstor:.*/cloudstor: %s/g' %s"% (i['ip'], aim))
            print "?True?success"

        elif i['order'] == "store" and i['ip'] :
            os.system("sed -i 's/store:.*/store: %s/g' %s"% (i['ip'], aim))
            items = run_playbook("yml/store.yml")
            results_select(items, i['order'])

        elif i['order'] == "web" and i['ip'] :
            items = run_playbook("yml/web.yml")
            results_select(items, i['order']) 

        elif i['order'] == "fileserver" and i['ip'] :
            items = run_playbook("yml/fileserver.yml")
            results_select(items, i['order'])

        elif i['order'] == "gateway" and i['ip'] :
            items = run_playbook("yml/gateway.yml")
            results_select(items, i['order'])


if __name__ == '__main__':
    config=[]
    if len(sys.argv) > 2:
        print "enter too much "
    elif len(sys.argv) == 2:
        single = sys.argv[1].split('=')
        config.append(dict(order=single[0], ip=single[1]))
        server_start(config)
    else:
        print "please add ip address"

