package job

import (
	"hello-k8s/pkg/kubernetes/client"
	"hello-k8s/pkg/utils/errno"
	"reflect"

	"hello-k8s/pkg/utils/tool"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	batchv1 "k8s.io/api/batch/v1"
	api "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// @Summary 创建Job对象
// @Description 创建Job对象
// @Tags resource
// @Accept json
// @Produce json
// @param data body job.CreateJobRequest true "创建Job对象所需参数."
// @Success 200 {object} tool.Response "{"code":200, "message":"OK", "data":{""}}"
// @Router /resource/job/create [post]
func Create(c *gin.Context) {
	log.Info("调用创建 Job 对象的函数")

	var r CreateJobRequest
	if err := c.BindJSON(&r); err != nil {
		tool.SendResponse(c, errno.ErrBind, err)
		return
	}

	// Init kubernetes client
	clientset, err := client.New()
	if err != nil {
		tool.SendResponse(c, errno.ErrCreateK8sClientSet, nil)
		return
	}

	tool.CreateNamespace(r.Namespace, clientset)

	job := newJob(r)
	result, err := clientset.BatchV1().Jobs(r.Namespace).Create(job)
	if err != nil {
		tool.SendResponse(c, errno.ErrCreateJob, err)
		return
	}

	tool.SendResponse(c, errno.OK, result)
}

func newJob(r CreateJobRequest) *batchv1.Job {

	labels := tool.GetLabelsMap(r.Job.PodTemplate.Labels)
	objectMeta := metaV1.ObjectMeta{
		Name:   r.Name,
		Labels: labels,
	}

	podSpec := tool.CreatePodSpec(r.Name, r.Job.PodTemplate)

	podTemplate := api.PodTemplateSpec{
		ObjectMeta: objectMeta,
		Spec:       *podSpec,
	}

	jobSpec := batchv1.JobSpec{
		Template: podTemplate,
	}

	job := batchv1.Job{
		ObjectMeta: objectMeta,
		Spec:       jobSpec,
	}

	if reflect.ValueOf(r.Job).FieldByName("Parallelism").IsValid() {
		job.Spec.Parallelism = r.Job.Parallelism
	}
	if reflect.ValueOf(r.Job).FieldByName("Completions").IsValid() {
		job.Spec.Completions = r.Job.Completions
	}
	if reflect.ValueOf(r.Job).FieldByName("ActiveDeadlineSeconds").IsValid() {
		job.Spec.ActiveDeadlineSeconds = r.Job.ActiveDeadlineSeconds
	}

	return &job
}
