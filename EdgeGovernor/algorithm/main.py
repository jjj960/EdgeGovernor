import uvicorn
from fastapi import FastAPI
from fastapi import Request

from algorithm.forcast import start
from algorithm.scheduler.edgegovernor import edgegovernor
# 创建FastAPI实例
app = FastAPI()

# 定义POST请求处理函数
@app.get("/forcast")
async def process_forcast_msg():
    result = start.startForcast()
    print(result)
    return result

@app.post("/scheduler")
async def process_scheduler_msg(request: Request):
    data = await request.json()  # 获取请求的JSON数据
    print("请求数据为：",data)
    result = edgegovernor.StartScheduler(data)  # 将数据传递给处理函数
    print("结果为：",result)
    return result

# 运行FastAPI应用
if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=50052)
