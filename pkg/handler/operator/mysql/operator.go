package mysql

import (
	. "hello-k8s/pkg/api/v1"
	"hello-k8s/pkg/kubernetes/client"
	"hello-k8s/pkg/utils/errno"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
)

const (
	CRDGROUP         = "mysql.oracle.com"
	CRDVERSION       = "v1alpha1"
	DEFAULTNAMESPACE = "mysql-operator"
	MYSQLOPERATOR    = "mysql-operator"
	MYSQLAGENT       = "mysql-agent"

	MYSQLOPERATORCONTAINERNAME = "mysql-operator-controller"
	MYSQLOPERATORIMAGE         = "n1ce37/mysql-operator:0.3.0"
	MYSQLAGENTIMAGE            = "n1ce37/mysql-agent"

	MYSQLCLUSTERNAME        = "mysqlclusters.mysql.oracle.com"
	MYSQLBACKUPNAME         = "mysqlbackups.mysql.oracle.com"
	MYSQLRESTORENAME        = "mysqlrestores.mysql.oracle.com"
	MYSQLBACKUPSCHEDULENAME = "mysqlbackupschedules.mysql.oracle.com"

	SUBJECTKIND = "ServiceAccount"
	ROLEREF     = "ClusterRole"

	MYSQLAGENTIMAGEARG = "--mysql-agent-image=n1ce37/mysql-agent"
	MYSQLVERSIONARG    = "--v=4"
)

// @Summary 安装 mysql operator.
// @Description 安装 mysql operator.
// @Tags operator
// @Accept json
// @Produce json
// @param data body mysql.CreateOperatorRequest true "安装MySQL Operator组件时所需参数."
// @Success 200 {object} handler.Response  "{"code":0,"message":"OK","data":{""}}"
// @Router /operator/mysqloperator [post]
func CreateOperator(c *gin.Context) {
	log.Debug("调用安装 MySQL Operator 组件的函数.")

	clientset, err := client.New()
	if err != nil {
		SendResponse(c, errno.ErrCreateK8sClientSet, nil)
		return
	}

	apiExtensionsClientset, err := client.NewApiExtensionsClient()
	if err != nil {
		SendResponse(c, errno.ErrCreateApiClientSet, nil)
	}

	var r CreateOperatorRequest
	if err := c.Bind(&r); err != nil {
		SendResponse(c, errno.ErrBind, nil)
		return
	}

	if err := Deploy(r.Namespace, clientset, apiExtensionsClientset); err != nil {
		log.Error("create mysql operator err: %v", err)
		SendResponse(c, errno.ErrCreateMySQLOperator, err)
	}

	SendResponse(c, errno.OK, nil)
}

// @Summary 从Kubernetes集群中删除mysql operator组件
// @Description 从Kubernetes集群中删除mysql operator组件
// @Tags operator
// @Accept json
// @Produce json
// @param data body mysql.DeleteOperatorRequest true "删除MySQL Operator组件时所需的参数"
// @Success 200 {object} handler.Response "{"code":200,"message":"OK","data":{""}}"
// @Router /operator/mysqloperator [delete]
func DeleteOperator(c *gin.Context) {
	log.Debug("调用删除MySQL Operator组件的函数.")

	clientset, err := client.New()
	if err != nil {
		SendResponse(c, errno.ErrCreateK8sClientSet, nil)
		return
	}

	apiExtensionsClientset, err := client.NewApiExtensionsClient()
	if err != nil {
		SendResponse(c, errno.ErrCreateApiClientSet, nil)
	}

	var r DeleteOperatorRequest
	if err := c.Bind(&r); err != nil {
		SendResponse(c, errno.ErrBind, nil)
		return
	}

	opt := &metav1.DeleteOptions{}
	clientset.AppsV1().Deployments(r.Namespace).Delete(MYSQLOPERATOR, opt)
	clientset.RbacV1().ClusterRoleBindings().Delete(MYSQLAGENT, opt)
	clientset.RbacV1().ClusterRoleBindings().Delete(MYSQLOPERATOR, opt)
	clientset.RbacV1().ClusterRoles().Delete(MYSQLAGENT, opt)
	clientset.RbacV1().ClusterRoles().Delete(MYSQLOPERATOR, opt)
	clientset.CoreV1().ServiceAccounts(r.Namespace).Delete(MYSQLAGENT, opt)
	clientset.CoreV1().ServiceAccounts(r.Namespace).Delete(MYSQLOPERATOR, opt)

	apiExtensionsClientset.ApiextensionsV1beta1().CustomResourceDefinitions().Delete(MYSQLBACKUPSCHEDULENAME, opt)
	apiExtensionsClientset.ApiextensionsV1beta1().CustomResourceDefinitions().Delete(MYSQLRESTORENAME, opt)
	apiExtensionsClientset.ApiextensionsV1beta1().CustomResourceDefinitions().Delete(MYSQLBACKUPNAME, opt)
	apiExtensionsClientset.ApiextensionsV1beta1().CustomResourceDefinitions().Delete(MYSQLCLUSTERNAME, opt)

	SendResponse(c, errno.OK, nil)
}

