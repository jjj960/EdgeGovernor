package resource

import (
	"EdgeGovernor/pkg/utils"
	"context"
	"github.com/docker/docker/api/types"
	"io"
	"log"
	"os"
)

func ImagePull(name string) string {
	out, err := utils.DockerCli.ImagePull(context.Background(), name, types.ImagePullOptions{})
	if err != nil {
		log.Println(err)
		return err.Error()
	}

	defer out.Close()

	io.Copy(os.Stdout, out)

	return "success pull"
}

func DeleteImage(imageName string) error {
	_, err := utils.DockerCli.ImageRemove(context.Background(), imageName, types.ImageRemoveOptions{})
	if err != nil {
		return err
	}

	return nil
}

func CheckImageExist(imageName string) bool { //检查该镜像是否存在
	// 获取主机上的所有镜像
	images, err := utils.DockerCli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		log.Println("Failed to obtain mirror list:", err)
	}

	// 遍历镜像列表，检查是否存在指定的镜像
	for _, image := range images {
		for _, tag := range image.RepoTags {
			if tag == imageName {
				return true
			}
		}
	}
	return false
}
