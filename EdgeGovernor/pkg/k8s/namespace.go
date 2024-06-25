package k8s

import (
	"EdgeGovernor/pkg/utils"
	"context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
)

func CreateNamespace() {
	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "menet",
		},
	}

	createdNamespace, err := utils.K8sClientset.CoreV1().Namespaces().Create(context.TODO(), namespace, metav1.CreateOptions{})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Namespace %s created successfully\n", createdNamespace.Name)
}

func DeleteNamespace() {
	err := utils.K8sClientset.CoreV1().Namespaces().Delete(context.TODO(), "menet", metav1.DeleteOptions{})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Namespace menet created successfully\n")
}
