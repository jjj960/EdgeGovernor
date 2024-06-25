import matplotlib.pyplot as plt


# 定义一个函数来读取 txt 文件并将数据放入数组
def read_array_from_txt(filename):
    data = []
    with open(filename, 'r') as file:
        for line in file:
            # 假设每行数据是一个数字，如果不是数字需要进行类型转换
            value = float(line.strip())  # 去除换行符并转换为数字
            data.append(value)
    return data


podNum = "800"

# 调用函数读取数据并放入数组
bracpuLoad = read_array_from_txt("./result/" + podNum + "/bra/cpuNodeLoad.txt")
bramemLoad = read_array_from_txt("./result/" + podNum + "/bra/memNodeLoad.txt")
branetLoad = read_array_from_txt("./result/" + podNum + "/bra/netNodeLoad.txt")
bradiskLoad = read_array_from_txt("./result/" + podNum + "/bra/diskNodeLoad.txt")
bracpu = read_array_from_txt("./result/" + podNum + "/bra/CPU.txt")
bramem = read_array_from_txt("./result/" + podNum + "/bra/Mem.txt")
branet = read_array_from_txt("./result/" + podNum + "/bra/Net.txt")
bradisk = read_array_from_txt("./result/" + podNum + "/bra/Disk.txt")
brastd = read_array_from_txt("./result/" + podNum + "/bra/STD.txt")
brautl = read_array_from_txt("./result/" + podNum + "/bra/UTL.txt")

lrpcpuLoad = read_array_from_txt("./result/" + podNum + "/lrp/cpuNodeLoad.txt")
lrpmemLoad = read_array_from_txt("./result/" + podNum + "/lrp/memNodeLoad.txt")
lrpnetLoad = read_array_from_txt("./result/" + podNum + "/lrp/netNodeLoad.txt")
lrpdiskLoad = read_array_from_txt("./result/" + podNum + "/lrp/diskNodeLoad.txt")
lrpcpu = read_array_from_txt("./result/" + podNum + "/lrp/CPU.txt")
lrpmem = read_array_from_txt("./result/" + podNum + "/lrp/Mem.txt")
lrpnet = read_array_from_txt("./result/" + podNum + "/lrp/Net.txt")
lrpdisk = read_array_from_txt("./result/" + podNum + "/lrp/Disk.txt")
lrpstd = read_array_from_txt("./result/" + podNum + "/lrp/STD.txt")
lrputl = read_array_from_txt("./result/" + podNum + "/lrp/UTL.txt")

topcpuLoad = read_array_from_txt("./result/" + podNum + "/topsis/cpuNodeLoad.txt")
topmemLoad = read_array_from_txt("./result/" + podNum + "/topsis/memNodeLoad.txt")
topnetLoad = read_array_from_txt("./result/" + podNum + "/topsis/netNodeLoad.txt")
topdiskLoad = read_array_from_txt("./result/" + podNum + "/topsis/diskNodeLoad.txt")
topcpu = read_array_from_txt("./result/" + podNum + "/topsis/CPU.txt")
topmem = read_array_from_txt("./result/" + podNum + "/topsis/Mem.txt")
topnet = read_array_from_txt("./result/" + podNum + "/topsis/Net.txt")
topdisk = read_array_from_txt("./result/" + podNum + "/topsis/Disk.txt")
topstd = read_array_from_txt("./result/" + podNum + "/topsis/STD.txt")
toputl = read_array_from_txt("./result/" + podNum + "/topsis/UTL.txt")


print(brautl[799],lrputl[799],toputl[799])

# CPU负载比较
finalBRACPULoad = bracpuLoad[30]
finalLRPCPULoad = lrpcpuLoad[30]
finalTopCPULoad = topcpuLoad[30]
averageBRACPULoad = sum(bracpuLoad) / len(bracpuLoad)
averageLRPCPULoad = sum(lrpcpuLoad) / len(lrpcpuLoad)
averageTopCPULoad = sum(topcpuLoad) / len(topcpuLoad)

finalTopToBRACPU = (finalTopCPULoad - finalBRACPULoad) / finalBRACPULoad * 100
finalTopToLRPCPU = (finalTopCPULoad - finalLRPCPULoad) / finalLRPCPULoad * 100

averageTopToBRACPU = (averageTopCPULoad - averageBRACPULoad) / averageBRACPULoad * 100
averageTopToLRPCPU = (averageTopCPULoad - averageLRPCPULoad) / averageLRPCPULoad * 100


print(sum(lrpcpuLoad) / len(lrpcpuLoad),sum(bracpuLoad) / len(bracpuLoad),sum(topcpuLoad) / len(topcpuLoad))


# 计算 topstd 的平均值
average_topstd = sum(topstd) / len(topstd)

# 计算 brastd 相对于 topstd 降低的百分比
average_brastd = sum(brastd) / len(brastd)
percentage_decrease_brastd = ((average_topstd - average_brastd) / average_topstd) * 100

# 计算 lrpstd 相对于 topstd 降低的百分比
average_lrpstd = sum(lrpstd) / len(lrpstd)
percentage_decrease_lrpstd = ((average_topstd - average_lrpstd) / average_topstd) * 100

