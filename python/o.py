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

loader = DataLoader()
variable_manager = VariableManager()
inventory = Inventory(loader=loader, variable_manager=variable_manager)
variable_manager.set_inventory(inventory)


class ResultsCollector(CallbackBase):
    def __init__(self, *args, **kwargs):
        super(ResultsCollector, self).__init__(*args, **kwargs)
        self.host_ok = []
        self.host_unreachable = {}
        self.host_failed = {}

    def v2_runner_on_unreachable(self, result):
        self.host_unreachable[result._host.get_name()] = result

    def v2_runner_on_ok(self, result,  *args, **kwargs):
        name = result._host.get_name()
        task = result._task.get_name()
        if task =='setup':
            pass
        else:
            self.host_ok.append(dict(ip=name, task=task, result=result))

    def v2_runner_on_failed(self, result,  *args, **kwargs):
        self.host_failed[result._host.get_name()] = result

class Options(object):
    def __init__(self):
        self.connection = "smart"
        self.forks = 10
        self.check = False
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
        #import pdb
        #pdb.set_trace()
        #for i in callback.host_ok:
           # print i['ip'], i['task'], i['result']._result['stdout']
            #return i['ip'], i['task'], i['result']._result['stdout']
    finally:
        if tqm is not None:
            tqm.cleanup()

def run_playbook(books):
    results_callback = callback_loader.get('json')
    #playbooks=['yml/docker.yml']
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
        #import pdb
        #pdb.set_trace()
        #for i in callback.host_ok:
        #    print i['ip'], i['task'], i['result']._result['stdout']
        return callback
        #for host, results in callback.host_ok.items():
            #device = trans(result._result['stdout'])
            #print host, results._result['out.stdout_lines']

    except Exception as e:
        print ('run_playbook:%s'%e)

if __name__ == '__main__':
    #run_playbook("yml/store.yml")
    order= "docker swarm join     --token SWMTKN-1-2iz0i3evtuous8rksj5mc9uuhs0ytwdcnkke6407dmpl69187a-8svnkz3dqykisk14ust2n0ku4     192.168.2.148:2377"
    run_adhoc("192.168.2.149", order)
    #docker_init()