// 部署 mysql operator 组件.
func Deploy(namespace string, client clientset.Interface, apiExtensionsClient *apiextensionsclientset.Clientset) error {
	CreateNamespace(namespace, client)

	mysqlClusterCRD := newMySQLClusterCRD()
	if _, err := apiExtensionsClient.ApiextensionsV1beta1().CustomResourceDefinitions().Create(mysqlClusterCRD); err != nil {
		return err
	}
	mysqlBackupCRD := newMySQLBackupCRD()
	if _, err := apiExtensionsClient.ApiextensionsV1beta1().CustomResourceDefinitions().Create(mysqlBackupCRD); err != nil {
		return err
	}
	mysqlRestoreCRD := newMySQLRestoreCRD()
	if _, err := apiExtensionsClient.ApiextensionsV1beta1().CustomResourceDefinitions().Create(mysqlRestoreCRD); err != nil {
		return err
	}
	mysqlBackupScheduleCRD := newMySQLBackupScheduleCRD()
	if _, err := apiExtensionsClient.ApiextensionsV1beta1().CustomResourceDefinitions().Create(mysqlBackupScheduleCRD); err != nil {
		return err
	}
	mysqlOperatorServiceAccount := NewServiceAccount(MYSQLOPERATOR)
	if _, err := client.CoreV1().ServiceAccounts(namespace).Create(mysqlOperatorServiceAccount); err != nil {
		return err
	}
	mysqlAgentServiceAccount := NewServiceAccount(MYSQLAGENT)
	if _, err := client.CoreV1().ServiceAccounts(namespace).Create(mysqlAgentServiceAccount); err != nil {
		return err
	}
	mysqlOperatorClusterRole := newMySQLOperatorClusterRole()
	if _, err := client.RbacV1().ClusterRoles().Create(mysqlOperatorClusterRole); err != nil {
		return err
	}
	mysqlAgentClusterRole := newMySQLAgentClusterRole()
	if _, err := client.RbacV1().ClusterRoles().Create(mysqlAgentClusterRole); err != nil {
		return err
	}
	mysqlOperatorClusterRoleBinding := newMySQLOperatorClusterRoleBinding(namespace)
	if _, err := client.RbacV1().ClusterRoleBindings().Create(mysqlOperatorClusterRoleBinding); err != nil {
		return err
	}
	mysqlAgentClusterRoleBinding := newMySQLAgentClusterRoleBinding(namespace)
	if _, err := client.RbacV1().ClusterRoleBindings().Create(mysqlAgentClusterRoleBinding); err != nil {
		return err
	}

	mysqlOperatorDeployment := newDeployment(namespace)
	if _, err := client.AppsV1().Deployments(namespace).Create(mysqlOperatorDeployment); err != nil {
		return err
	}

	return nil
}

func newMySQLClusterCRD() *v1beta1.CustomResourceDefinition {
	return InitCustomResourceDefinition(
		MYSQLCLUSTERNAME,
		"Cluster",
		"mysqlcluster",
		"mysqlclusters",
		CRDGROUP,
		CRDVERSION)
}

func newMySQLBackupCRD() *v1beta1.CustomResourceDefinition {
	return InitCustomResourceDefinition(
		MYSQLBACKUPNAME,
		"Backup",
		"mysqlbackup",
		"mysqlbackups",
		CRDGROUP,
		CRDVERSION)
}