# 打印结果
print(f"topstd 的平均值: {average_topstd:.4f}")
print(f"bra 的平均值: {average_brastd:.4f}")

print(f"lrp 的平均值: {average_lrpstd:.4f}")

print(f"brastd 相对于 topstd 降低了 {percentage_decrease_brastd:.2f}%")
print(f"lrpstd 相对于 topstd 降低了 {percentage_decrease_lrpstd:.2f}%")

# # 计算每个数组的平均值
# average_toputl = sum(toputl) / len(toputl)
# average_lrputl = sum(lrputl) / len(lrputl)
# average_brautl = sum(brautl) / len(brautl)
#
# # 计算 toputl 相对于 lrputl 和 brautl 变化的百分比
# percent_change_toputl_lrputl = ((average_toputl - average_lrputl) / average_lrputl) * 100
# percent_change_toputl_brautl = ((average_toputl - average_brautl) / average_brautl) * 100
#
# print(f"toputl 相对于 lrputl 提高了 {percent_change_toputl_lrputl:.2f}%")
# print(f"toputl 相对于 brautl 提高了 {percent_change_toputl_brautl:.2f}%")

# # 绘制 CPU Load 对比图
# plt.figure()
# plt.plot(bracpuLoad, label='BRA')
# plt.plot(lrpcpuLoad, label='LRP')
# plt.plot(topcpuLoad, label='TOPSIS')
# plt.title('Final CPU Load Comparison')
# plt.legend()
# plt.xlabel('Node ID')
# plt.ylabel('Load Usage')
# plt.savefig('cpu_load_comparison.png')
# plt.show()
#
#
# # 绘制 Mem Load 对比图
# plt.figure()
# plt.plot(bramemLoad, label='BRA')
# plt.plot(lrpmemLoad, label='LRP')
# plt.plot(topmemLoad, label='TOPSIS')
# plt.title('Final Memory Load Comparison')
# plt.legend()
# plt.xlabel('Node ID')
# plt.ylabel('Load Usage')
# plt.savefig('mem_load_comparison.png')
#
# plt.show()
#
#
# # 绘制 Net Load 对比图
# plt.figure()
# plt.plot(branetLoad, label='BRA')
# plt.plot(lrpnetLoad, label='LRP')
# plt.plot(topnetLoad, label='TOPSIS')
# plt.title('Final Network Load Comparison')
# plt.legend()
# plt.xlabel('Node ID')
# plt.ylabel('Load Usage')
# plt.savefig('net_load_comparison.png')
#
# plt.show()
#
#
# # 绘制 Disk Load 对比图
# plt.figure()
# plt.plot(bradiskLoad, label='BRA')
# plt.plot(lrpdiskLoad, label='LRP')
# plt.plot(topdiskLoad, label='TOPSIS')
# plt.title('Final Disk Load Comparison')
# plt.legend()
# plt.xlabel('Node ID')
# plt.ylabel('Load Usage')
# plt.savefig('disk_load_comparison.png')
#
# plt.show()
#
#
# plt.figure()
# plt.plot(bracpu, label='BRA')
# plt.plot(lrpcpu, label='LRP')
# plt.plot(topcpu, label='TOPSIS')
# plt.title('Cpu STD Comparison')
# plt.legend()
# plt.xlabel('Pod Number')
# plt.savefig('cpu_std_comparison.png')
#
# plt.show()
#
#
#
# plt.figure()
# plt.plot(bramem, label='BRA')
# plt.plot(lrpmem, label='LRP')
# plt.plot(topmem, label='TOPSIS')
# plt.title('Mem STD Comparison')
# plt.legend()
# plt.xlabel('Pod Number')
# plt.savefig('mem_std_comparison.png')
#
# plt.show()
#
#
# plt.figure()
# plt.plot(branet, label='BRA')
# plt.plot(lrpnet, label='LRP')
# plt.plot(topnet, label='TOPSIS')
# plt.title('Net STD Comparison')
# plt.legend()
# plt.xlabel('Pod Number')
# plt.savefig('net_std_comparison.png')
#
# plt.show()
#
#
#
# plt.figure()
# plt.plot(bradisk, label='BRA')
# plt.plot(lrpdisk, label='LRP')
# plt.plot(topdisk, label='TOPSIS')
# plt.title('Disk STD Comparison')
# plt.legend()
# plt.xlabel('Pod Number')
# plt.savefig('disk_std_comparison.png')
#
# plt.show()
#
#
#
# plt.figure()
# plt.plot(brastd, label='BRA')
# plt.plot(lrpstd, label='LRP')
# plt.plot(topstd, label='TOPSIS')
# plt.title('Comprehensive Average STD Comparison')
# plt.legend()
# plt.xlabel('Pod Number')
# plt.savefig('comprehensive_std_comparison.png')
#
# plt.show()
#
#
# plt.figure()
# plt.plot(brautl, label='BRA')
# plt.plot(lrputl, label='LRP')
# plt.plot(toputl, label='TOPSIS')
# plt.title('UTL Load Comparison')
# plt.legend()
# plt.xlabel('Pod Number')
# plt.savefig('comprehensive_utl_comparison.png')
#
# plt.show()
#
#
# # 调整布局
# plt.tight_layout()
#
# # 显示图形
# plt.show()
