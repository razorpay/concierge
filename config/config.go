package config

import (
	"flag"
	"os"
	"path/filepath"
	"strconv"

	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	"github.com/joho/godotenv"
)

//DBConfig ...
var DBConfig DatabaseConfig

//KubeClients ...
var KubeClients map[string]KubenetesClientSet

//KubeConfig ...
var KubeConfig *string

//LoadConfig ...
func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Error("Error loading .env file")
	}
	initilizeDBConfig()
	initilizeKubeConfigFromFile()
	InitilizeKubeConfig()
}

func initilizeDBConfig() {
	maxidleconns, _ := strconv.Atoi(os.Getenv("DB_MAX_IDLE_CONN"))
	maxopenconns, _ := strconv.Atoi(os.Getenv("DB_MAX_OPEN_CONN"))
	DBConfig = DatabaseConfig{
		Host:         os.Getenv("DB_HOST"),
		DBName:       os.Getenv("DB_DATABASE"),
		DBUsername:   os.Getenv("DB_USERNAME"),
		DBPassword:   os.Getenv("DB_PASSWORD"),
		DBPort:       os.Getenv("DB_PORT"),
		DBDatabase:   os.Getenv("DB_DATABASE"),
		MaxIdleConns: maxidleconns,
		MaxOpenConns: maxopenconns,
	}
}

//InitilizeKubeConfig ...
func InitilizeKubeConfig() {
	var err error
	var config *rest.Config
	var clientset *kubernetes.Clientset
	var contexts = []string{"prod-green", "prod-blue", "stage-white"}
	KubeClients = make(map[string]KubenetesClientSet)

	if os.Getenv("APP_ENV") == "dev" {
		for _, context := range contexts {
			log.Info(context)
			config, err = customBuildConfigFromFlags(context, *KubeConfig)
			if err != nil {
				log.Error(err)
			}
			clientset, _ = kubernetes.NewForConfig(config)
			kubeclient := KubenetesClientSet{
				ClientSet: clientset,
			}
			log.Info(kubeclient)
			KubeClients[context] = kubeclient
		}

	} else {
		config = initilizeKubeConfigFromServiceAccount()
		clientset, err = kubernetes.NewForConfig(config)
		if err != nil {
			log.Error(err)
		}
		kubeclient := KubenetesClientSet{
			ClientSet: clientset,
		}
		KubeClients = map[string]KubenetesClientSet{
			"prod-green": kubeclient,
		}
	}
}

func initilizeKubeConfigFromFile() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()
	KubeConfig = kubeconfig
	// config, _ := clientcmd.LoadFromFile("/Users/ankitjain/.kube/config")
	// log.Info(config)
}

func initilizeKubeConfigFromServiceAccount() *rest.Config {
	var err error
	var config *rest.Config
	config, err = rest.InClusterConfig()
	if err != nil {
		log.Error(err)
	}
	return config
}

func customBuildConfigFromFlags(context, kubeconfigPath string) (*rest.Config, error) {
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfigPath},
		&clientcmd.ConfigOverrides{
			CurrentContext: context,
		}).ClientConfig()
}
