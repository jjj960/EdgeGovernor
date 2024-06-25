package resource

import (
	"MENet/pkg/docker"
	"MENet/pkg/utils"
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"log"
	"testing"
)

func Test(t *testing.T) {
	docker.GetDockerCli()
	//contain := []string{"/etc/apt", "/root"}
	//CreateContainer1("test", "nginx:latest", 50, 100, contain)

}

func CreateContainer1(taskName string, imageName string, cpuSize int64, memSize int64, containerDataDir []string) {
	if !CheckImageExist(imageName) {
		ImagePull(imageName)
	}

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

	resp, err := docker.DockerCli.ContainerCreate(context.Background(), config, hostConfig, nil, nil, taskName)
	if err != nil {
		log.Println("Container creation error:", err)
	}

	// 启动容器
	if err := docker.DockerCli.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{}); err != nil {
		log.Println("Container startup error:", err)
	}
}
