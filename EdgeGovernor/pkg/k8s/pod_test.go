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
	"testing"
)

func Test(t *testing.T) {
	utils.GetK8sCli()
	dataDit := []string{"/root", "/etc/apt"}
	_, err := CreatePod1("master", "test", "nginx:latest", "50", "50", 0, "", dataDit)
	if err != nil {
		fmt.Println(err)
	}
	//DeletePod("test", false)
}

func CreatePod1(nodename string, podname string, podimage string, cpu string, mem string, needPersistence int, storage string, dataDir []string) (string, error) {
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
}
