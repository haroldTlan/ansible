# -*- coding: utf-8 -*-

import sys
sys.path.append('/root/code')
from o import *
import time
'''
while True:
    items = run_playbook("yml/info/process.yml")

    for i in items.host_ok:
        if i['task'] == "process":
        
            print "?process?" + i['result']._result['stdout']

        elif i['task'] == "networkflow":
            pass
#        print "?flow?" + i['result']._result['stdout_lines']
    #time.sleep(2)
'''
#import pdb
#pdb.set_trace()
items = run_playbook("/root/code/yml/info/flow.yml")

try:
    for i in items.host_ok:
#        if i['task'] == "process":
#            print "?info?" + i['result']._result['stdout']
        print "?info?" + i['result']._result['stdout']
#        elif i['task'] == "networkflow":
#            print "?flow?" + i['result']._result['stdout']


except:
    print "error"
