package mysql5

import (
	"context"
	"fmt"
	"hello-k8s/pkg/kubernetes/client"
	"hello-k8s/pkg/kubernetes/component/addons"
	"hello-k8s/pkg/utils/tool"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	v1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog"
)

type MySQL5Operator struct {
	Namespace  string
	Image      string
	AgentImage string
}

var _ addons.AddOn = &MySQL5Operator{}

func New() *MySQL5Operator {
	return &MySQL5Operator{
		Namespace:  "mysql5-operator",
		Image:      "cuijx/mysql5-operator:v0.18.2.1",
		AgentImage: "cuijx/mysql5-agent",
	}
}

func (o *MySQL5Operator) Deploy(c *client.HelloK8SClient, options addons.AddOnOptions) (err error) {
	klog.Info("开始安装MySQL V5 Operator......")
	tool.CreateNamespace(o.Namespace, c.K8sClientset)
	klog.Info("命名空间检查完成......")

	err = createClusterCRD(c.ApiExtensionsClient)
	defer func() {
		if err != nil {
			destroyClusterCRD(c.ApiExtensionsClient)
		}
	}()
	klog.Info("创建Cluster CustomResourceDefine资源对象完成......")

	err = createBackupCRD(c.ApiExtensionsClient)
	defer func() {
		if err != nil {
			destoryBackupCRD(c.ApiExtensionsClient)
		}
	}()
	klog.Info("创建Backup CustomResourceDefine资源对象完成......")

	err = createRestoreCRD(c.ApiExtensionsClient)
	defer func() {
		if err != nil {
			destoryRestoreCRD(c.ApiExtensionsClient)
		}
	}()
	klog.Info("创建Restore CustomResourceDefine资源对象完成......")

	err = createBackupScheduleCRD(c.ApiExtensionsClient)
	defer func() {
		if err != nil {
			destoryBackupScheduleCRD(c.ApiExtensionsClient)
		}
	}()
	klog.Info("创建BackupSchedule CustomResourceDefine资源对象完成......")

	err = createOperatorSA(o.Namespace, c.K8sClientset)
	defer func() {
		if err != nil {
			destoryOperatorSA(o.Namespace, c.K8sClientset)
		}
	}()
	klog.Info("创建Operator ServiceAccount资源对象完成......")

	err = createAgentSA(o.Namespace, c.K8sClientset)
	defer func() {
		if err != nil {
			destoryAgentSA(o.Namespace, c.K8sClientset)
		}
	}()
	klog.Info("创建Agent ServiceAccount资源对象完成......")

	err = createOperatorClusterRole(c.K8sClientset)
	defer func() {
		if err != nil {
			destoryOperatorClusterRole(c.K8sClientset)
		}
	}()
	klog.Info("创建MySQL V5 Operator 的ClusterRole资源对象完成......")

	err = createAgentClusterRole(c.K8sClientset)
	defer func() {
		if err != nil {
			destoryAgentClusterRole(c.K8sClientset)
		}
	}()
	klog.Info("创建MySQL V5 Operator 的Agent的ClusterRole资源对象完成......")

	err = createOperatorClusterRoleBinding(c.K8sClientset)
	defer func() {
		if err != nil {
			destoryOperatorClusterRoleBinding(c.K8sClientset)
		}
	}()
	klog.Info("创建MySQL V5 Operator 的ClusterRoleBinding资源对象完成......")

	err = createAgentClusterRoleBinding(c.K8sClientset)
	defer func() {
		if err != nil {
			destoryAgentClusterRoleBinding(c.K8sClientset)
		}
	}()
	klog.Info("创建MySQL V5 Operator 的Agent的ClusterRoleBinding资源对象完成......")

	err = createDeployment(o.Namespace, o.Image, o.AgentImage, c.K8sClientset)
	defer func() {
		if err != nil {
			destoryDeployment(o.Namespace, c.K8sClientset)
		}
	}()
	klog.Info("创建MySQL V5 Operator 的Deployment资源对象完成......")

	return nil
}

