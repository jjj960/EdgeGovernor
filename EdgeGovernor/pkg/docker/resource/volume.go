package resource

import (
	"EdgeGovernor/pkg/utils"
	"context"
	"fmt"
	"github.com/docker/docker/api/types/volume"
	"log"
	"regexp"
)

func CreateVolume(taskName string, nfsServerIP string, containerDir []string) error {
	volumePrefix := taskName + "-volume" // 卷名称前缀

	for i, _ := range containerDir {
		volumeName := fmt.Sprintf("%s%d", volumePrefix, i) //根据需要数据持久化的文件夹的个数,卷名递增,volume0,volume1.....
		volumePath := "/data/menet/" + taskName + "/" + volumeName

		err := utils.CreateFolder(volumePath)
		if err != nil {
			return err
		}

		nfsOptions := map[string]string{
			"type":   "nfs",
			"o":      "addr=" + nfsServerIP + ",rw,nolock",
			"device": ":/data/menet/" + taskName + "/" + volumeName,
		}

		options := volume.CreateOptions{
			Name:       volumeName,
			Driver:     "local",
			DriverOpts: nfsOptions,
		}

		_, err = utils.DockerCli.VolumeCreate(context.Background(), options)
		if err != nil {
			return err
		}
	}

	return nil
}

func DeleteVolume(taskVolumeList []string) error {
	if len(taskVolumeList) == 0 {
		return fmt.Errorf("There are no volumes to delete")
	}

	for _, name := range taskVolumeList {
		err := utils.DockerCli.VolumeRemove(context.Background(), name, true)
		if err != nil {
			return fmt.Errorf("Delete volume failed: %s", err)
		}
	}

	return nil
}

func GetVolumeList(taskName string) []string {
	volumeList, err := utils.DockerCli.VolumeList(context.Background(), volume.ListOptions{})
	if err != nil {
		log.Println("Failed to obtain volume list:", err)
	}
	var taskVolumeList []string
	for _, volume1 := range volumeList.Volumes {
		match, _ := regexp.MatchString(taskName+"-volume\\d+", volume1.Name)
		if match {
			taskVolumeList = append(taskVolumeList, volume1.Name)
		}
	}

	return taskVolumeList
}
