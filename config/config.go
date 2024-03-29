package config

import (
	"flag"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"concierge/constants"

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

//CronConfig ...
var CronConfig Cron

//AppCfg ...
var AppCfg Application

//LookerConfig
var LookerConfig Looker

//RedisDBConfig ...
var RedisDBConfig RedisConfig

var MutexPrefix string

var AWSS3Buckets string

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
	initilizeMutexConfig()
	initializeCronConfig()
	initializeAWSS3BucketConfig()
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

func initilizeMutexConfig() {
	database := getEnv(constants.MutexDatabase, 0).(int)
	maxIdle := getEnv(constants.MutexMaxIdle, 10).(int)
	maxActive := getEnv(constants.MutexMaxActive, 100).(int)

	RedisDBConfig = RedisConfig{
		Host:      os.Getenv(constants.MutexHost),
		Database:  database,
		Password:  os.Getenv(constants.MutexPassword),
		Port:      os.Getenv(constants.MutexPort),
		MaxIdle:   maxIdle,
		MaxActive: maxActive,
	}

	MutexPrefix = getEnv(os.Getenv("Name"), constants.MutexPrefix).(string)
}

func initilizeKubeContext() {
	k8sContexts := strings.Split(os.Getenv("KUBE_CONTEXTS"), ",")
	Contexts = append(Contexts, k8sContexts...)
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

func initializeCronConfig() {
	CronConfig = Cron{
		CronUsername: getEnv(os.Getenv("CRON_USERNAME"), "cron").(string),
		CronPassword: os.Getenv("CRON_PASSWORD"),
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
	baseUrl := getEnv(os.Getenv("LOOKER_BASE_URL"), "").(string)
	clientId := getEnv(os.Getenv("LOOKER_CLIENT_ID"), "").(string)
	clientSecret := getEnv(os.Getenv("LOOKER_CLIENT_SECRET"), "").(string)
	datumHostname := getEnv(os.Getenv("DATUM_HOSTNAME"), "").(string)
	datumAuthSecret := getEnv(os.Getenv("DATUM_AUTH_SECRET"), "").(string)

	isEnabled := false

	if baseUrl != "" && clientId != "" && clientSecret != "" && datumHostname != "" && datumAuthSecret != "" {
		isEnabled = true
	}

	LookerConfig = Looker{
		BaseUrl:         baseUrl,
		ClientId:        clientId,
		ClientSecret:    clientSecret,
		IsEnabled:       isEnabled,
		DatumHostname:   datumHostname,
		DatumAuthSecret: datumAuthSecret,
	}
}

func initializeAWSS3BucketConfig() {
	AWSS3Buckets = strings.TrimSpace(getEnv(os.Getenv("AWS_S3_BUCKETS"), "").(string))
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