// Delete 在Kubernetes集群中删除MySQL V5 Operator.
func (o *MySQL5Operator) Delete(c *client.HelloK8SClient) (err error) {
	klog.Info("开始删除mysql5 operator对象......")

	if err = destoryDeployment(o.Namespace, c.K8sClientset); err != nil {
		return
	}
	klog.Info("删除mysql5 operator的deployment对象.")

	if err = destoryAgentClusterRoleBinding(c.K8sClientset); err != nil {
		return
	}
	klog.Info("删除mysql5 operator的agent的clusterrolebinding对象.")

	if err = destoryOperatorClusterRoleBinding(c.K8sClientset); err != nil {
		return
	}
	klog.Info("删除mysql5 operator的clusterrolebinding对象.")

	if err = destoryAgentClusterRole(c.K8sClientset); err != nil {
		return
	}
	klog.Info("删除mysql5 operator的agent的clusterrole对象.")

	if err = destoryOperatorClusterRole(c.K8sClientset); err != nil {
		return
	}
	klog.Info("删除mysql5 operator的clusterrole对象.")

	if err = destoryAgentSA(o.Namespace, c.K8sClientset); err != nil {
		return
	}
	klog.Info("删除mysql5 operator的agent的ServiceAccount对象.")

	if err = destoryOperatorSA(o.Namespace, c.K8sClientset); err != nil {
		return
	}
	klog.Info("删除mysql5 operator的ServiceAccount对象.")

	if err = destoryBackupScheduleCRD(c.ApiExtensionsClient); err != nil {
		return
	}
	klog.Info("删除mysql5 operator对象的BackupSchedule CustomResourceDefine对象.")

	if err = destoryRestoreCRD(c.ApiExtensionsClient); err != nil {
		return
	}
	klog.Info("删除mysql5 operator对象的Restore CustomResourceDefine对象.")

	if err = destoryBackupCRD(c.ApiExtensionsClient); err != nil {
		return
	}
	klog.Info("删除mysql5 operator对象的Backup CustomResourceDefine对象.")

	if err = destroyClusterCRD(c.ApiExtensionsClient); err != nil {
		return
	}
	klog.Info("删除mysql5 operator对象的Cluster CustomResourceDefine对象.")

	klog.Info("成功删除mysql5 operator对象......")

	return nil
}

func createClusterCRD(client *apiextensionsclientset.Clientset) (err error) {
	crd := newClusterCRD()
	if _, err = client.ApiextensionsV1beta1().CustomResourceDefinitions().Create(context.TODO(), crd, metav1.CreateOptions{}); err != nil {
		return
	}
	return nil
}

func newClusterCRD() *v1beta1.CustomResourceDefinition {
	return tool.InitCustomResourceDefinition(
		"mysql5clusters.mysql.oracle.com",
		"MySQLCluster",
		"mysql5cluster",
		"mysql5clusters",
		"mysql.oracle.com",
		"v1")
}

func destroyClusterCRD(client *apiextensionsclientset.Clientset) (err error) {
	if err = client.ApiextensionsV1beta1().CustomResourceDefinitions().Delete(context.TODO(), "mysql5clusters.mysql.oracle.com", metav1.DeleteOptions{}); err != nil {
		return
	}
	return nil
}

func createBackupCRD(client *apiextensionsclientset.Clientset) (err error) {
	crd := newBackupCRD()
	if _, err = client.ApiextensionsV1beta1().CustomResourceDefinitions().Create(context.TODO(), crd, metav1.CreateOptions{}); err != nil {
		return
	}
	return nil
}

func newBackupCRD() *v1beta1.CustomResourceDefinition {
	return tool.InitCustomResourceDefinition(
		"mysql5backups.mysql.oracle.com",
		"MySQLBackup",
		"mysql5backup",
		"mysql5backups",
		"mysql.oracle.com",
		"v1")
}

func destoryBackupCRD(client *apiextensionsclientset.Clientset) (err error) {
	if err = client.ApiextensionsV1beta1().CustomResourceDefinitions().Delete(context.TODO(), "mysql5backups.mysql.oracle.com", metav1.DeleteOptions{}); err != nil {
		return
	}
	return nil
}

