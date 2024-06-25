import copy
import json
from dataset.podList import PodsList
from dataset.nodeList import NodesList
from dataset.utils import getpodsMem, getpodsDisk, getnodesDisk, getpodsNet, getpodsCpu, getnodesNet, getnodesMem, \
    getnodesCpu
from algorithm.lrp import getLRPScore
from algorithm.bra import getBRAScore
from algorithm.topsis import StartScheduler
import matplotlib.pyplot as plt


def arrangePodToNode(podName, nodes, algorithm):
    targetNode = ""
    if algorithm == "topsis":
        targetNode = StartScheduler()
    if algorithm == "bra":
        targetNode = getBRAScore(podName)
    if algorithm == "lrp":
        targetNode = getLRPScore(podName)
    podCpu = getpodsCpu(podName)
    podMem = getpodsMem(podName)
    podNet = getpodsNet(podName)
    podDisk = getpodsDisk(podName)
    nodeCpu = getnodesCpu(targetNode)
    nodeMem = getnodesMem(targetNode)
    nodeNet = getnodesNet(targetNode)
    nodeDisk = getnodesDisk(targetNode)
    nodes[targetNode][5].append(podName)
    nodes[targetNode][6][0] = nodes[targetNode][6][0] + (podCpu / nodeCpu)
    nodes[targetNode][6][1] = nodes[targetNode][6][1] + (podMem / nodeMem)
    nodes[targetNode][6][2] = nodes[targetNode][6][2] + (podNet / nodeNet)
    nodes[targetNode][6][3] = nodes[targetNode][6][3] + (podDisk / nodeDisk)
    return nodes


def StartBRA():
    STD = []
    UTL = []
    Net = []
    CPU = []
    Disk = []
    Mem = []
    U_avg = 0
    length = len(PodsList)
    needArrangePodName = list(PodsList.keys())
    L = 0
    for i in range(length):
        newNodes = arrangePodToNode(needArrangePodName[i], NodesList, "bra")
        beautiful_nodes_out = json.dumps(newNodes, indent=4, ensure_ascii=False)
        if i % 1 == 0:
            U_avg = 0
            std = 0
            for key in newNodes:
                # U_avg表示节点 i 上各个资源利用率总和的平均值
                U_avg = (newNodes[key][6][0] + newNodes[key][6][1] + newNodes[key][6][2]  #
                         + newNodes[key][6][3]) / 4
                # 表示节点 i 上各个资源利用率的标准差,即资源失衡度
                std += ((newNodes[key][6][0] - U_avg) ** 2 + (newNodes[key][6][1] - U_avg) ** 2 +
                        (newNodes[key][6][2] - U_avg) ** 2 + (newNodes[key][6][3] - U_avg) ** 2) ** 0.5
            # STD表示所有节点的平均资源失衡度
            STD.append(std / 31)
            UTL_avg = 0
            UTL_sum = 0
            net = 0
            disk = 0
            cpu = 0
            net_sum = 0
            disk_sum = 0
            cpu_sum = 0
            mem = 0
            mem_sum = 0
            for key in newNodes:
                UTL_avg = (newNodes[key][6][0] + newNodes[key][6][1] + newNodes[key][6][2]
                           + newNodes[key][6][3]) / 4
                UTL_sum += UTL_avg
                cpu = newNodes[key][6][0]
                cpu_sum += cpu
                mem = newNodes[key][6][1]
                mem_sum += mem
                net = newNodes[key][6][2]
                net_sum += net
                disk = newNodes[key][6][3]
                disk_sum += disk
            UTL.append(UTL_sum / 31)
            CPU.append(cpu_sum / 31)
            Mem.append(mem_sum / 31)
            Net.append(net_sum / 31)
            Disk.append(disk_sum / 31)

    cpuNodeLoad = []
    memNodeLoad = []
    netNodeLoad = []
    diskNodeLoad = []

    for host_name, data in NodesList.items():
        cpuNodeLoad.append(data[6][0])
        memNodeLoad.append(data[6][1])
        netNodeLoad.append(data[6][2])
        diskNodeLoad.append(data[6][3])

    return STD, UTL, CPU, Mem, Net, Disk, cpuNodeLoad, memNodeLoad, netNodeLoad, diskNodeLoad


