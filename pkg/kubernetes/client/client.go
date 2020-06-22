package client

import (
	"path/filepath"

	mysql5Clientset "github.com/cuijxin/mysql-operator/pkg/generated/clientset/versioned"
	pgClientset "github.com/cuijxin/postgres-operator-atom/pkg/generated/clientset/versioned"
	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/klog"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var MyClient *HelloK8SClient

type HelloK8SClient struct {
	KubeConfig          *rest.Config
	K8sClientset        kubernetes.Interface
	ApiExtensionsClient *apiextensionsclientset.Clientset

	Mysql5Client mysql5Clientset.Interface
	PgClient     pgClientset.Interface
}

func (c *HelloK8SClient) InitHelloK8SClient() {
	config, err := getKubernetesConfig()
	if err != nil {
		panic(err)
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	apiExtensionsClient, err := apiextensionsclientset.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	mysql5client, err := mysql5Clientset.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	postgresqlClient, err := pgClientset.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	MyClient = &HelloK8SClient{
		KubeConfig:          config,
		K8sClientset:        client,
		ApiExtensionsClient: apiExtensionsClient,
		Mysql5Client:        mysql5client,
		PgClient:            postgresqlClient,
	}
}

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

func getKubernetesConfig() (*rest.Config, error) {
	var kubeconfig *string
	var tmp string
	if home := homedir.HomeDir(); home != "" {
		klog.Infof("home dir is:%v", home)
		tmp = filepath.Join(home, ".kube", "config-tencent")
		klog.Infof("config path is:%v", tmp)
		kubeconfig = &tmp
	}

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		return nil, err
	}

	return config, nil
}
