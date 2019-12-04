package config

import (
	"flag"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	"github.com/joho/godotenv"
)

//DBConfig ...
var DBConfig DatabaseConfig

//KubeClient ...
var KubeClient KubenetesClientSet

//LoadConfig ...
func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Error("Error loading .env file")
	}
	initilizeDBConfig()
	initilizeKubeConfig()
}

func initilizeDBConfig() {
	DBConfig = DatabaseConfig{
		Host:       os.Getenv("DB_HOST"),
		DBName:     os.Getenv("DB_DATABASE"),
		DBUsername: os.Getenv("DB_USERNAME"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBPort:     os.Getenv("DB_PORT"),
		DBDatabase: os.Getenv("DB_DATABASE"),
	}
}

func initilizeKubeConfig() {
	var err error
	var config *rest.Config
	var clientset *kubernetes.Clientset

	log.Info(os.Getenv("APP_ENV"))
	if os.Getenv("APP_ENV") == "dev" {
		var kubeconfig *string
		if home := homedir.HomeDir(); home != "" {
			kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		} else {
			kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
		}
		flag.Parse()
		config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			log.Error(err)
		}
	} else {
		config, err = rest.InClusterConfig()
		if err != nil {
			log.Error(err)
		}
	}
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		log.Error(err)
	}
	KubeClient = KubenetesClientSet{
		ClientSet: clientset,
	}
}
