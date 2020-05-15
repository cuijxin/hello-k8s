package client

import (
	"path/filepath"

	"github.com/lexkong/log"
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
