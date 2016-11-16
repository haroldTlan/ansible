# -*- coding: utf-8 -*-


def errors_info(items, settingtype):
    if items.host_failed:
        try:
            stderr = items.host_failed[0]["result"]._result['stderr']

        except:
            stderr = items.host_failed[0]["result"]._result['msg']
    else:
        try:
            stderr = items.host_unreachable[0]["result"]._result['msg']

        except:
            stderr = items.host_unreachable[0]["result"]._result['stderr']


    try:
        if "already" in stderr:
            return "%s has been created"%settingtype
        elif "No such file" in stderr:
            return "Need docker swarm init"
        elif "Exception" in stderr:
             return stderr.split("Exception: ")[-1]
        else:
            return "%s"%stderr
    except:
        return "this is a big bug for stderr"

def success_info(items, settingtype):
    
    try:
        stdout = items.host_ok.values()[0]._result['stdout']
    except:
        pass

    try:
        if "To add a worker to this swarm" in stdout:
            return "docker init success"
        elif "No such file" in stderr:
            return "Need docker swarm init"
        else:
            return "%s"%stderr
    except:
        return "this is a big bug for stderr"


