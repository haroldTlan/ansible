import re
import os
import ruamel.yaml
import ruamel.yaml.util
import commands

sf = open("/etc/rozofs/storage.conf")
conf = open("/home/zonion/speedio/speedio.conf", "r")


try:
    storages = sf.read()

    pattern = r'root = "(.*?)"'
    root = re.search(pattern, storages)

    result, indent, block_seq_indent = ruamel.yaml.util.load_yaml_guess_indent(
        conf, preserve_quotes=True)
    conf.close()

    result['nas']['mount_dirs'][0]= ruamel.yaml.scalarstring.SingleQuotedScalarString(root.group(1) + '/0')
    with open('/home/zonion/speedio/speedio.conf', 'w') as conf:
        ruamel.yaml.round_trip_dump(result, conf, indent=indent,block_seq_indent=block_seq_indent)

finally:
    sf.close()
    conf.close()

def getPid(process):
    cmd = "ps aux | grep '%s' " % process
    info = commands.getoutput(cmd)
    print info
    infos = [i.split()[1] for i in info.split("\n")]
    return infos

pid = getPid("admd")
print pid
for i in pid:
    try:
        os.system("kill %s"%i)
    except:
        continue
