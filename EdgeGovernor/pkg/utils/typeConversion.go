package utils

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"strconv"
)

var Jsoniter = jsoniter.ConfigCompatibleWithStandardLibrary

func StringtoInt64(msg string) (int64, error) {
	// 将字符串转换为浮点数
	f, err := strconv.ParseFloat(msg, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to convert string to floating-point number: %v", err)
	}
	// 将浮点数转换为 int64
	i := int64(f)
	return i, nil
}

func StringtoFloat64(msg string) (float64, error) {
	f, err := strconv.ParseFloat(msg, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to convert string to float64: %v", err)
	}
	return f, nil
}

func Float64toString(msg float64) string {
	str := strconv.FormatFloat(msg, 'f', -1, 64)
	return str
}

func Int64toString(msg int64) string {
	str := strconv.FormatInt(msg, 10)
	return str
}

func InterfaceToString(data interface{}) (string, error) {
	// 使用类型断言将 interface{} 类型的数据转换为 string 类型
	if str, ok := data.(string); ok {
		return str, nil
	}
	return "", fmt.Errorf("Failed to convert interface{} to string")
}
