# -*- coding: utf-8 -*-

import sys
sys.path.append('/root/code')
from o import *
import time
items = run_playbook("/root/code/yml/info/process.yml")

try:
    for i in items.host_ok:
        print "?info?" + i['result']._result['stdout']


except:
    print "error"
