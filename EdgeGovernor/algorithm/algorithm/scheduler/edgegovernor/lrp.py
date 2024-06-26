
from dataset.utils import getpodsCpu,getpodsMem
from dataset.nodeList import NodesList


def getLRPScore(podName):
    score = {}
    # 需要部署的pod的cpu和内存
    NeedAddCpu = getpodsCpu(podName)
    NeedAddRam = getpodsMem(podName)
    totalMilliCPU = 0
    totalRam = 0
    # 单个node的总的cpu和内存
    for key in NodesList:
        # node上已有pod的cpu和内存
        ExistCpu = NodesList[key][0] * NodesList[key][6][0]
        ExistRam = NodesList[key][1] * NodesList[key][6][1]

        if NodesList[key][5]:
            # for pod in NodesList[key][5]:  # 'node10': [64, 49152, 40, 1200, 220, [], [0.004172704293255518, 0.006199530842266317, 0.005952556100582241, 0.0005354172486981479, 0.007666170462532992]],
            #     ExistCpu = ExistCpu + getpodsCpu(pod)
            #     ExistRam = ExistRam + getpodsRam(pod)
            # totalCPU = ExistCpu + NeedAddCpu    #总共使用的CPU数量
            #
            # totalRam = ExistRam + NeedAddRam    #总共使用的内存数量
            totalCPU = ExistCpu
            totalRam = ExistRam
            score[key] = (((NodesList[key][0] - totalCPU) * 10) / NodesList[key][0]
                          + ((NodesList[key][1] - totalRam) * 10) / NodesList[key][1]) / 2
        else:
            score[key] = 10
    # 返回字典中value最大的键
    return max(score, key=score.get)