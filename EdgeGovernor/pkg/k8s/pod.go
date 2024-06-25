package k8s

import (
	"EdgeGovernor/pkg/utils"
	"context"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
)

func CreateResourcePod(nodename string, podname string, podimage string, cpu string, mem string, needPersistence int, storage string, dataDir []string) (string, error) { //pod的创建
	var replicas int32 = 1
	switch needPersistence {
	case 1:
		var replicas int32 = 1

		utils.CreateFolder("/data/menet/task/" + podname)

		// 创建 Deployment 对象
		deployment := &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      podname,
				Namespace: "menet",
			},
			Spec: appsv1.DeploymentSpec{
				Replicas: &replicas,
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{
						"app": podname + "-app",
					},
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							"app": podname + "-app",
						},
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  podname + "-container",
								Image: podimage,
								Resources: corev1.ResourceRequirements{
									Requests: corev1.ResourceList{
										corev1.ResourceCPU:    resource.MustParse(cpu + "m"),
										corev1.ResourceMemory: resource.MustParse(mem + "Mi"),
									},
									Limits: corev1.ResourceList{
										corev1.ResourceCPU:    resource.MustParse(cpu + "m"),
										corev1.ResourceMemory: resource.MustParse(mem + "Mi"),
									},
								},
								VolumeMounts: getVolumeMounts(dataDir),
							},
						},
						Volumes: getVolumes(podname, dataDir),
					},
				},
			},
		}

		// 创建 Deployment
		createdDeployment, err := utils.K8sClientset.AppsV1().Deployments("menet").Create(context.Background(), deployment, metav1.CreateOptions{})
		if err != nil {
			log.Println(err)
			return "fail", err
		}
		log.Printf("Deployment %s created\n", createdDeployment.Name)
		fmt.Printf("Task %s successfully deployed on node %s, allocating CPU %s m, allocating memory %s MB", podname, nodename, cpu, mem)
		return "success", nil

	case 0: //任务不需要数据持久化

		deployment := &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      podname,
				Namespace: "menet",
			},
			Spec: appsv1.DeploymentSpec{
				Replicas: &replicas, // 副本数
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{
						"app": podname,
					},
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							"app": podname,
						},
					},
					Spec: corev1.PodSpec{
						NodeName: nodename,
						Containers: []corev1.Container{
							{
								Name:  podname,
								Image: podimage,
								Resources: corev1.ResourceRequirements{
									Limits: corev1.ResourceList{
										corev1.ResourceCPU:    resource.MustParse(cpu + "m"),
										corev1.ResourceMemory: resource.MustParse(mem + "Mi"),
									},
									Requests: corev1.ResourceList{
										corev1.ResourceCPU:    resource.MustParse(cpu + "m"),
										corev1.ResourceMemory: resource.MustParse(mem + "Mi"),
									},
								},
							},
						},
					},
				},
			},
		}
		deployment, err := utils.K8sClientset.AppsV1().Deployments("menet").Create(context.TODO(), deployment, metav1.CreateOptions{})
		if err != nil {
			log.Println(err)
			return "fail", err
		}
		log.Printf("Deployment %s created\n", deployment.Name)
		fmt.Printf("Task %s successfully deployed on node %s, allocating CPU %s m, allocating memory %s MB, allocating disk %s MB", podname, nodename, cpu, mem, storage)
		return "success", nil
	}
	return "", nil
}

func CreatePod(nodename string, podname string, podimage string, workflowName string, dataDir []string) (string, error) { //pod的创建
	var replicas int32 = 1

	utils.CreateFolder(fmt.Sprintf("/data/menet/workflow/%s/%s", workflowName, podname))

	// 创建 Deployment 对象
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podname,
			Namespace: "menet",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": podname + "-app",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": podname + "-app",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:         podname + "-container",
							Image:        podimage,
							VolumeMounts: getVolumeMounts(dataDir),
						},
					},
					Volumes: getVolumes(podname, dataDir),
				},
			},
		},
	}

	// 创建 Deployment
	createdDeployment, err := utils.K8sClientset.AppsV1().Deployments("menet").Create(context.Background(), deployment, metav1.CreateOptions{})
	if err != nil {
		log.Println(err)
		return "fail", err
	}
	log.Printf("Deployment %s created\n", createdDeployment.Name)
	fmt.Printf("Job %s successfully deployed on node %s", podname, nodename)
	return "success", nil
}

