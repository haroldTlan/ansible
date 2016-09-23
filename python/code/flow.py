# -*- coding: utf-8 -*-
import sys
sys.path.append('/root/code')
from o import *


while True:
    items = run_playbook("yml/info/flow.yml")
    for i in items.host_ok:
        print "?flow?" + i['result']._result['stdout']
