package redis

import (
	"hello-k8s/pkg/errno"
	"hello-k8s/pkg/kubernetes/client"

	. "hello-k8s/pkg/handler"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// RedisOperatorContainerName Redis operator 的容器名称.
	RedisOperatorContainerName = "app"

	// RedisOperatorImage Redis operator 镜像.
	RedisOperatorImage = "quay.io/spotahome/redis-operator:latest"

	// Redis operator 字符串.
	RedisOperator = "redisoperator"

	// DefaultNamespace 命名空间.
	DefaultNamespace = "default"

	// RedisOperatorCRD Redis operator  CRD 名称.
	RedisOperatorCRD = "redisfailovers.databases.spotahome.com"

	// SubjectKind SubjectKind 定义.
	SubjectKind = "ServiceAccount"

	// RoleRef RoleRef 定义.
	RoleRef = "ClusterRole"
)

// @Summary CreateOperator 安装 redis operator.
// @Description CreateOperator 安装 redis operator.
// @Tags operator
// @Accept json
// @Produce json
// @param data body redis.CreateOperatorRequest true "安装Redis operator所需参数."
// @Success 200 {object} handler.Response  "{"code":0,"message":"OK","data":{""}}"
// @Router /operator/redisoperator [post]
func CreateOperator(c *gin.Context) {
	log.Info("调用安装 Redis Operator 组件的函数.")

	clientset, err := client.New()
	if err != nil {
		SendResponse(c, errno.ErrCreateK8sClientSet, nil)
		return
	}

	var r CreateOperatorRequest
	if err := c.Bind(&r); err != nil {
		SendResponse(c, errno.ErrBind, nil)
		return
	}

	CreateNamespace(r.Namespace, clientset)

	sa := NewServiceAccount(RedisOperator)
	if _, err := clientset.CoreV1().ServiceAccounts(r.Namespace).Create(sa); err != nil {
		SendResponse(c, errno.ErrCreateServiceAccount, nil)
		return
	}

	clusterRole := newClusterRole()
	if _, err := clientset.RbacV1().ClusterRoles().Create(clusterRole); err != nil {
		SendResponse(c, errno.ErrCreateClusterRole, nil)
		return
	}

	clusterRoleBinding := newClusterRoleBinding(r.Namespace)
	if _, err := clientset.RbacV1().ClusterRoleBindings().Create(clusterRoleBinding); err != nil {
		SendResponse(c, errno.ErrCreateClusterRoleBinding, nil)
		return
	}

	deployment := newDeployment(r.Namespace, RedisOperator, RedisOperatorImage, RedisOperatorContainerName, RedisOperator)
	if _, err := clientset.AppsV1().Deployments(r.Namespace).Create(deployment); err != nil {
		SendResponse(c, errno.ErrCreateDeployment, err)
		return
	}

	SendResponse(c, errno.OK, nil)
}

// @Summary 从Kubernetes集群中删除redis operator组件
// @Description 从Kubernetes集群中删除redis operator组件
// @Tags operator
// @Accept json
// @Produce json
// @param data body redis.DeleteOperatorRequest true "删除Redis Operator组件时所需的参数"
// @Success 200 {object} handler.Response "{"code":200,"message":"OK","data":{""}}"
// @Router /operator/redisoperator [delete]
func DeleteOperator(c *gin.Context) {
	log.Info("调用删除Redis Operator组件的函数.")

	clientset, err := client.New()
	if err != nil {
		SendResponse(c, errno.ErrCreateK8sClientSet, nil)
		return
	}

	apiClientset, err := client.NewApiExtensionsClient()
	if err != nil {
		SendResponse(c, errno.ErrCreateApiClientSet, nil)
		return
	}

	var r DeleteOperatorRequest
	if err := c.Bind(&r); err != nil {
		SendResponse(c, errno.ErrBind, nil)
		return
	}

	opt := &metav1.DeleteOptions{}
	clientset.AppsV1().Deployments(r.Namespace).Delete(RedisOperator, opt)
	clientset.RbacV1().ClusterRoles().Delete(RedisOperator, opt)
	clientset.RbacV1().ClusterRoleBindings().Delete(RedisOperator, opt)
	clientset.CoreV1().ServiceAccounts(r.Namespace).Delete(RedisOperator, opt)

	apiClientset.ApiextensionsV1beta1().CustomResourceDefinitions().Delete(RedisOperatorCRD, opt)

	SendResponse(c, errno.OK, nil)
}

func newClusterRole() *rbacv1.ClusterRole {
	return &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: RedisOperator,
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{"databases.spotahome.com"},
				Resources: []string{"redisfailovers"},
				Verbs:     []string{"*"},
			},
			{
				APIGroups: []string{"apiextensions.k8s.io"},
				Resources: []string{"customresourcedefinitions"},
				Verbs:     []string{"*"},
			},
			{
				APIGroups: []string{""},
				Resources: []string{"pods", "services", "endpoints", "events", "configmaps"},
				Verbs:     []string{"*"},
			},
			{
				APIGroups: []string{""},
				Resources: []string{"secrets"},
				Verbs:     []string{"get"},
			},
			{
				APIGroups: []string{"apps"},
				Resources: []string{"deployments", "statefulsets"},
				Verbs:     []string{"*"},
			},
			{
				APIGroups: []string{"policy"},
				Resources: []string{"poddisruptionbudgets"},
				Verbs:     []string{"*"},
			},
		},
	}
}

func newClusterRoleBinding(namespace string) *rbacv1.ClusterRoleBinding {
	return &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: RedisOperator,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      SubjectKind,
				Name:      RedisOperator,
				Namespace: namespace,
			},
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: rbacv1.GroupName,
			Kind:     RoleRef,
			Name:     RedisOperator,
		},
	}
}

func newDeployment(namespace, redisoperatorName, image, containerName, serviceAccountName string) *appsv1.Deployment {
	deployment := CreateBasicDeployment(namespace, redisoperatorName, "app", 1)

	rqCPULimit, _ := resource.ParseQuantity("100")
	rqMemoryLimit, _ := resource.ParseQuantity("50Mi")
	rqCPURequest, _ := resource.ParseQuantity("10m")
	rqMemoryRequest, _ := resource.ParseQuantity("50Mi")

	containers := []corev1.Container{
		{
			Image:           image,
			ImagePullPolicy: corev1.PullIfNotPresent,
			Name:            containerName,
			Resources: corev1.ResourceRequirements{
				Limits:   corev1.ResourceList{corev1.ResourceCPU: rqCPULimit, corev1.ResourceMemory: rqMemoryLimit},
				Requests: corev1.ResourceList{corev1.ResourceCPU: rqCPURequest, corev1.ResourceMemory: rqMemoryRequest},
			},
		},
	}

	spec := corev1.PodSpec{
		ServiceAccountName: serviceAccountName,
		Containers:         containers,
		RestartPolicy:      corev1.RestartPolicyAlways,
	}

	deployment.Spec.Template.Spec = spec

	return deployment
}
