# -*- coding: utf-8 -*-
import json
from ansible.parsing.dataloader import DataLoader
from ansible.vars import VariableManager
from ansible.inventory import Inventory
from ansible.playbook.play import Play
from ansible.executor.task_queue_manager import TaskQueueManager
from ansible.executor.playbook_executor import PlaybookExecutor

from ansible.plugins import callback_loader
from ansible.plugins.callback import CallbackBase

import os
import logging

import shutil
import ansible.constants as C


loader = DataLoader()
variable_manager = VariableManager()
inventory = Inventory(loader=loader, variable_manager=variable_manager)
variable_manager.set_inventory(inventory)


FIELDS = ['cmd', 'command', 'start', 'end', 'delta', 'msg', 'stdout', 'stderr']
STD = ['stdout', 'stderr', 'msg']


logger = logging.getLogger('ansible')
fh = logging.FileHandler('/var/log/ansible.log')
logger.setLevel(logging.INFO)
formatter = logging.Formatter('%(asctime)s - %(name)s - %(message)s')
fh.setFormatter(formatter)
logger.addHandler(fh)
    
#push log to /var/log/ansible.log
def ansible_log(res_obj):
    res=res_obj._result
    
    info = ""
    if type(res) == type(dict()):
        for field in STD:
            if field in res.keys() and len(res[field]) > 0:
                info = u'{0}:{1}'.format(field, res[field])
    log_info = u'{0}: {1} {2}'.format(res_obj._host, res_obj._task, info)
    logger.info(log_info)

#get result output
class ResultsCollector(CallbackBase):
    def __init__(self, *args, **kwargs):
        super(ResultsCollector, self).__init__(*args, **kwargs)
        self.host_ok = []
        self.host_unreachable = []
        self.host_failed = []

    def v2_runner_on_unreachable(self, result, ignore_errors=False):
        name = result._host.get_name()
        task = result._task.get_name()
        ansible_log(result)
        #self.host_unreachable[result._host.get_name()] = result
        self.host_unreachable.append(dict(ip=name, task=task, result=result))

    def v2_runner_on_ok(self, result,  *args, **kwargs):
        name = result._host.get_name()
        task = result._task.get_name()
        if task == "setup":
            pass
        elif "Info" in task:
            self.host_ok.append(dict(ip=name, task=task, result=result))
        else:
            ansible_log(result)
            self.host_ok.append(dict(ip=name, task=task, result=result))

    def v2_runner_on_failed(self, result,   *args, **kwargs):
        name = result._host.get_name()
        task = result._task.get_name()
        ansible_log(result)
        #self.host_failed[result._host.get_name()] = result
        self.host_failed.append(dict(ip=name, task=task, result=result))

class Options(object):
    def __init__(self):
        self.connection = "ssh"
        self.forks = 10
        self.check = False
        self.become = None
        self.become_method = None
        self.become_user=None
        module_path='/usr/local/lib/python2.7/dist-packages/ansible-2.2.0-py2.7.egg/ansible/modules'
    def __getattr__(self, name):
        return None

options = Options()

def trans(strings):
    string = strings.replace('\'','\"')
    return json.loads(string)

def run_adhoc(ip,order):
    variable_manager.extra_vars={"ansible_ssh_user":"root" , "ansible_ssh_pass":"passwd"}
    play_source = {"name":"Ansible Ad-Hoc","hosts":"%s"%ip,"gather_facts":"no","tasks":[{"action":{"module":"command","args":"%s"%order}}]}
#    play_source = {"name":"Ansible Ad-Hoc","hosts":"192.168.2.160","gather_facts":"no","tasks":[{"action":{"module":"command","args":"python ~/store.py del"}}]}   
    play = Play().load(play_source, variable_manager=variable_manager, loader=loader)
    tqm = None
    callback = ResultsCollector()

    try:
        tqm = TaskQueueManager(
            inventory=inventory,
            variable_manager=variable_manager,
            loader=loader,
            options=options,
            passwords=None,
            #stdout_callback='minimal',
            #stdout_callback=results_callback,
            run_tree=False,
        )
        tqm._stdout_callback = callback
        result = tqm.run(play)
        return callback

    finally:
        if tqm is not None:
            tqm.cleanup()

def run_playbook(books):
    results_callback = callback_loader.get('json')
    playbooks = [books]

    variable_manager.extra_vars={"ansible_ssh_user":"root" , "ansible_ssh_pass":"passwd"}
    callback = ResultsCollector()

    pd = PlaybookExecutor(
        playbooks=playbooks,
        inventory=inventory,
        variable_manager=variable_manager,
        loader=loader,
        options=options,
        passwords=None,

        )
    #pd._tqm._stdout_callback = results_callback
    pd._tqm._stdout_callback = callback
   
    try:
        result = pd.run()
        shutil.rmtree(C.DEFAULT_LOCAL_TMP, True)
        return callback

    except Exception as e:
        # ('run_playbook:%s'%e)
        print "error"
        print e

if __name__ == '__main__':
    run_playbook("yml/info/store.yml")
    #order= "docker swarm join     --token SWMTKN-1-2iz0i3evtuous8rksj5mc9uuhs0ytwdcnkke6407dmpl69187a-8svnkz3dqykisk14ust2n0ku4     192.168.2.148:2377"
    #run_adhoc("192.168.2.149", "ifconfig")
    #docker_init()
