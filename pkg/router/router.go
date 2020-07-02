package router

import (
	_ "hello-k8s/docs"
	mysql5cluster "hello-k8s/pkg/api/v1/clusters/mysql5"
	pgcluster "hello-k8s/pkg/api/v1/clusters/postgres"
	"hello-k8s/pkg/api/v1/helm/charts"
	"hello-k8s/pkg/api/v1/helm/envs"
	"hello-k8s/pkg/api/v1/helm/releases"
	"hello-k8s/pkg/api/v1/helm/repositories"
	"hello-k8s/pkg/api/v1/operators/mysql5"
	"hello-k8s/pkg/api/v1/operators/postgres"
	"hello-k8s/pkg/api/v1/resources/configmap"
	"hello-k8s/pkg/api/v1/resources/container"
	"hello-k8s/pkg/api/v1/resources/cronjob"
	"hello-k8s/pkg/api/v1/resources/deployment"
	"hello-k8s/pkg/api/v1/resources/job"
	"hello-k8s/pkg/api/v1/resources/persistentvolumeclaim"
	"hello-k8s/pkg/api/v1/resources/pod"
	"hello-k8s/pkg/api/v1/resources/secret"
	"hello-k8s/pkg/api/v1/resources/service"
	"hello-k8s/pkg/api/v1/resources/storageclass"
	"hello-k8s/pkg/api/v1/sd"
	"hello-k8s/pkg/api/v1/user"
	"hello-k8s/pkg/router/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// Load loads the middlewares, routes, handlers.
func Load(g *gin.Engine, mw ...gin.HandlerFunc) *gin.Engine {
	// Middlewares.
	// 在处理某些请求时可能因为程序 bug 或者其他异常情况导致程序 panic，
	// 这时候为了不影响下一次请求的调用，需要通过 gin.Recovery()来恢复 API 服务器
	g.Use(gin.Recovery())

	// 强制浏览器不使用缓存
	g.Use(middleware.NoCache)

	// 浏览器跨域 OPTIONS 请求设置
	g.Use(middleware.Options)

	// 一些安全设置
	g.Use(middleware.Secure)
	g.Use(mw...)
	// 404 Handler.
	g.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "The incorrect API route.")
	})

	// swagger api docs
	g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	u := g.Group("/v1/user")
	{
		u.POST("", user.Create)
	}

	a := g.Group("/v1/addon/")
	{
		// mysql5 operator
		a.POST("/mysql5/operator", mysql5.InstallOperator)
		a.DELETE("/mysql5/operator", mysql5.UnInstallOperator)

		// mysql5 cluster
		a.POST("/mysql5/cluster/create", mysql5cluster.CreateCluster)
		a.DELETE("/mysql5/cluster/delete", mysql5cluster.DeleteCluster)

		// postgres-operator
		a.POST("/postgres/operator", postgres.InstallOperator)
		a.DELETE("/postgres/operator", postgres.UnInstallOperator)

		// postgres-cluster
		a.POST("/postgresql/cluster/create", pgcluster.CreateCluster)
	}

	r := g.Group("/v1/resource")
	{
		r.POST("/persistentvolumeclaim/create", persistentvolumeclaim.Create)
		r.DELETE("/persistentvolumeclaim/delete", persistentvolumeclaim.Delete)
		r.GET("/persistentvolumeclaim/detail/:name/:namespace", persistentvolumeclaim.GetPersistentVolumeClaim)
		r.GET("/persistentvolumeclaim/list/:namespace", persistentvolumeclaim.GetPersistentVolumeClaimList)

		r.POST("/job/create", job.Create)
		r.DELETE("/job/delete", job.DeleteJob)
		r.GET("/job/detail/:name/:namespace", job.GetJob)
		r.GET("/job/list/:namespace", job.GetJobList)
		r.GET("/job/pods/:name/:namespace", job.GetJobPods)

		r.POST("/cronjob/create", cronjob.Create)
		r.DELETE("cronjob/delete", cronjob.DeleteCronJob)
		r.GET("/cronjob/detail/:name/:namespace", cronjob.GetCronJob)
		r.GET("/cronjob/list/:namespace", cronjob.GetCronJobList)

		r.DELETE("/deployment/delete", deployment.Delete)
		r.GET("/deployment/detail/:name/:namespace", deployment.GetDeployment)
		r.GET("/deployment/list/:namespace", deployment.GetDeploymentList)
		r.GET("/deployment/pods/:name/:namespace", deployment.GetDeploymentPods)

		r.DELETE("/service/delete", service.Delete)
		r.GET("/service/detail/:name/:namespace", service.GetService)
		r.GET("/service/list/:namespace", service.GetServiceList)
		r.GET("/service/pods/:name/:namespace", service.GetServicePods)

		r.GET("/storageclass/detail/:name", storageclass.GetStorageClass)
		r.GET("/storageclass/list", storageclass.GetStorageClassList)

		r.POST("/secret/create", secret.Create)
		r.DELETE("/secret/delete", secret.Delete)
		r.GET("/secret/detail/:name/:namespace", secret.GetSecret)
		r.GET("/secret/list/:namespace", secret.GetSecretList)

		r.POST("/configmap/create", configmap.Create)
		r.GET("/configmap/detail/:name/:namespace", configmap.GetConfigMap)
		r.GET("/configmap/list/:namespace", configmap.GetConfigMapList)
		r.DELETE("/configmap/delete", configmap.Delete)

		r.GET("/pod/detail/:name/:namespace", pod.GetPod)
		r.GET("/pod/list/:namespace", pod.GetPodList)
		r.GET("/pod/container/:podId/:namespace", container.GetPodContainers)
		r.GET("/container/logs/:namespace/:podId/:container", container.GetLogs)
	}

	h := g.Group("/v1/helm")
	{
		h.GET("/repositories/charts", repositories.ListRepoCharts)
		h.PUT("/repositories", repositories.UpdateRepositories)

		h.GET("/charts", charts.GetChartInfo)

		h.GET("/envs", envs.GetHelmEnvs)

		h.GET("/namespaces/:namespace/releases", releases.ListReleases)
		h.GET("/namespaces/:namespace/releases/:release", releases.ShowReleaseInfo)
	}

	// The health check handlers
	svcd := g.Group("/sd")
	{
		svcd.GET("/health", sd.HealthCheck)
		svcd.GET("/disk", sd.DiskCheck)
		svcd.GET("/cpu", sd.CPUCheck)
		svcd.GET("/ram", sd.RAMCheck)
	}

	return g
}
