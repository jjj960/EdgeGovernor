package utils

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func GenerateTaskNameRandomString(length int) (string, error) {
	if length <= 0 {
		return "", errors.New("invalid length")
	}

	rand.Seed(time.Now().UnixNano())
	charset := "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result), nil
}

func GenerateAESRandomString(length int) (string, error) {
	if length <= 0 {
		return "", errors.New("invalid length")
	}

	rand.Seed(time.Now().UnixNano())
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result), nil
}

func GetSubString(str string) (string, error) {
	firstIndex := strings.Index(str, "_")
	if firstIndex == -1 {
		return "", fmt.Errorf("underscore not found")
	}

	secondIndex := strings.Index(str[firstIndex+1:], "_")
	if secondIndex == -1 {
		return "", fmt.Errorf("second underscore not found")
	}

	return str[firstIndex+1 : firstIndex+1+secondIndex], nil
}
