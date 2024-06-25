package utils

import (
	"EdgeGovernor/pkg/sec"
	"context"
	"fmt"
)

func BackUpNodeMap() error { //备份关键map
	data, err := Jsoniter.Marshal(NodeTables)

	if err != nil {
		return fmt.Errorf("Error marshalling JSON: %s", err)
	}

	encryptMap := sec.Safer.Encrypt(data) //加密map

	_, err = ETCDCli.Put(context.Background(), "/menet/backup/NodeTables", string(encryptMap))
	if err != nil {
		return fmt.Errorf("failed to put map in etcd: %s", err)
	}

	return nil
}

func GetNodeTablesMap() error {
	data, err := ETCDCli.Get(context.Background(), "/menet/backup/NodeTables")
	if err != nil {
		return fmt.Errorf("failed to get map in etcd: %s", err)
	}

	for _, kv := range data.Kvs {
		decryptMap := sec.Safer.Decrypt(kv.Value)

		var tempMap map[string]Node
		err := Jsoniter.Unmarshal([]byte(decryptMap), &tempMap)
		if err != nil {
			return fmt.Errorf("Error unmarshalling Map JSON: %s", err)
		}

		for key, value := range tempMap {
			NodeTables.Set(key, &value)
		}
	}

	return nil
}
