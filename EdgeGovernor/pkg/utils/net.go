package utils

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os/exec"
	"regexp"
	"strings"
)

func BytesToMB(bytes int64) (int64, error) {
	if bytes < 0 {
		return 0, fmt.Errorf("bytes should be a positive value")
	}

	mb := bytes / (1024 * 1024)
	return mb, nil
}

func GetIP() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve network interfaces: %v", err)
	}

	var ip string

	for _, iface := range interfaces {
		if iface.Flags&net.FlagLoopback != 0 || iface.Flags&net.FlagUp == 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			return "", fmt.Errorf("failed to retrieve interface addresses: %v", err)
		}

		for _, addr := range addrs {
			if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
				if ipNet.IP.To4() != nil {
					ip = ipNet.IP.String()
					break
				}
			}
		}

		if ip != "" {
			break
		}
	}

	if ip == "" {
		out, err := exec.Command("hostname", "-I").Output()
		if err != nil {
			return "", fmt.Errorf("failed to execute command: %v", err)
		}

		ip = strings.TrimSpace(string(out))
	}

	if ip == "" {
		return "", errors.New("failed to retrieve IP address")
	}

	return ip, nil
}

func GetNetbandwidth(name string) (float64, error) {
	// 执行命令
	cmd := exec.Command("ethtool", name)
	// 捕获命令输出
	output, err := cmd.Output()
	if err != nil {
		log.Println("命令执行出错:", err)
	}
	// 提取Speed字段内容
	re := regexp.MustCompile(`Speed: (\d+)`)
	result := re.FindStringSubmatch(string(output))

	// 检查是否找到匹配项
	if len(result) < 2 {
		fmt.Println("未找到Speed字段")
	}

	// 输出Speed字段内容
	speed := result[1]
	return StringtoFloat64(speed)
}
