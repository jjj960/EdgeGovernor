import socket
from concurrent.futures import ProcessPoolExecutor

from algorithm.forcast.randomForest.randomForest import fitModel
from utils.readDB import getNodeWorkload


def predict_all_nodes_load(data):
    print("开始进行负载预测")
    # 存储所有结果的列表
    all_results = []
    node = socket.gethostname()
    process_pool = ProcessPoolExecutor(max_workers=4)
    cpu = data[node]["cpu"]
    mem = data[node]["mem"]
    net = data[node]["net"]
    disk = data[node]["disk"]
    print("开始预测节点" + node + "的数据")
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

    node_loads = {}

    for node, category, value in all_results:
        if node not in node_loads:
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


def getDBData():    #获取每个节点的负载列表
    data = {}
    time, cpu, mem, net, disk = getNodeWorkload()
    data[socket.gethostname()] = {
        "time": time,
        "cpu": cpu,
        "mem": mem,
        "net": net,
        "disk": disk
    }
    return data

def startForcast():
    # 获取所有节点的列表
    data = getDBData()
    # 并发预测每个节点的负载
    finialresult = predict_all_nodes_load(data)
    #返回的数据为[['hys', 46.52000000000004, 70.338, 54.7900111756453, 4305.053333333335, 0]]    [[hostname, cpupercent, mempercent, diskpercent, net, 0]]

    return finialresult

if __name__ == '__main__':
    startForcast()