func createRestoreCRD(client *apiextensionsclientset.Clientset) (err error) {
	crd := newRestoreCRD()
	if _, err = client.ApiextensionsV1beta1().CustomResourceDefinitions().Create(context.TODO(), crd, metav1.CreateOptions{}); err != nil {
		return
	}
	return nil
}

func newRestoreCRD() *v1beta1.CustomResourceDefinition {
	return tool.InitCustomResourceDefinition(
		"mysql5restores.mysql.oracle.com",
		"MySQLRestore",
		"mysql5restore",
		"mysql5restores",
		"mysql.oracle.com",
		"v1")
}

func destoryRestoreCRD(client *apiextensionsclientset.Clientset) (err error) {
	if err = client.ApiextensionsV1beta1().CustomResourceDefinitions().Delete(context.TODO(), "mysql5restores.mysql.oracle.com", metav1.DeleteOptions{}); err != nil {
		return
	}
	return nil
}

func createBackupScheduleCRD(client *apiextensionsclientset.Clientset) (err error) {
	crd := newBackupScheduleCRD()
	if _, err = client.ApiextensionsV1beta1().CustomResourceDefinitions().Create(context.TODO(), crd, metav1.CreateOptions{}); err != nil {
		return
	}
	return nil
}

func newBackupScheduleCRD() *v1beta1.CustomResourceDefinition {
	return tool.InitCustomResourceDefinition(
		"mysql5backupschedules.mysql.oracle.com",
		"MySQLBackupSchedule",
		"mysql5backupschedule",
		"mysql5backupschedules",
		"mysql.oracle.com",
		"v1")
}

func destoryBackupScheduleCRD(client *apiextensionsclientset.Clientset) (err error) {
	if err = client.ApiextensionsV1beta1().CustomResourceDefinitions().Delete(context.TODO(), "mysql5backupschedules.mysql.oracle.com", metav1.DeleteOptions{}); err != nil {
		return
	}
	return nil
}

func createOperatorSA(namespace string, client kubernetes.Interface) (err error) {
	account := tool.NewServiceAccount("mysql5-operator")
	if _, err = client.CoreV1().ServiceAccounts(namespace).Create(context.TODO(), account, metav1.CreateOptions{}); err != nil {
		return
	}
	return nil
}

func destoryOperatorSA(namespace string, client kubernetes.Interface) (err error) {
	if err = client.CoreV1().ServiceAccounts(namespace).Delete(context.TODO(), "mysql5-operator", metav1.DeleteOptions{}); err != nil {
		return
	}
	return nil
}

func createAgentSA(namespace string, client kubernetes.Interface) (err error) {
	account := tool.NewServiceAccount("mysql5-agent")
	if _, err = client.CoreV1().ServiceAccounts(namespace).Create(context.TODO(), account, metav1.CreateOptions{}); err != nil {
		return
	}
	return nil
}

func destoryAgentSA(namespace string, client kubernetes.Interface) (err error) {
	if err = client.CoreV1().ServiceAccounts(namespace).Delete(context.TODO(), "mysql5-agent", metav1.DeleteOptions{}); err != nil {
		return
	}
	return nil
}

func createOperatorClusterRole(client kubernetes.Interface) (err error) {
	cr := newOperatorClusterRole()
	if _, err = client.RbacV1().ClusterRoles().Create(context.TODO(), cr, metav1.CreateOptions{}); err != nil {
		return
	}
	return nil
}

func newOperatorClusterRole() *rbacv1.ClusterRole {
	return &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: "mysql5-operator",
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
				Resources: []string{"services", "configmaps"},
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
				Resources: []string{"mysql5backups", "mysql5backupschedules", "mysql5clusters", "mysql5clusters/finalizers", "mysql5restores"},
				Verbs:     []string{"get", "list", "patch", "update", "watch"},
			},
			{
				APIGroups: []string{"mysql.oracle.com"},
				Resources: []string{"mysql5backups"},
				Verbs:     []string{"create"},
			},
		},
	}
}

func destoryOperatorClusterRole(client kubernetes.Interface) (err error) {
	if err = client.RbacV1().ClusterRoles().Delete(context.TODO(), "mysql5-operator", metav1.DeleteOptions{}); err != nil {
		return
	}
	return nil
}

