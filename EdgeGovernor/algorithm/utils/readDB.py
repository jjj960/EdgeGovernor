import duckdb


def getLiveNode():
    # 使用DuckDB的Connection连接到数据库文件
    conn = duckdb.connect('test.db')  #改成实际数据库路径
    cursor = conn.cursor()

    # 执行查询语句
    cursor.execute("SELECT Hostname FROM nodes WHERE Status = 'active' AND Role = 'follower'")

    # 获取所有查询结果
    results = cursor.fetchall()
    # 将某一列的值打包为列表
    liveNodes = [row[0] for row in results]

    # 关闭游标（DuckDB连接在Python中通常是自动关闭的）
    cursor.close()

    return liveNodes


def getNodeWorkload(name):
    # 使用DuckDB的Connection连接到数据库文件
    conn = duckdb.connect('test.db')    #改成实际数据库路径
    cursor = conn.cursor()

    table_name = f"{name}_load"

    query = f"SELECT Timestamp, CPUUsagePercent, MemoryUsedPercent, BytesReceived, DiskUsedPercent FROM {table_name} ORDER BY Timestamp DESC LIMIT 500"

    # 执行查询语句
    cursor.execute(query)
    # 获取所有查询结果
    results = cursor.fetchall()
    # 将某一列的值打包为列表，注意这里我们不需要反转结果，因为我们已经按降序排序了
    time = [row[0] for row in results]
    cpu = [round(float(row[1])) for row in results]
    mem = [round(float(row[2])) for row in results]
    net = [round(float(row[3])) for row in results]
    disk = [round(float(row[4])) for row in results]

    # 关闭游标（DuckDB连接在Python中通常是自动关闭的）
    cursor.close()
    return time, cpu, mem, net, disk


def getNodeMsg(name):
    # 使用DuckDB的Connection连接到数据库文件
    conn = duckdb.connect('test.db')
    cursor = conn.cursor()

    # 执行查询语句，使用参数化查询避免SQL注入
    cursor.execute("SELECT IP, CPU, Memory, DiskSpace, Net FROM nodes WHERE Hostname = ?", (name,))
    # 获取所有查询结果
    results = cursor.fetchall()
    # 将每一行的结果打包为列表的列表
    node_info = [list(row) for row in results]

    # 关闭游标（DuckDB连接在Python中通常是自动关闭的）
    cursor.close()

    return node_info