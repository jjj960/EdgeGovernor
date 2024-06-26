from concurrent.futures import ProcessPoolExecutor

from algorithm.forcast.randomForest.randomForest import fitModel
from utils.readDB import getLiveNode
from utils.readDB import getNodeMsg
from utils.readDB import getNodeWorkload


def predict_all_nodes_load(nodes, data):
    print("开始进行负载预测")
    # 存储所有结果的列表
    all_results = []

    for node in nodes:
        process_pool = ProcessPoolExecutor(max_workers=4)
        cpu = data[node]["cpu"]
        mem = data[node]["mem"]
        net = data[node]["net"]
        disk = data[node]["disk"]
        print("开始预测节点"+node+"的数据")
        future1 = process_pool.submit(fitModel, cpu)
        future2 = process_pool.submit(fitModel, mem)
        future3 = process_pool.submit(fitModel, net)
        future4 = process_pool.submit(fitModel, disk)

        clf1, scaler1, pred_future1 = future1.result()
        clf2, scaler2, pred_future2 = future2.result()
        clf3, scaler3, pred_future3 = future3.result()
        clf4, scaler4, pred_future4 = future4.result()

        process_pool.shutdown()

        all_results.append([node, "cpu", pred_future1.flatten()[4]])
        all_results.append([node, "mem", pred_future2.flatten()[4]])
        all_results.append([node, "net", pred_future3.flatten()[4]])
        all_results.append([node, "disk", pred_future4.flatten()[4]])


    finialresult = []

    # 创建一个字典用于存储每个节点的负载信息
    node_loads = {}

    for node, category, value in all_results:
        if node not in node_loads:
            # 初始化一个包含节点名称和四个负载类型的空列表
            node_loads[node] = [node, None, None, None, None, 0]

        # 根据负载类型将负载值插入到相应的位置
        if category == "cpu":
            node_loads[node][1] = value
        elif category == "mem":
            node_loads[node][2] = value
        elif category == "net":
            node_loads[node][4] = value
        elif category == "disk":
            node_loads[node][3] = value

    # 将负载信息从字典转换为有序的列表
    for node in node_loads.values():
        finialresult.append(node)

    return finialresult

def dataprocess(finialresult):
    for item in finialresult:
        nodeName = item[0]
        nodeMsg = getNodeMsg(nodeName)

        insert_index = 1
        ip_address = nodeMsg[0][0]
        cpu_capacity = nodeMsg[0][1] * 1000
        remaining_cpu_capacity = cpu_capacity * ((100 - item[1]) / 100)
        memory_capacity = nodeMsg[0][2]
        remaining_memory_capacity = memory_capacity * ((100 - item[2]) / 100)
        disk_capacity = nodeMsg[0][3]
        remaining_disk_capacity = disk_capacity * ((100 - item[3]) / 100)
        print(item[4])
        print(disk_capacity)
        print(remaining_disk_capacity)
        network_bandwidth = nodeMsg[0][4]

        item.insert(insert_index, ip_address)
        item.insert(insert_index + 2, cpu_capacity)
        item.insert(insert_index + 3, remaining_cpu_capacity)
        item.insert(insert_index + 5, memory_capacity)
        item.insert(insert_index + 6, remaining_memory_capacity)
        item.insert(insert_index + 8, disk_capacity)
        item.insert(insert_index + 9, remaining_disk_capacity)
        item.append(network_bandwidth)

def getDBData():    #获取每个节点的负载列表
    data = {}
    liveNode = getLiveNode()
    for node in liveNode:
       time, cpu, mem, net, disk = getNodeWorkload(node)
       data[node] = {
           "time": time,
           "cpu": cpu,
           "mem": mem,
           "net": net,
           "disk": disk
       }
    return liveNode ,data

def startForcast():
    # 获取所有节点的列表
    nodes, data = getDBData()
    # 并发预测每个节点的负载
    finialresult = predict_all_nodes_load(nodes, data)
    dataprocess(finialresult)
    print(finialresult)
    return finialresult

if __name__ == '__main__':
    startForcast()