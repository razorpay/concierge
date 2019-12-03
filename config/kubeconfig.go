package config

import (
	"k8s.io/client-go/kubernetes"
)

//KubenetesClientSet ...
type KubenetesClientSet struct {
	ClientSet *kubernetes.Clientset
}
