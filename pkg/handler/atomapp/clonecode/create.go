package clonecode

import (
	"fmt"
	"hello-k8s/pkg/errno"
	"hello-k8s/pkg/kubernetes/client"
	"net/url"
	"strings"

	. "hello-k8s/pkg/api/v1"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// CloneCodeJobImage 负责克隆代码的镜像
	CloneCodeJobImage = "cuijx/atom-git:v0.0.2"

	// VolumeMountPath 克隆代码存储卷挂载根路径
	VolumeMountPath = "/code/go-demo"
)

// @Summary 创建克隆Github代码的Job
// @Description 创建构建Docker镜像的Job
// @Tags atomapp
// @Accept json
// @Produce json
// @param data body clonecode.CreateCloneCodeJobRequest true "创建克隆代码的Job时所需参数."
// @Success 200 {object} handler.Response "{"code":0,"message":"OK","data":{""}}"
// @Router /atomapp/clonecode/create [post]
func Create(c *gin.Context) {
	log.Info("调用创建克隆代码Job对象的函数")

	var r CreateCloneCodeJobRequest
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

	cloneCodeJob := newCloneCodeJob(r)
	result, err := clientset.BatchV1().Jobs(r.Namespace).Create(cloneCodeJob)
	if err != nil {
		SendResponse(c, errno.ErrCreateCloneCodeJob, err)
		return
	}

	SendResponse(c, errno.OK, result)
}

func newCloneCodeJob(r CreateCloneCodeJobRequest) *batchv1.Job {
	job := batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: r.Name,
		},
	}

	var repo string
	if r.GithubRepoURL != "" {
		u, _ := url.Parse(r.GithubRepoURL)
		if u.Scheme == "https" {
			tmp := strings.Split(r.GithubRepoURL, "//")
			repo = tmp[len(tmp)-1]
		}
	}

	var cmdStr string
	if r.GithubRepoBranchOrTagName != "" {
		cmdStr = fmt.Sprintf("git clone -b %s https://%s@%s", r.GithubRepoBranchOrTagName, r.GithubAuth.Token, repo)
	} else {
		cmdStr = fmt.Sprintf("git clone https://%s@%s", r.GithubAuth.Token, repo)
	}
	log.Debugf("git cmd is:%s", cmdStr)

	containers := []v1.Container{
		{
			Name:  r.Name,
			Image: CloneCodeJobImage,
			Command: []string{
				"sh",
				"-c",
			},
			Args: []string{
				cmdStr,
			},
			VolumeMounts: []v1.VolumeMount{
				{
					MountPath: VolumeMountPath,
					Name:      r.CodePersistentVolumeClaim,
				},
			},
		},
	}

	pvc := v1.PersistentVolumeClaimVolumeSource{
		ClaimName: r.CodePersistentVolumeClaim,
	}

	template := v1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Name: r.Name,
		},
		Spec: v1.PodSpec{
			Volumes: []v1.Volume{
				{
					Name: r.CodePersistentVolumeClaim,
					VolumeSource: v1.VolumeSource{
						PersistentVolumeClaim: &pvc,
					},
				},
			},
			RestartPolicy: v1.RestartPolicyNever,
			Containers:    containers,
		},
	}

	job.Spec.Template = template

	return &job
}
