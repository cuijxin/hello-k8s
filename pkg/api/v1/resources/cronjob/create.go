package cronjob

import (
	"hello-k8s/pkg/errno"
	"hello-k8s/pkg/kubernetes/client"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	. "hello-k8s/pkg/handler"
	batchv1 "k8s.io/api/batch/v1"
	batch2 "k8s.io/api/batch/v1beta1"
	api "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// @Summary 创建 CronJob 对象
// @Description 创建 CronJob 对象
// @Tags resource
// @Accept json
// @Produce json
// @param data body cronjob.CreateCronJobRequest true "创建 CronJob 对象所需参数."
// @Success 200 {object} handler.Response "{"code":200, "message":"OK", "data":{""}}"
// @Router /resource/cronjob/create [post]
func Create(c *gin.Context) {
	log.Info("调用创建 Job 对象的函数")

	var r CreateCronJobRequest
	if err := c.BindJSON(&r); err != nil {
		SendResponse(c, errno.ErrBind, err)
		return
	}

	// Init kubernetes client
	clientset, err := client.New()
	if err != nil {
		SendResponse(c, errno.ErrCreateK8sClientSet, nil)
		return
	}

	CreateNamespace(r.Namespace, clientset)

	cronjob := newCronJob(r)
	result, err := clientset.BatchV1beta1().CronJobs(r.Namespace).Create(cronjob)
	if err != nil {
		SendResponse(c, errno.ErrCreateCronJob, err)
	}

	SendResponse(c, errno.OK, result)
}

func newCronJob(r CreateCronJobRequest) *batch2.CronJob {
	labels := GetLabelsMap(r.CronJob.CronJobTemplate.PodTemplate.Labels)
	objectMeta := metav1.ObjectMeta{
		Name:   r.Name,
		Labels: labels,
	}

	podSpec := CreatePodSpec(r.Name, r.CronJob.CronJobTemplate.PodTemplate)

	podTemplate := api.PodTemplateSpec{
		ObjectMeta: objectMeta,
		Spec:       *podSpec,
	}

	jobSpec := batchv1.JobSpec{
		Template: podTemplate,
	}

	if reflect.ValueOf(r.CronJob.CronJobTemplate).FieldByName("Parallelism").IsValid() {
		jobSpec.Parallelism = r.CronJob.CronJobTemplate.Parallelism
	}
	if reflect.ValueOf(r.CronJob.CronJobTemplate).FieldByName("Completions").IsValid() {
		jobSpec.Completions = r.CronJob.CronJobTemplate.Completions
	}
	if reflect.ValueOf(r.CronJob.CronJobTemplate).FieldByName("ActiveDeadlineSeconds").IsValid() {
		jobSpec.ActiveDeadlineSeconds = r.CronJob.CronJobTemplate.ActiveDeadlineSeconds
	}

	spec := batch2.CronJobSpec{
		Schedule: r.CronJob.Schedule,
		JobTemplate: batch2.JobTemplateSpec{
			Spec: jobSpec,
		},
	}

	cronJob := batch2.CronJob{
		ObjectMeta: objectMeta,
		Spec:       spec,
	}

	return &cronJob
}
