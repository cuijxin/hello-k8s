package main

import (
	"errors"
	"hello-k8s/pkg/config"
	"hello-k8s/pkg/kubernetes/client"
	"hello-k8s/pkg/router"
	"hello-k8s/pkg/storage/database"
	"hello-k8s/pkg/storage/mongo"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"k8s.io/klog"
)

var (
	cfg = pflag.StringP("config", "c", "", "hello-k8s apiserver config file path.")
)

func main() {
	pflag.Parse()

	// init config
	if err := config.Init(*cfg); err != nil {
		panic(err)
	}

	database.DB = &mongo.MongoDB{}
	if err := database.DB.Init(database.DBInitOptions{
		User:     viper.GetString("db.mongodb.user"),
		Password: viper.GetString("db.mongodb.password"),
		Address:  viper.GetString("db.mongodb.addr"),
	}); err != nil {
		klog.Errorf("failed to connect db server: %v", err)
		panic(err)
	}

	// init kubernetes client
	client.MyClient.InitHelloK8SClient()

	// Set gin mode.
	gin.SetMode(viper.GetString("runmode"))

	// Create the Gin engine.
	g := gin.New()

	// gin middlewares
	middlewares := []gin.HandlerFunc{}

	// Routes.
	router.Load(
		// Cores.
		g,

		// Middlewares.
		middlewares...,
	)

	// Ping the server to make sure the router is working.
	go func() {
		if err := pingServer(); err != nil {
			klog.Fatalf("The router has no response, or it might took too long to start up: %v", err)
		}

		klog.Info("The router has been deployed successfully.")
	}()

	klog.Infof("Start to listening the incoming requests on http address: %s", viper.GetString("addr"))
	klog.Info(http.ListenAndServe(viper.GetString("addr"), g).Error())
}

// pingServer pings the http server to make sure the router is working.
func pingServer() error {
	for i := 0; i < viper.GetInt("max_ping_count"); i++ {
		// Ping the server by sending a GET request to `/health`.
		resp, err := http.Get(viper.GetString("url") + "/sd/health")
		if err == nil && resp.StatusCode == 200 {
			return nil
		}

		// Sleep for a second to continue the next ping.
		klog.Info("Waiting for the router, retry in 1 second.")
		time.Sleep(time.Second)
	}
	return errors.New("Cannot connect to the router.")
}
