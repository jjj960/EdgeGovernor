package algorithm

import (
	"EdgeGovernor/pkg/sec"
	"EdgeGovernor/pkg/utils"
	"context"
	"fmt"
)

func AddAlgorithmStatus(name, url, types, status string) error {
	mutex.Lock()
	defer mutex.Unlock()
	AlgorithmStatusMap[name] = AlgorithmStatus{Name: name, URL: url, Type: types, Status: status}

	err := BackUpAlgorithmStatusMap()
	if err != nil {
		return fmt.Errorf("AlgorithmMap backup fail: %s\n", err)
	}

	return nil
}

// 获取特定算法的状态
func GetAlgorithmStatus(name string) string {
	mutex.RLock()
	defer mutex.RUnlock()
	if algStatus, ok := AlgorithmStatusMap[name]; ok {
		return algStatus.Status
	}
	return "Unknown"
}

// 获取特定算法的url
func GetAlgorithmURL(name, types string) string {
	mutex.RLock()
	defer mutex.RUnlock()

	if algStatus, ok := AlgorithmStatusMap[name]; ok {
		if algStatus.Type == types {
			return algStatus.URL
		}
	}
	return "Unknown"
}

// 更新指定算法的状态
func UpdateAlgorithmStatus(name, newStatus string) error {
	mutex.Lock()
	defer mutex.Unlock()

	if algStatus, ok := AlgorithmStatusMap[name]; ok {
		algStatus.Status = newStatus
		AlgorithmStatusMap[name] = algStatus

		err := BackUpAlgorithmStatusMap()
		if err != nil {
			return fmt.Errorf("AlgorithmMap backup fail: %s\n", err)
		}

		return nil
	}

	return fmt.Errorf("Algorithm %s not found", name)
}

// 删除指定算法的状态信息
func DeleteAlgorithmStatus(name string) error {
	mutex.Lock()
	defer mutex.Unlock()
	if _, ok := AlgorithmStatusMap[name]; ok {
		delete(AlgorithmStatusMap, name)

		err := BackUpAlgorithmStatusMap()
		if err != nil {
			return fmt.Errorf("AlgorithmMap backup fail: %s\n", err)
		}

		return nil
	}
	return fmt.Errorf("Algorithm %s not found", name)
}

// 输出整个全局 map 的内容
func PrintAlgorithmStatusMap() []byte {
	mutex.RLock()
	defer mutex.RUnlock()
	jsonData, err := utils.Jsoniter.Marshal(AlgorithmStatusMap)
	if err != nil {
		fmt.Println("Error marshalling map to JSON:", err)
		return nil
	}

	return jsonData
}

func BackUpAlgorithmStatusMap() error { //备份关键map
	data, err := utils.Jsoniter.Marshal(AlgorithmStatusMap)

	if err != nil {
		return fmt.Errorf("Error marshalling JSON: %s", err)
	}

	encryptMap := sec.Safer.Encrypt(data) //加密map

	_, err = utils.ETCDCli.Put(context.Background(), "/menet/backup/algorithmStatusMap", string(encryptMap))
	if err != nil {
		return fmt.Errorf("failed to put map in etcd: %s", err)
	}

	return nil
}

func GetAlgorithmStatusMap() error {
	data, err := utils.ETCDCli.Get(context.Background(), "/menet/backup/algorithmStatusMap")
	if err != nil {
		return fmt.Errorf("failed to get map in etcd: %s", err)
	}

	for _, kv := range data.Kvs {
		decryptMap := sec.Safer.Decrypt(kv.Value)

		var tempMap map[string]AlgorithmStatus
		err := utils.Jsoniter.Unmarshal([]byte(decryptMap), &tempMap)
		if err != nil {
			return fmt.Errorf("Error unmarshalling Map JSON: %s", err)
		}

		for key, value := range tempMap {
			AlgorithmStatusMap[key] = value
		}

	}

	return nil
}
