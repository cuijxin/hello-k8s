package cronjob

import (
	"hello-k8s/pkg/handler/resources/job"
)

// CronJobArgs 定义了构建一个 CronJob 对象时所需要的参数.
type CronJobArgs struct {
	// The schedule in Cron format, see https://en.wikipedia.org/wiki/Cron.
	Schedule string `json:"schedule" protobuf:"bytes,1,opt,name=schedule"`

	CronJobTemplate job.JobArgs `json:"jobTemplate"`
}

// CreateCronJobRequest 定义了创建一个 CronJob 对象时的请求参数.
type CreateCronJobRequest struct {
	// Name CronJob 对象名称.
	Name string `json:"name"`

	// Namespace 命名空间.
	Namespace string `json:"namespace"`

	//
	CronJob CronJobArgs `json:"cronjob"`
}

// DeleteCronJobRequest 定义了删除CronJob对象时所需参数
type DeleteCronJobRequest struct {
	// CronJob 对象名称.
	Name string `json:"name"`

	// Namespace 命名空间.
	Namespace string `json:"namespace"`
}