func getVolumes(podname string, dataDir []string) []corev1.Volume { //宿主机
	volumes := []corev1.Volume{}
	for i, path := range dataDir {
		utils.CreateFolder("/data/menet/task/" + podname + path)
		volume := corev1.Volume{
			Name: fmt.Sprintf("volume-%d", i+1),
			VolumeSource: corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: "/data/menet/task/" + podname + path,
				},
			},
		}
		volumes = append(volumes, volume)
	}
	return volumes
}

func getVolumeMounts(dataDir []string) []corev1.VolumeMount { //pod内部

	volumeMounts := []corev1.VolumeMount{}
	for i, path := range dataDir {
		volumeMount := corev1.VolumeMount{
			Name:      fmt.Sprintf("volume-%d", i+1),
			MountPath: path,
		}
		volumeMounts = append(volumeMounts, volumeMount)
	}
	return volumeMounts
}

func DeletePod(podName string, shouldDeleteVolumeAndPVC int) { //pod的删除,并判断是否是执行完毕或者迁移
	deleteOptions := metav1.DeleteOptions{}
	// 删除Pod
	err := utils.K8sClientset.AppsV1().Deployments("menet").Delete(context.TODO(), podName, deleteOptions)
	if err != nil {
		panic(err.Error())
	}

	// 如果需要删除Volume和PVC，则执行下面的代码
	if shouldDeleteVolumeAndPVC == 1 {
		volumes, err := utils.K8sClientset.CoreV1().Pods("menet").Get(context.TODO(), podName, metav1.GetOptions{})
		if err != nil {
			log.Println(err)
		}

		for _, volume := range volumes.Spec.Volumes {
			if volume.PersistentVolumeClaim != nil {
				pvcName := volume.PersistentVolumeClaim.ClaimName
				err = utils.K8sClientset.CoreV1().PersistentVolumeClaims("menet").Delete(context.TODO(), pvcName, deleteOptions)
				if err != nil {
					log.Println(err)
				}
			}
		}
	}
}

func UpdatePodNode(podName string, targetNode string) (string, error) {
	// 获取原始的Deployment对象
	deployment, err := utils.K8sClientset.AppsV1().Deployments("menet").Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		log.Println(err)
		return "fail", err
	}

	// 更新Deployment的PodTemplateSpec的NodeName属性
	deployment.Spec.Template.Spec.NodeName = targetNode

	// 更新Deployment对象
	_, err = utils.K8sClientset.AppsV1().Deployments("menet").Update(context.TODO(), deployment, metav1.UpdateOptions{})
	if err != nil {
		log.Println(err)
		return "fail", err
	}
	return "success", nil
}

func SelectPod(namespace string) map[string]corev1.Pod { //获取pod的信息
	// 构建列表选项
	listOptions := metav1.ListOptions{}

	// 调用API获取指定Namespace下的所有Pod
	pods, err := utils.K8sClientset.CoreV1().Pods(namespace).List(context.TODO(), listOptions)
	if err != nil {
		panic(err.Error())
	}

	// 将Pod信息存储到Map中
	podMap := make(map[string]corev1.Pod)
	for _, pod := range pods.Items {
		podMap[pod.Name] = pod
	}

	//// 打印每个Pod的信息
	//for podName, pod := range podMap {
	//	fmt.Printf("Pod Name: %s\n", podName)
	//	fmt.Printf("Pod Status: %s\n", pod.Status.Phase)
	//	fmt.Println("-------------------------")
	//}

	return podMap
}
