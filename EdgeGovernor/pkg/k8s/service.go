package k8s

import (
	"EdgeGovernor/pkg/utils"
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func CreateService(serviceName string, podName string, sourcePort int32, targetPort int32) {
	// 指定Service的名称和Namespace
	namespace := "menet"

	// 创建一个Service对象
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceName,
			Namespace: namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app": podName,
			},
			Ports: []corev1.ServicePort{
				{
					Name:     "http",
					Port:     targetPort, //暴露的端口号
					Protocol: corev1.ProtocolTCP,
					TargetPort: intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: sourcePort, // Pod容器中暴露的端口号
					},
				},
			},
		},
	}

	// 创建Service
	result, err := utils.K8sClientset.CoreV1().Services(namespace).Create(context.TODO(), service, metav1.CreateOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Service %s created\n", result.GetName())
}

func DeleteService(serviceName string) {
	// 指定Service的名称和Namespace
	namespace := "menet"
	// 删除Service
	err := utils.K8sClientset.CoreV1().Services(namespace).Delete(context.TODO(), serviceName, metav1.DeleteOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Service %s deleted\n", serviceName)
}
