package resource

import (
	"EdgeGovernor/pkg/utils"
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	clientv1 "github.com/docker/docker/client"
	"log"
)

func CreateResourceContainer(taskName string, imageName string, cpuSize int64, memSize int64, needPersistence int, containerDataDir []string) (string, error) {
	if !CheckImageExist(imageName) {
		ImagePull(imageName)
	}

	// 配置容器资源限制
	resources := &container.Resources{
		Memory:     memSize * 1024 * 1024, // 内存限制 单位为MB
		MemorySwap: -1,                    // 禁止交换内存
		CPUQuota:   cpuSize * 100,         // CPU配额为50%
		CPUPeriod:  100000,
		CPUShares:  cpuSize,
	}

	switch needPersistence {
	case 0: //不需要数据持久化

		// 创建容器的配置
		config := &container.Config{
			Image: imageName,
			Labels: map[string]string{
				"app":       taskName,
				"namespace": "menet",
			},
		}

		hostConfig := &container.HostConfig{
			Resources: *resources,
		}

		resp, err := utils.DockerCli.ContainerCreate(context.Background(), config, hostConfig, nil, nil, taskName)
		if err != nil {
			log.Println("Container creation error:", err)
			return "", err
		}

		// 启动容器
		if err := utils.DockerCli.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{}); err != nil {
			log.Println("Container startup error:", err)
			return "", err
		}

	case 1: //需要数据持久化,  使用方法, containerDataDir数组中存储需要数据持久化的文件夹位置, volume会默认生成在/var/lib/docker/volumes/文件夹下

		utils.CreateFolder("/data/menet/task/" + taskName)

		// 配置容器资源限制
		resources := &container.Resources{
			Memory:     memSize * 1024 * 1024, // 内存限制 单位为MB
			MemorySwap: -1,                    // 禁止交换内存
			CPUQuota:   cpuSize * 100,         // CPU配额为50%
			CPUPeriod:  100000,
			CPUShares:  cpuSize,
		}
		// 创建容器的配置
		config := &container.Config{
			Image: imageName,
			Labels: map[string]string{
				"app":       taskName,
				"namespace": "menet",
			},
		}

		mounts := []mount.Mount{}

		for _, targetPath := range containerDataDir {
			utils.CreateFolder("/data/menet/task/" + taskName + targetPath)
			m := mount.Mount{
				Type:   mount.TypeBind,
				Source: "/data/menet/task/" + taskName + targetPath, //表示要挂载到容器内部的宿主机目录。
				Target: targetPath,                                  //容器内的文件夹。
			}

			mounts = append(mounts, m)
		}

		hostConfig := &container.HostConfig{
			Resources: *resources,
			Mounts:    mounts,
		}

		resp, err := utils.DockerCli.ContainerCreate(context.Background(), config, hostConfig, nil, nil, taskName)
		if err != nil {
			log.Println("Container creation error:", err)
			return "fail", err
		}

		// 启动容器
		if err := utils.DockerCli.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{}); err != nil {
			log.Println("Container startup error:", err)
			return "fail", err
		}
	}
	return "success", nil
}

func CreateContainer(taskName string, imageName string, workflowName string, containerDataDir []string) (string, error) {
	if !CheckImageExist(imageName) {
		ImagePull(imageName)
	}

	utils.CreateFolder(fmt.Sprintf("/data/menet/workflow/%s/%s", workflowName, taskName))

	// 创建容器的配置
	config := &container.Config{
		Image: imageName,
		Labels: map[string]string{
			"app":       taskName,
			"namespace": "menet",
		},
	}

	mounts := []mount.Mount{}

	for _, targetPath := range containerDataDir {
		utils.CreateFolder("/data/menet/task/" + taskName + targetPath)
		m := mount.Mount{
			Type:   mount.TypeBind,
			Source: "/data/menet/task/" + taskName + targetPath, //表示要挂载到容器内部的宿主机目录。
			Target: targetPath,                                  //容器内的文件夹。
		}

		mounts = append(mounts, m)
	}

	hostConfig := &container.HostConfig{
		Mounts: mounts,
	}

	resp, err := utils.DockerCli.ContainerCreate(context.Background(), config, hostConfig, nil, nil, taskName)
	if err != nil {
		log.Println("Container creation error:", err)
		return "fail", err
	}

	// 启动容器
	if err := utils.DockerCli.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{}); err != nil {
		log.Println("Container startup error:", err)
		return "fail", err
	}

	return "success", nil
}

