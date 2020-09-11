package config

import (
	"flag"
	"os"
	"path/filepath"
	"strconv"
	"strings"

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

//Contexts ...
var Contexts = []string{}

//CSRFConfig ...
var CSRFConfig CSRF

//AppCfg ...
var AppCfg Application

//LookerConfig
var LookerConfig Looker

//LoadConfig ...
func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Error("Error loading .env file")
	}

	initilizeDBConfig()
	initilizeKubeContext()
	initilizeKubeConfigFromFile()
	initilizeKubeConfig()
	initilizeLogging()
	initilizeCSRFConfig()
	initilizeAppConfig()
	initializeLookerConfig()
}

func initilizeDBConfig() {
	maxidleconns, _ := strconv.Atoi(os.Getenv("DB_MAX_IDLE_CONN"))
	maxopenconns, _ := strconv.Atoi(os.Getenv("DB_MAX_OPEN_CONN"))
	dblogmode, _ := strconv.ParseBool(os.Getenv("DB_LOG_MODE"))
	DBConfig = DatabaseConfig{
		Host:         os.Getenv("DB_HOST"),
		DBName:       os.Getenv("DB_DATABASE"),
		DBUsername:   os.Getenv("DB_USERNAME"),
		DBPassword:   os.Getenv("DB_PASSWORD"),
		DBPort:       os.Getenv("DB_PORT"),
		DBDatabase:   os.Getenv("DB_DATABASE"),
		MaxIdleConns: maxidleconns,
		MaxOpenConns: maxopenconns,
		DBLogMode:    dblogmode,
	}
}

func initilizeKubeContext() {
	k8sContexts := strings.Split(os.Getenv("KUBE_CONTEXTS"), ",")
	for _, context := range k8sContexts {
		Contexts = append(Contexts, context)
	}
}

func initilizeKubeConfig() {
	var err error
	var config *rest.Config
	var clientset *kubernetes.Clientset
	KubeClients = make(map[string]KubenetesClientSet)

	if os.Getenv("AUTH_TYPE") == "KUBECONFIG" {
		for _, context := range Contexts {
			config, err = customBuildConfigFromFlags(context, *KubeConfig)
			if err != nil {
				log.Fatal(err)
			}
			clientset, err = kubernetes.NewForConfig(config)
			if err != nil {
				log.Fatal(err)
			}
			kubeclient := KubenetesClientSet{
				ClientSet: clientset,
			}
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
			"context": kubeclient,
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

func initilizeLogging() {
	log.SetFormatter(&log.JSONFormatter{})
	loglevel := os.Getenv("LOG_LEVEL")
	if loglevel == "debug" {
		log.SetLevel(log.DebugLevel)

	} else {
		log.SetLevel(log.InfoLevel)
	}
	log.Debug("Logging in debug mode.")
}

func initilizeCSRFConfig() {
	secure := getEnv(os.Getenv("CSRF_SECURE"), false).(bool)
	CSRFConfig = CSRF{
		AuthKey: os.Getenv("CSRF_AUTH_KEY"),
		Secure:  secure,
	}
}

func initilizeAppConfig() {
	maxExpiry := getEnv(os.Getenv("APP_MAX_EXPIRY"), 32400).(int)
	cookieSecure := getEnv(os.Getenv("COOKIE_SECURE"), false).(bool)
	cookieHTTPOnly := getEnv(os.Getenv("COOKIE_HTTPONLY"), true).(bool)
	appPort := getEnv(os.Getenv("APP_PORT"), "8990").(string)
	appIP := getEnv(os.Getenv("APP_IP"), "0.0.0.0").(string)
	AppCfg = Application{
		Name:           os.Getenv("NAME"),
		Mode:           os.Getenv("APP_ENV"),
		ListenIP:       appIP,
		ListenPort:     appPort,
		MaxExpiry:      maxExpiry,
		CookieHTTPOnly: cookieHTTPOnly,
		CookieSecure:   cookieSecure,
	}
}

func initializeLookerConfig() {
	baseUrl := getEnv(os.Getenv("LOOKER_CLIENT_ID"), "https://looker.razorpay.com:19999/").(string)
	clientId := getEnv(os.Getenv("LOOKER_CLIENT_ID"), "").(string)
	clientSecret := getEnv(os.Getenv("LOOKER_CLIENT_SECRET"), "").(string)

	LookerConfig = Looker{
		BaseUrl:      baseUrl,
		ClientId:     clientId,
		ClientSecret: clientSecret,
	}
}

func getEnv(value string, x interface{}) interface{} {
	if value != "" {
		switch v := x.(type) {
		case string:
			return value
		case bool:
			val, _ := strconv.ParseBool(value)
			return val
		case int:
			val, _ := strconv.Atoi(value)
			return val
		default:
			return v
		}

	}
	return x
}
