from dataset.nodeList import NodesList
from dataset.podList import PodsList

'''
根据pod名获取cpu
'''
def getpodsCpu(name):
    cpu = PodsList.get(name)[0]
    return cpu

'''
根据pod名获取内存
'''
def getpodsMem(name):
    mem = PodsList.get(name)[1]
    return mem

'''
根据pod名获取net
'''
def getpodsNet(name):
    net = PodsList.get(name)[2]
    return net

'''
根据pod名获取磁盘
'''
def getpodsDisk(name):
    disk = PodsList.get(name)[3]
    return disk



'''
根据node名获取cpu
'''
def getnodesCpu(name):
    cpu = NodesList.get(name)[0]
    return cpu

'''
根据node名获取内存
'''
def getnodesMem(name):
    mem = NodesList.get(name)[1]
    return mem

'''
根据node名获取net
'''
def getnodesNet(name):
    net = NodesList.get(name)[2]
    return net

'''
根据node名获取磁盘
'''

def getnodesDisk(name):
    disk = NodesList.get(name)[3]
    return disk


