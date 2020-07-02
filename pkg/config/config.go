package config

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/gofrs/flock"
	"github.com/spf13/viper"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/kube"
	"helm.sh/helm/v3/pkg/repo"
	"k8s.io/client-go/util/homedir"
	"k8s.io/klog"
	"sigs.k8s.io/yaml"
)

type HelmConfig struct {
	HelmRepos []*repo.Entry
}

var (
	Settings = cli.New()
	HelmConf = &HelmConfig{}
)

type Config struct {
	Name string
}

func Init(cfg string) error {
	c := Config{
		Name: cfg,
	}

	// 初始化配置文件
	if err := c.initConfig(); err != nil {
		return err
	}

	if Settings.KubeConfig == "" {
		if home := homedir.HomeDir(); home != "" {
			Settings.KubeConfig = filepath.Join(home, ".kube", "config-tencent")
		}
	}

	// 初始化Helm Repo
	if err := c.initHelmRepository(); err != nil {
		return err
	}

	// 初始化日志包
	// c.initLog()

	// 监控配置文件变化并热加载程序
	c.watchConfig()

	return nil
}

func (c *Config) initConfig() error {
	if c.Name != "" {
		viper.SetConfigFile(c.Name) // 如果指定了配置文件，则解析指定的配置文件
	} else {
		viper.AddConfigPath("conf") // 如果没有指定配置文件，则解析默认的配置文件
		viper.SetConfigName("config")
	}
	viper.SetConfigType("yaml")     // 设置配置文件格式为YAML
	viper.AutomaticEnv()            // 读取匹配的环境变量
	viper.SetEnvPrefix("APISERVER") // 读取环境变量的前缀为APISERVER
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	if err := viper.ReadInConfig(); err != nil { // viper解析配置文件
		return err
	}

	return nil
}

func (c *Config) initHelmRepository() error {
	configBody, err := ioutil.ReadFile(fmt.Sprintf("./conf/%s", viper.GetString("helmConfig.config")))
	if err != nil {
		klog.Fatalln(err)
		return err
	}
	err = yaml.Unmarshal(configBody, HelmConf)
	if err != nil {
		klog.Fatalln(err)
		return err
	}

	// init repo
	for _, c := range HelmConf.HelmRepos {
		err = initRepository(c)
		if err != nil {
			klog.Fatalln(err)
			return err
		}
	}
	return nil
}

func initRepository(c *repo.Entry) error {
	// Ensure the file directory exists as it is required for file locking
	err := os.MkdirAll(filepath.Dir(Settings.RepositoryConfig), os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return err
	}

	// Acquire a file lock for process synchronization
	fileLock := flock.New(strings.Replace(Settings.RepositoryConfig, filepath.Ext(Settings.RepositoryConfig), ".lock", 1))
	lockCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	locked, err := fileLock.TryLockContext(lockCtx, time.Second)
	if err == nil && locked {
		defer fileLock.Unlock()
	}
	if err != nil {
		return err
	}

	b, err := ioutil.ReadFile(Settings.RepositoryConfig)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	var f repo.File
	if err := yaml.Unmarshal(b, &f); err != nil {
		return err
	}

	r, err := repo.NewChartRepository(c, getter.All(Settings))
	if err != nil {
		return err
	}

	if _, err := r.DownloadIndexFile(); err != nil {
		return err
	}

	f.Update(c)

	if err := f.WriteFile(Settings.RepositoryConfig, 0644); err != nil {
		return err
	}

	return nil
}

func ActionConfigInit(namespace string) (*action.Configuration, error) {
	actionConfig := new(action.Configuration)
	clientConfig := kube.GetConfig(Settings.KubeConfig, Settings.KubeContext, namespace)
	if Settings.KubeToken != "" {
		clientConfig.BearerToken = &Settings.KubeToken
	}
	if Settings.KubeAPIServer != "" {
		clientConfig.APIServer = &Settings.KubeAPIServer
	}
	err := actionConfig.Init(clientConfig, namespace, os.Getenv("HELM_DRIVER"), klog.Infof)
	if err != nil {
		klog.Errorf("%+v", err)
		return nil, err
	}

	return actionConfig, nil
}

// 监控配置文件变化并热加载程序
func (c *Config) watchConfig() {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		klog.Infof("Config file changed: %s", e.Name)
	})
}