func createAgentClusterRole(client kubernetes.Interface) (err error) {
	cr := newAgentClusterRole()
	if _, err = client.RbacV1().ClusterRoles().Create(context.TODO(), cr, metav1.CreateOptions{}); err != nil {
		return
	}
	return nil
}

func newAgentClusterRole() *rbacv1.ClusterRole {
	return &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: "mysql5-agent",
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
				Resources: []string{"mysql5backups", "mysql5backupschedules", "mysql5clusters", "mysql5clusters/finalizers", "mysql5restores"},
				Verbs:     []string{"get", "list", "patch", "update", "watch"},
			},
		},
	}
}

func destoryAgentClusterRole(client kubernetes.Interface) (err error) {
	if err = client.RbacV1().ClusterRoles().Delete(context.TODO(), "mysql5-agent", metav1.DeleteOptions{}); err != nil {
		return
	}
	return nil
}

func createOperatorClusterRoleBinding(client kubernetes.Interface) (err error) {
	crb := newOperatorClusterRoleBinding()
	if _, err = client.RbacV1().ClusterRoleBindings().Create(context.TODO(), crb, metav1.CreateOptions{}); err != nil {
		return
	}
	return nil
}

func newOperatorClusterRoleBinding() *rbacv1.ClusterRoleBinding {
	return tool.InitClusterRoleBinding("ServiceAccount", "ClusterRole", "mysql5-operator",
		"mysql5-operator", "mysql5-operator", "mysql5-operator")
}

func destoryOperatorClusterRoleBinding(client kubernetes.Interface) (err error) {
	if err = client.RbacV1().ClusterRoleBindings().Delete(context.TODO(), "mysql5-operator", metav1.DeleteOptions{}); err != nil {
		return
	}
	return nil
}

func createAgentClusterRoleBinding(client kubernetes.Interface) (err error) {
	crb := newAgentClusterRoleBinding()
	if _, err = client.RbacV1().ClusterRoleBindings().Create(context.TODO(), crb, metav1.CreateOptions{}); err != nil {
		return
	}
	return nil
}

func newAgentClusterRoleBinding() *rbacv1.ClusterRoleBinding {
	return tool.InitClusterRoleBinding("ServiceAccount", "ClusterRole", "mysql5-agent",
		"mysql5-agent", "mysql5-agent", "mysql5-operator")
}

func destoryAgentClusterRoleBinding(client kubernetes.Interface) (err error) {
	if err = client.RbacV1().ClusterRoleBindings().Delete(context.TODO(), "mysql5-agent", metav1.DeleteOptions{}); err != nil {
		return
	}
	return nil
}

func createDeployment(namespace, image, agentImage string, client kubernetes.Interface) (err error) {
	d := newDeployment(namespace, image, agentImage)
	if _, err = client.AppsV1().Deployments(namespace).Create(context.TODO(), d, metav1.CreateOptions{}); err != nil {
		return
	}
	return nil
}

func newDeployment(namespace, image, agentImage string) *appsv1.Deployment {
	d := tool.CreateBasicDeployment(namespace, "mysql5-operator", "app", 1)
	formatAgentImage := fmt.Sprintf("--mysql-agent-image=%s", agentImage)

	containers := []corev1.Container{
		{
			Name:            "mysql5-operator",
			ImagePullPolicy: corev1.PullIfNotPresent,
			Image:           image,
			Ports: []corev1.ContainerPort{
				{
					ContainerPort: 10254,
				},
			},
			Args: []string{
				"--v=4",
				formatAgentImage,
			},
		},
	}

	spec := corev1.PodSpec{
		ServiceAccountName: "mysql5-operator",
		Containers:         containers,
	}

	d.Spec.Template.Spec = spec

	return d
}

func destoryDeployment(namespace string, client kubernetes.Interface) (err error) {
	if err = client.AppsV1().Deployments(namespace).Delete(context.TODO(), "mysql5-operator", metav1.DeleteOptions{}); err != nil {
		return
	}
	return nil
}
