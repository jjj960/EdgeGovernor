#BalancedResourceAllocation 算法

from dataset.nodeList import NodesList
from dataset.utils import getpodsCpu, getpodsMem

'''
根据pod分配节点
输入：pod的名字，如pod1以及nodes字典
输出：返回score最大的单个node
'''

def getBRAScore(podName):
    score = {}
    # 需要部署的pod的cpu和内存
    NeedAddCpu = getpodsCpu(podName)
    NeedAddMem = getpodsMem(podName)
    totalMilliCPU = 0
    totalRam = 0
    # 单个node的总的cpu和内存
    for key in NodesList:
        # node上已有pod的cpu和内存
        ExistCpu = NodesList[key][0] * NodesList[key][6][0]
        ExistMem = NodesList[key][1] * NodesList[key][6][1]
        if NodesList[key][5]:
            # for pod in NodesList[key][5]:
            #     ExistCpu = ExistCpu + getpodsCpu(pod)
            #     ExistMem = ExistMem + getpodsMem(pod)
            # totalCPU = ExistCpu + NeedAddCpu    #总共使用的CPU数量
            #
            # totalRam = ExistRam + NeedAddRam    #总共使用的内存数量
            totalCPU = ExistCpu
            totalMem = ExistMem
            score[key] = ((10-abs((totalCPU/NodesList[key][0])-(totalMem/NodesList[key][1]))*10) + ((((NodesList[key][0] - totalCPU) * 10) / NodesList[key][0] + ((NodesList[key][1] - totalMem) * 10) / NodesList[key][1]) / 2)) / 2
        else:
            score[key] = 10
    # 返回字典中value最大的键
    return max(score, key=score.get)