func newMySQLRestoreCRD() *v1beta1.CustomResourceDefinition {
	return InitCustomResourceDefinition(
		MYSQLRESTORENAME,
		"Restore",
		"mysqlrestore",
		"mysqlrestores",
		CRDGROUP,
		CRDVERSION)
}

func newMySQLBackupScheduleCRD() *v1beta1.CustomResourceDefinition {
	return InitCustomResourceDefinition(
		MYSQLBACKUPSCHEDULENAME,
		"BackupSchedule",
		"mysqlbackupschedule",
		"mysqlbackupschedules",
		CRDGROUP,
		CRDVERSION)
}

func newMySQLOperatorClusterRole() *rbacv1.ClusterRole {
	return &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: MYSQLOPERATOR,
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"pods"},
				Verbs:     []string{"get", "list", "patch", "update", "watch"},
			},
			{
				APIGroups: []string{""},
				Resources: []string{"secrets"},
				Verbs:     []string{"get", "create"},
			},
			{
				APIGroups: []string{""},
				Resources: []string{"services"},
				Verbs:     []string{"create", "get", "list", "watch"},
			},
			{
				APIGroups: []string{""},
				Resources: []string{"events"},
				Verbs:     []string{"create", "update", "patch"},
			},
			{
				APIGroups: []string{"apps"},
				Resources: []string{"statefulsets"},
				Verbs:     []string{"create", "get", "list", "patch", "update", "watch"},
			},
			{
				APIGroups: []string{"mysql.oracle.com"},
				Resources: []string{"mysqlbackups", "mysqlbackupschedules", "mysqlclusters", "mysqlclusters/finalizers", "mysqlrestores"},
				Verbs:     []string{"get", "list", "patch", "update", "watch"},
			},
			{
				APIGroups: []string{"mysql.oracle.com"},
				Resources: []string{"mysqlbackups"},
				Verbs:     []string{"create"},
			},
		},
	}
}

func newMySQLAgentClusterRole() *rbacv1.ClusterRole {
	return &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: MYSQLAGENT,
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"pods"},
				Verbs:     []string{"get", "list", "patch", "update", "watch"},
			},
			{
				APIGroups: []string{""},
				Resources: []string{"secrets"},
				Verbs:     []string{"get"},
			},
			{
				APIGroups: []string{""},
				Resources: []string{"events"},
				Verbs:     []string{"create", "update", "patch"},
			},
			{
				APIGroups: []string{"mysql.oracle.com"},
				Resources: []string{"mysqlbackups", "mysqlbackupschedules", "mysqlclusters", "msyqlclusters/finalizers", "mysqlrestores"},
				Verbs:     []string{"get", "list", "patch", "update", "watch"},
			},
		},
	}
}

func newMySQLOperatorClusterRoleBinding(namespace string) *rbacv1.ClusterRoleBinding {
	return InitClusterRoleBinding(SUBJECTKIND, ROLEREF, MYSQLOPERATOR, MYSQLOPERATOR, MYSQLOPERATOR, namespace)
}

func newMySQLAgentClusterRoleBinding(namespace string) *rbacv1.ClusterRoleBinding {
	return InitClusterRoleBinding(SUBJECTKIND, ROLEREF, MYSQLAGENT, MYSQLAGENT, MYSQLAGENT, namespace)
}

func newDeployment(namespace string) *appsv1.Deployment {
	deployment := CreateBasicDeployment(namespace, MYSQLOPERATOR, "app", 1)

	containers := []corev1.Container{
		{
			Name:            MYSQLOPERATORCONTAINERNAME,
			ImagePullPolicy: corev1.PullIfNotPresent,
			Image:           MYSQLOPERATORIMAGE,
			Ports: []corev1.ContainerPort{
				{
					ContainerPort: 10254,
				},
			},
			Args: []string{
				MYSQLVERSIONARG,
				MYSQLAGENTIMAGEARG,
			},
		},
	}

	spec := corev1.PodSpec{
		ServiceAccountName: MYSQLOPERATOR,
		Containers:         containers,
	}

	deployment.Spec.Template.Spec = spec

	return deployment
}