func DeleteContainer(containerName string, isDeleteVolume bool) error {
	if CheckContainerExist(containerName) {
		if CheckContainerStatus(containerName) == "running" { //容器正在运行,先停止
			err := StopContainer(containerName)
			if err != nil {
				return err
			}
		}
	} else {
		return fmt.Errorf("The container does not exist")
	}

	err := utils.DockerCli.ContainerRemove(
		context.Background(),
		containerName,
		types.ContainerRemoveOptions{
			Force: false,
		},
	)

	if err != nil {
		log.Println("Failed to delete container:", err)
		return err
	}

	if isDeleteVolume { //需要删除卷
		err := utils.DeleteFolder("/data/menet/task/" + containerName)
		return err
	}

	return nil
}

func StopContainer(containerName string) error {
	err := utils.DockerCli.ContainerStop(context.Background(), containerName, container.StopOptions{})
	if err != nil {
		log.Println("Container stop failed:", err)
		return err
	}

	return nil
}

func StartContainer(containerName string) error {
	err := utils.DockerCli.ContainerStart(context.Background(), containerName, types.ContainerStartOptions{})
	if err != nil {
		log.Println("Container startup failed:", err)
		return err
	}

	return nil
}

func CheckContainerExist(containerName string) bool {
	_, err := utils.DockerCli.ContainerInspect(context.Background(), containerName)
	if err == nil {
		return true
	} else if clientv1.IsErrNotFound(err) {
		return false
	} else {
		log.Println("Find Container Error:", err)
	}
	return false
}

func CheckContainerStatus(containerName string) string {
	containerInfo, err := utils.DockerCli.ContainerInspect(context.Background(), containerName)
	if err != nil {
		log.Println(err)
	}

	if containerInfo.State.Running {
		return "running"
	} else if containerInfo.State.Status == "exited" {
		fmt.Println("Container exited normally")
		return "exited"
	} else {
		return "error"
	}
}

//func UpdateContainer(containername string) {
//	dockerClient := client.GetDockerCli()
//
//	err := dockerClient.ContainerStop(context.Background(), containername, container.StopOptions{Signal: containername})
//	if err != nil {
//		fmt.Printf("停止容器失败：%v\n", err)
//		return
//	}
//
//	fmt.Printf("成功停止容器：%s\n", containername)
//	resp, err := dockerClient.ContainerInspect(
//		context.Background(),
//		containername,
//	)
//	if err != nil {
//		fmt.Printf("获取容器信息失败：%v\n", err)
//		return
//	}
//
//	// 获取原始的资源限制和请求配置
//	resources := resp.HostConfig.Resources
//	limits := resources.CPUQuota
//	newLimits := map[string]types.ConfigCreateResponse{
//		"cpu":    "50",   // 新的CPU限制，以百分比形式表示
//		"memory": "256m", // 新的内存限制
//	}
//
//	newRequests := map[string]types.NanoCPUPercent{
//		"cpu":    20,   // 新的CPU请求，以百分比形式表示
//		"memory": 128m, // 新的内存请求
//	}
//
//	// 创建更新后的容器配置
//	containerConfig := &container.UpdateConfig{
//		CPULimits:     newLimits["cpu"],
//		MemoryLimit:   newLimits["memory"],
//		CPUShares:     newRequests["cpu"],
//		MemoryRequest: newRequests["memory"],
//	}
//
//	// 更新容器
//	err = dockerClient.ContainerUpdate(
//		context.Background(),
//		containername,
//		resp.Container.Config,
//		containerConfig,
//	)
//	if err != nil {
//		fmt.Printf("更新容器失败：%v\n", err)
//		return
//	}
//
//	fmt.Printf("成功更新容器：%s 的资源限制和请求\n", containername)
//}
