package etcd

import (
	"EdgeGovernor/pkg/constants"
	"fmt"
	"io/ioutil"
	"os/exec"
)

func GetAvailableETCDEndPoints() []string { //获取可用的ETCD节点
	var EndPoints []string
	if constants.Etcd1Enable {
		EndPoints = append(EndPoints, constants.Etcd1URL)
	}
	if constants.Etcd2Enable {
		EndPoints = append(EndPoints, constants.Etcd2URL)
	}
	if constants.Etcd3Enable {
		EndPoints = append(EndPoints, constants.Etcd3URL)
	}
	if constants.Etcd4Enable {
		EndPoints = append(EndPoints, constants.Etcd4URL)
	}
	if len(EndPoints) == 0 {
		panic("There are no available ETCD databases in the cluster !")
	}
	return EndPoints
}

func ServiceCreate(serviceName string, text string) error {
	err := ioutil.WriteFile("/usr/lib/systemd/system/"+serviceName+".service", []byte(text), 0644)
	if err != nil {
		return fmt.Errorf("Failed to write to the etcd service file:", err)
	}

	err = exec.Command("systemctl", "daemon-reload").Run()
	if err != nil {
		return fmt.Errorf("Failed to execute the systemctl daemon reload command:", err)
	}

	err = exec.Command("systemctl", "enable", "etcd").Run()
	if err != nil {
		return fmt.Errorf("Failed to execute the systemctl enable etcd command:", err)
	}

	err = exec.Command("systemctl", "start", "etcd").Run()
	if err != nil {
		return fmt.Errorf("Failed to execute the systemctl start etcd command:", err)
	}

	err = ServiceCheck(serviceName)
	if err != nil {
		return fmt.Errorf("Etcd service is not running properly:", err)
	}

	return nil
}

func ServiceCheck(serviceName string) error {
	cmd := exec.Command("systemctl", "is-active", serviceName)
	output, err := cmd.Output()
	if err != nil {
		return err
	}

	status := string(output)
	if status == "active\n" { // 运行状态为 "active\n" 表示服务正在运行
		return nil
	}

	return nil
}
