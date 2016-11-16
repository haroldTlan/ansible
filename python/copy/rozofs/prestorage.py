import re
import os
import ruamel.yaml
import ruamel.yaml.util
import commands

import argparse
import yaml


parser = argparse.ArgumentParser(description="rozofs")
parser.add_argument("--sequence", help="default: --sequence=0", default='None')
#parser.add_argument("--ip", help="default: --ip=192.168.2.149", default='192.168.2.149')
#parser.add_argument("--expand", help="default: --expand=None", default='None')


def preStorage(seq="None"):

    if seq=="None":
        return "?False?None"

    #cid, sid = rule(seq)
    #sf = open("/etc/rozofs/storage.conf")
    conf = open("/home/zonion/speedio/speedio.conf", "r")


    try:
       # storages = sf.read()

        #pattern = r'root = "(.*?)"'
        #root = re.search(pattern, storages)

        result, indent, block_seq_indent = ruamel.yaml.util.load_yaml_guess_indent(
            conf, preserve_quotes=True)
        conf.close()
        mountFile = "/srv/rozofs/storages/storage_%s/0"%(seq)
        os.system("mkdir -p %s"%mountFile)

        result['nas']['mount_dirs'][0]= ruamel.yaml.scalarstring.SingleQuotedScalarString(mountFile)
        with open('/home/zonion/speedio/speedio.conf', 'w') as conf:
            ruamel.yaml.round_trip_dump(result, conf, indent=indent,block_seq_indent=block_seq_indent)

    finally:
        #sf.close()
        conf.close()

    pid = getPid("admd")
    for i in pid:
        try:
            os.system("kill %s"%i)
        except:
            continue

def getPid(process):
    cmd = "ps aux | grep '%s' " % process
    info = commands.getoutput(cmd)
    infos = [i.split()[1] for i in info.split("\n")]
    return infos

def rule(seq):
    num = int(seq)
    
    temp = num % 4
    if temp == 0:
        return num/4, 4
    else:
        return num/4+1, temp



if __name__ == '__main__':
    args = parser.parse_args()
    preStorage(args.sequence)
