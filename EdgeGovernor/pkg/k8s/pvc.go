package k8s

import (
	"EdgeGovernor/pkg/constants"
	"EdgeGovernor/pkg/utils"
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/storage/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
	"log"
)

func CreatePV(size string) {
	pvName := "menet"
	storageSize := size + "Gi"
	storageClassName := "nfs"
	accessModes := []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce}
	nfsServer := constants.Hostname
	nfsPath := "/data/menet"

	pv := &corev1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name: pvName,
		},
		Spec: corev1.PersistentVolumeSpec{
			StorageClassName: storageClassName,
			Capacity: corev1.ResourceList{
				corev1.ResourceStorage: resource.MustParse(storageSize),
			},
			AccessModes: accessModes,
			PersistentVolumeSource: corev1.PersistentVolumeSource{
				NFS: &corev1.NFSVolumeSource{
					Server: nfsServer,
					Path:   nfsPath,
				},
			},
		},
	}

	// 创建PV对象
	_, err := utils.K8sClientset.CoreV1().PersistentVolumes().Create(context.TODO(), pv, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	log.Println("PV", pvName, "已成功创建")
}

func CreatePVC(pvcName string, storageSize string) error {

	storageCapacity := storageSize + "Gi" // 替换为你需要的存储容量，例如 "1Gi"

	// 创建一个PVC对象
	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name: pvcName,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{
				corev1.ReadWriteMany,
			},
			StorageClassName: pointer.String("nfs"),
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse(storageCapacity),
				},
			},
		},
	}

	// 创建PVC
	result, err := utils.K8sClientset.CoreV1().PersistentVolumeClaims("menet").Create(context.TODO(), pvc, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	fmt.Printf("PersistentVolumeClaim %s created\n", result.GetName())
	return nil
}

func GetPVCList() (*corev1.PersistentVolumeClaimList, error) {
	pvcList, err := utils.K8sClientset.CoreV1().PersistentVolumeClaims("menet").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Println(err)
	}
	return pvcList, err
}

func GetPVC(pvcName string) (*corev1.PersistentVolumeClaim, error) {
	pvcInfo, err := utils.K8sClientset.CoreV1().PersistentVolumeClaims("menet").Get(context.TODO(), pvcName, metav1.GetOptions{})
	if err != nil {
		log.Println(err)
	}
	return pvcInfo, err
}

func DeletePVC(pvcName string) string {
	err := utils.K8sClientset.CoreV1().PersistentVolumeClaims("menet").Delete(context.TODO(), pvcName, metav1.DeleteOptions{})
	if err != nil {
		log.Println(err)
	}
	return "Successfully delete pvc"
}

func CreateStorageClass(storageClassName string, nfsServer string, nfsPath string) {
	// 指定StorageClass的名称和NFS配置
	volumeBindingMode := "WaitForFirstConsumer"
	// 创建一个NFS类型的StorageClass对象
	storageClass := &v1.StorageClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: storageClassName,
		},
		Provisioner:          "kubernetes.io/nfs",
		VolumeBindingMode:    (*v1.VolumeBindingMode)(&volumeBindingMode),
		AllowVolumeExpansion: pointer.Bool(true),
		MountOptions:         []string{"vers=4.1"},
		Parameters: map[string]string{
			"server": nfsServer,
			"path":   nfsPath,
		},
	}

	// 创建StorageClass
	_, err := utils.K8sClientset.StorageV1().StorageClasses().Create(context.TODO(), storageClass, metav1.CreateOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("StorageClass %s created with NFS provisioner\n", storageClassName)
}

func UpdatePVC(pvcName string, newStorageCapacity string) { //动态更新pvc容量,只支持扩容!!!
	// 指定PVC的名称、Namespace和新的存储容量
	namespace := "menet"

	// 获取现有的PVC
	existingPVC, err := utils.K8sClientset.CoreV1().PersistentVolumeClaims(namespace).Get(context.TODO(), pvcName, metav1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}
	//capacity := existingPVC.Status.Capacity.Storage()      //比较容量,防止pvc缩容
	//if int(capacity) >= int(newStorageCapacity) {
	//
	//}
	// 更新PVC的存储容量
	existingPVC.Spec.Resources.Requests[corev1.ResourceStorage] = resource.MustParse(newStorageCapacity + "Gi") // 替换为新的存储容量，例如 "2Gi"

	// 执行PVC的更新操作
	_, err = utils.K8sClientset.CoreV1().PersistentVolumeClaims(namespace).Update(context.TODO(), existingPVC, metav1.UpdateOptions{})
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("PersistentVolumeClaim %s updated with new storage capacity\n", pvcName)
}
