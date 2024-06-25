package utils

import (
	"fmt"

	"os"
	"path/filepath"
)

func CreateFolder(folderPath string) error { //不能创建多级文件夹
	// 检查文件夹是否存在
	_, err := os.Stat(folderPath)
	if os.IsNotExist(err) {
		// 文件夹不存在，创建文件夹
		err := os.MkdirAll(folderPath, 0755)
		if err != nil {
			return fmt.Errorf("Failed to create folder: %s", err.Error())
		}
	} else if err != nil {
		return fmt.Errorf("Error checking folder: %s", err.Error())
	}

	return nil
}

func DeleteFolder(folderPath string) error {
	// 检查文件夹是否存在
	if _, err := os.Stat(folderPath); err == nil {
		// 文件夹存在，删除文件夹及其内容
		err := os.RemoveAll(folderPath)
		if err != nil {
			return fmt.Errorf("Failed to delete folder: %s", err.Error())
		}
	} else if os.IsNotExist(err) {
		return fmt.Errorf("Folder does not exist")
	} else {
		return fmt.Errorf("Error checking folder: %s", err.Error())
	}

	return nil
}

func StaticFolder(filePath string) (float64, error) {
	totalSize := int64(0)

	err := filepath.Walk(filePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return nil
		}
		if !info.IsDir() {
			totalSize += info.Size()
		}
		return nil
	})

	if err != nil {
		fmt.Println(err)
		return 0, nil
	}

	totalSizeMB := float64(totalSize) / 1024 / 1024

	return totalSizeMB, nil
}

func IsFirstRun() bool {
	// 定义标记文件路径
	flagFilePath := "firstrun.flag"

	// 检查标记文件是否存在
	_, err := os.Stat(flagFilePath)
	if os.IsNotExist(err) {
		// 标记文件不存在，表示是第一次启动
		// 创建标记文件
		file, err := os.Create(flagFilePath)
		if err != nil {
			fmt.Println("Error creating flag file:", err)
		}
		defer file.Close()
		return true
	}

	// 标记文件存在，表示不是第一次启动
	return false
}
