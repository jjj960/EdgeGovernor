# 随机森林算法进行负载预测
import numpy as np
import pandas
from sklearn.ensemble import RandomForestRegressor
from sklearn.preprocessing import MinMaxScaler


def fitModel(data):
    # (1) 初始化模型
    clf = RandomForestRegressor(n_estimators=50, random_state=42)
    # (2) 获取数据集
    sourceChange = data
    # (3) 平滑处理
    first, second = sourceChange[0], sourceChange[1]
    # 定义滑动窗口大小
    window_size = 3
    # 创建平滑卷积核
    kernel = np.ones(window_size) / window_size
    # 使用convolve函数进行滑动平均处理
    sourceChange = np.convolve(sourceChange, kernel, mode='valid')
    sourceChange = np.insert(sourceChange, 0, [first, second])
    # (4) 数据归一化，归一化至0~1之间
    sourceChange = sourceChange.reshape(-1, 1)
    scaler = MinMaxScaler(feature_range=(0, 1))
    sourceChange = scaler.fit_transform(sourceChange)
    df = np.array(sourceChange)
    # (5) 构造数据集
    num_samples = len(df)  # 样本个数
    kim = 6  # 延时步长(kim个历史数据作为自变量)
    zim = 1  # 跨zim个时间点进行预测
    df1 = np.zeros((num_samples - zim - kim + 1, kim + zim))
    # 循环构造
    for i in range(num_samples - kim - zim):
        df2 = df[i:i + kim]
        df2 = np.append(df2, df[i + kim + zim - 1])
        df1[i, :] = df2
    df1 = pandas.DataFrame(df1)
    dataSet_X = df1.iloc[:, 0:len(df1.columns) - 1]
    dataSet_Y = df1.iloc[:, len(df1.columns) - 1]
    # 训练模型
    clf.fit(dataSet_X, dataSet_Y)

    # 预测未来5个样本的值
    future_X = np.zeros((1, kim + zim))
    future_X[0, :kim] = df[-kim:].reshape(-1)  # 使用最后kim个历史数据作为输入特征
    pred_future = clf.predict(dataSet_X[-6:])
    pred_future = scaler.inverse_transform(pred_future.reshape(-1, 1))  # 反归一化处理
    return clf, scaler, pred_future