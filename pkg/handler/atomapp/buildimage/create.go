package buildimage

import (
	"fmt"
	"hello-k8s/pkg/errno"
	"hello-k8s/pkg/kubernetes/client"

	. "hello-k8s/pkg/api/v1"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	KanikoImage = "cuijx/executor:20200311"

	VolumeMountPath = "/code"
)

// @Summary 创建用来构建Docker Image 的Job.
// @Description 创建用来构建Docker Image 的Job.
// @Tags atomapp
// @Accept json
// @Produce json
// @param data body buildimage.CreateBuildImageRequest true "创建用来构建Docker Image 的Job时所需的参数."
// @Success 200 {object} handler.Response "{"code":0,"message":"OK","data":{""}}"
// @Router /atomapp/buildimage/create [post]
func Create(c *gin.Context) {
	log.Info("调用创建克隆代码Job对象的函数")

	var r CreateBuildImageRequest
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

	job := newKanikoJob(r)
	result, err := clientset.BatchV1().Jobs(r.Namespace).Create(job)
	if err != nil {
		SendResponse(c, errno.ErrCreateBuildImageJob, err)
		return
	}

	SendResponse(c, errno.OK, result)
}

func newKanikoJob(r CreateBuildImageRequest) *batchv1.Job {
	job := batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: r.Name,
		},
	}

	dockerFileArg := fmt.Sprintf("--dockerfile=%s/%s/Dockerfile", VolumeMountPath, r.ContextPath)
	contextArg := fmt.Sprintf("--context=%s/%s", VolumeMountPath, r.ContextPath)
	var destinationArg string
	if r.ImageTag != "" {
		destinationArg = fmt.Sprintf("--destination=%s:%s", r.Repository, r.ImageTag)
	} else {
		destinationArg = fmt.Sprintf("--destination=%s", r.Repository)
	}

	containers := []corev1.Container{
		{
			Name:  r.Name,
			Image: KanikoImage,
			Args: []string{
				dockerFileArg,
				contextArg,
				destinationArg,
				"--skip-tls-verify=true",
				"--insecure=true",
			},
			VolumeMounts: []corev1.VolumeMount{
				{
					MountPath: VolumeMountPath,
					Name:      r.PersistentVolumeClaimName,
				},
			},
		},
	}

	pvc := corev1.PersistentVolumeClaimVolumeSource{
		ClaimName: r.PersistentVolumeClaimName,
	}

	podSpec := corev1.PodSpec{
		Volumes: []corev1.Volume{
			{
				Name: r.PersistentVolumeClaimName,
				VolumeSource: corev1.VolumeSource{
					PersistentVolumeClaim: &pvc,
				},
			},
		},
		RestartPolicy: corev1.RestartPolicyNever,
		Containers:    containers,
	}

	template := corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Name: r.Name,
		},
		Spec: podSpec,
	}

	job.Spec.Template = template

	return &job
}
