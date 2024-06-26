import requests
import json

data = [[1, 2, 3], [4, 5, 6], [7, 8, 9]]  # 二维数组数据

url = "http://localhost:50052/scheduler"  # 目标URL

payload = json.dumps(data)  # 将二维数组转换为JSON字符串

headers = {
    "Content-Type": "application/json"
}

response = requests.post(url, data=payload, headers=headers)

print(response.text)  # 打印服务器返回的响应内容
