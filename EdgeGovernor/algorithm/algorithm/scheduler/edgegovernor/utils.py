

def ArrayToDict(data):
    result = {}
    check = False
    if len(data) > 1 :
        check = True
        for item in data:
            host_name = item[0]
            cpu_percent = float(item[2])
            cpu_total = float(item[3])
            cpu_remaining = float(item[4])
            mem_percent = float(item[5])
            mem_total = float(item[6])
            mem_remaining = float(item[7])
            disk_percent = float(item[8])
            disk_total = float(item[9])
            disk_remaining = float(item[10])
            bitrate = float(item[11])
            if (cpu_percent > 80) or (mem_percent > 80) or (disk_percent > 80):
                cpu_remaining, mem_remaining, disk_remaining = 0, 0, 0

            result[host_name] = {
                'cpu_percent': cpu_percent,
                'cpu_total': cpu_total,
                'cpu_remaining': cpu_remaining,
                'mem_percent': mem_percent,
                'mem_total': mem_total,
                'mem_remaining': mem_remaining,
                'disk_percent': disk_percent,
                'disk_total': disk_total,
                'disk_remaining': disk_remaining,
                'bitrate': bitrate
            }

    return result, check




