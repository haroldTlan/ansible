# -*- coding: utf-8 -*-


def errors_info(items, settingtype):
    if items.host_failed:
        try:
            stderr = items.host_failed.values()[0]._result['stderr']

        except:
            stderr = items.host_failed.values()[0]._result['msg']
    else:
        try:
            stderr = items.host_unreachable.values()[0]._result['msg']

        except:
            stderr = items.host_unreachable.values()[0]._result['stderr']


    try:
        if "already" in stderr:
            return "?False?%s has been created"%settingtype
        elif "No such file" in stderr:
            return "?False?Need docker swarm init"
        else:
            return "?False?%s"%stderr
    except:
        return "?False?this is a big bug for stderr"

def success_info(items, settingtype):
    pass
    
    try:
        stdout = items.host_ok.values()[0]._result['stdout']
    except:
        pass

    try:
        if "To add a worker to this swarm" in stdout:
            return "?True?docker init success"
        elif "No such file" in stderr:
            return "?False?Need docker swarm init"
        else:
            return "?False?%s"%stderr
    except:
        return "?False?this is a big bug for stderr"


