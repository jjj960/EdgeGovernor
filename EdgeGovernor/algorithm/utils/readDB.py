import socket

import duckdb

def getNodeWorkload():
    conn = duckdb.connect('clouddb')
    conn.execute("SELECT Timestamp, CPUUsagePercent, MemoryUsedPercent, BytesRecv, DiskUsedPercent FROM Hostload ORDER BY Timestamp DESC LIMIT 500")
    results = conn.fetchall()
    print(results)
    time = [row[0] for row in results]
    cpu = [round(float(row[1])) for row in results]
    mem = [round(float(row[2])) for row in results]
    net = [round(float(row[3])) for row in results]
    disk = [round(float(row[4])) for row in results]

    return time, cpu, mem, net, disk