def StartLRP():
    STD = []
    UTL = []
    Net = []
    CPU = []
    Disk = []
    Mem = []
    U_avg = 0
    length = len(PodsList)
    needArrangePodName = list(PodsList.keys())
    L = 0
    for i in range(length):
        newNodes = arrangePodToNode(needArrangePodName[i], NodesList, "lrp")
        beautiful_nodes_out = json.dumps(newNodes, indent=4, ensure_ascii=False)
        if i % 1 == 0:
            U_avg = 0
            std = 0
            for key in newNodes:
                # U_avg表示节点 i 上各个资源利用率总和的平均值
                U_avg = (newNodes[key][6][0] + newNodes[key][6][1] + newNodes[key][6][2]  #
                         + newNodes[key][6][3]) / 4
                # 表示节点 i 上各个资源利用率的标准差,即资源失衡度
                std += ((newNodes[key][6][0] - U_avg) ** 2 + (newNodes[key][6][1] - U_avg) ** 2 +
                        (newNodes[key][6][2] - U_avg) ** 2 + (newNodes[key][6][3] - U_avg) ** 2) ** 0.5
            # STD表示所有节点的平均资源失衡度
            STD.append(std / 31)
            UTL_avg = 0
            UTL_sum = 0
            net = 0
            disk = 0
            cpu = 0
            net_sum = 0
            disk_sum = 0
            cpu_sum = 0
            mem = 0
            mem_sum = 0
            for key in newNodes:
                UTL_avg = (newNodes[key][6][0] + newNodes[key][6][1] + newNodes[key][6][2]
                           + newNodes[key][6][3]) / 4
                UTL_sum += UTL_avg
                cpu = newNodes[key][6][0]
                cpu_sum += cpu
                mem = newNodes[key][6][1]
                mem_sum += mem
                net = newNodes[key][6][2]
                net_sum += net
                disk = newNodes[key][6][3]
                disk_sum += disk
            UTL.append(UTL_sum / 31)
            CPU.append(cpu_sum / 31)
            Mem.append(mem_sum / 31)
            Net.append(net_sum / 31)
            Disk.append(disk_sum / 31)

    cpuNodeLoad = []
    memNodeLoad = []
    netNodeLoad = []
    diskNodeLoad = []

    for host_name, data in NodesList.items():
        cpuNodeLoad.append(data[6][0])
        memNodeLoad.append(data[6][1])
        netNodeLoad.append(data[6][2])
        diskNodeLoad.append(data[6][3])

    return STD, UTL, CPU, Mem, Net, Disk, cpuNodeLoad, memNodeLoad, netNodeLoad, diskNodeLoad


def StartTopsis():
    STD = []
    UTL = []
    Net = []
    CPU = []
    Disk = []
    Mem = []
    U_avg = 0
    length = len(PodsList)
    needArrangePodName = list(PodsList.keys())
    L = 0
    for i in range(length):
        newNodes = arrangePodToNode(needArrangePodName[i], NodesList, "topsis")
        beautiful_nodes_out = json.dumps(newNodes, indent=4, ensure_ascii=False)
        if i % 1 == 0:
            U_avg = 0
            std = 0
            for key in newNodes:
                # U_avg表示节点 i 上各个资源利用率总和的平均值
                U_avg = (newNodes[key][6][0] + newNodes[key][6][1] + newNodes[key][6][2]  #
                         + newNodes[key][6][3]) / 4
                # 表示节点 i 上各个资源利用率的标准差,即资源失衡度
                std += ((newNodes[key][6][0] - U_avg) ** 2 + (newNodes[key][6][1] - U_avg) ** 2 +
                        (newNodes[key][6][2] - U_avg) ** 2 + (newNodes[key][6][3] - U_avg) ** 2) ** 0.5
            # STD表示所有节点的平均资源失衡度
            STD.append(std / 31)
            UTL_avg = 0
            UTL_sum = 0
            net = 0
            disk = 0
            cpu = 0
            net_sum = 0
            disk_sum = 0
            cpu_sum = 0
            mem = 0
            mem_sum = 0
            for key in newNodes:
                UTL_avg = (newNodes[key][6][0] + newNodes[key][6][1] + newNodes[key][6][2]
                           + newNodes[key][6][3]) / 4
                UTL_sum += UTL_avg
                cpu = newNodes[key][6][0]
                cpu_sum += cpu
                mem = newNodes[key][6][1]
                mem_sum += mem
                net = newNodes[key][6][2]
                net_sum += net
                disk = newNodes[key][6][3]
                disk_sum += disk
            UTL.append(UTL_sum / 31)
            CPU.append(cpu_sum / 31)
            Mem.append(mem_sum / 31)
            Net.append(net_sum / 31)
            Disk.append(disk_sum / 31)

    cpuNodeLoad = []
    memNodeLoad = []
    netNodeLoad = []
    diskNodeLoad = []

    for host_name, data in NodesList.items():
        cpuNodeLoad.append(data[6][0])
        memNodeLoad.append(data[6][1])
        netNodeLoad.append(data[6][2])
        diskNodeLoad.append(data[6][3])

    return STD, UTL, CPU, Mem, Net, Disk, cpuNodeLoad, memNodeLoad, netNodeLoad, diskNodeLoad


if __name__ == '__main__':
    # BRASTD, BRAUTL, BRACPU, BRAMem, BRANet, BRADisk,BRAcpuNodeLoad, BRAmemNodeLoad, BRAnetNodeLoad, BRAdiskNodeLoad = StartBRA()
    LRPSTD, LRPUTL, LRPCPU, LRPMem, LRPNet, LRPDisk,LRPcpuNodeLoad, LRPmemNodeLoad, LRPnetNodeLoad, LRPdiskNodeLoad = StartLRP()
    # TopSTD, TopUTL, TopCPU, TopMem, TopNet, TopDisk,TopcpuNodeLoad, TopmemNodeLoad, TopnetNodeLoad, TopdiskNodeLoad = StartTopsis()

    # 将每个数组的数据写入到单独的 txt 文件
    arrays = [LRPSTD, LRPUTL, LRPCPU, LRPMem, LRPNet, LRPDisk,LRPcpuNodeLoad, LRPmemNodeLoad, LRPnetNodeLoad, LRPdiskNodeLoad]
    array_names = ['STD', 'UTL', 'CPU', 'Mem', 'Net', 'Disk', 'cpuNodeLoad', 'memNodeLoad',
                   'netNodeLoad', 'diskNodeLoad']

    for array, array_name in zip(arrays, array_names):
        with open(f'./result/800/lrp/{array_name}.txt', 'w') as file:
            for item in array:
                file.write(f"{item}\n")








