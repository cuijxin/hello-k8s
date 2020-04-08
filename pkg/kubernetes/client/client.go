package client

import (
	"path/filepath"

	pgsqlClientset "github.com/cuijxin/postgres-operator-atom/pkg/generated/clientset/versioned"
	"github.com/lexkong/log"
	mysqlClientset "github.com/oracle/mysql-operator/pkg/generated/clientset/versioned"
	redisClientset "github.com/spotahome/redis-operator/client/k8s/clientset/versioned"
	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func New() (kubernetes.Interface, error) {
	config, err := getKubernetesConfig()
	if err != nil {
		return nil, err
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func NewApiExtensionsClient() (*apiextensionsclientset.Clientset, error) {
	config, err := getKubernetesConfig()
	if err != nil {
		return nil, err
	}

	client, err := apiextensionsclientset.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func NewPostgresClientSet() (pgsqlClientset.Interface, error) {
	config, err := getKubernetesConfig()
	if err != nil {
		return nil, err
	}

	client, err := pgsqlClientset.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func NewRedisClientSet() (redisClientset.Interface, error) {
	config, err := getKubernetesConfig()
	if err != nil {
		return nil, err
	}

	client, err := redisClientset.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func NewMySQLClientSet() (mysqlClientset.Interface, error) {
	config, err := getKubernetesConfig()
	if err != nil {
		return nil, err
	}

	client, err := mysqlClientset.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func getKubernetesConfig() (*rest.Config, error) {
	var kubeconfig *string
	var tmp string
	if home := homedir.HomeDir(); home != "" {
		log.Infof("home dir is:%v", home)
		tmp = filepath.Join(home, ".kube", "config")
		log.Infof("config path is:%v", tmp)
		kubeconfig = &tmp
	}

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		return nil, err
	}

	return config, nil
}
