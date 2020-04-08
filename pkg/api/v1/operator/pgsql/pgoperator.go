package pgsql

import (
	"hello-k8s/pkg/handler/operator"
	"hello-k8s/pkg/kubernetes/client"
	"hello-k8s/pkg/utils/errno"

	. "hello-k8s/pkg/api/v1"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// DefaultConfigMapName is the name of the postgres-opeartor configmap.
	DefaultConfigMapName = "postgres-operator"

	// DefaultDeploymentName is the name of the postgres-operator deployment.
	DefaultDeploymentName = "postgres-operator"

	// DefaultNamespace is the name of the namespace which the postgres-operator deployed.
	DefaultNamespace = "default"

	// DefaultServiceAccountName is the name of the postgres-opeartor's serviceaccount.
	DefaultServiceAccountName = "zalando-postgres-operator"

	// DefaultClusterRoleName is the name of the postgres-operator's clusterrole.
	DefaultClusterRoleName = "zalando-postgres-operator"

	// DefaultClusterRoleBindingName is the name of the postgres-operator's clusterrolebinding.
	DefaultClusterRoleBindingName = "zalando-postgres-operator"

	// DefaultPostgresOperatorImage is the image of the postgres-operator.
	DefaultPostgresOperatorImage = "cuijx/postgres-operator:58a7746"
)

// @Summary CreateOperator deploy pgsqloperator to the kubernetes cluster
// @Description CreateOperator deploy pgsqloperator to the kubernetes cluster
// @Tags operator
// @Accept json
// @Produce json
// @param data body pgsql.CreateOpeartorRequest true "Deploy pgsqloperator params"
// @Success 200 {object} pgsql.CreateOperatorResponse "{"code":0,"message":"OK","data":{"Namespace":"defalt", "DeploymentName": "postgres-operator"}}"
// @Router /operator/pgsqloperator [post]
func CreateOperator(c *gin.Context) {
	log.Info("Pgsql Operator deploy function called.")

	clientset, err := client.New()
	if err != nil {
		SendResponse(c, errno.ErrCreateK8sClientSet, nil)
		return
	}

	var r CreateOpeartorRequest
	if err := c.Bind(&r); err != nil {
		SendResponse(c, errno.ErrBind, nil)
		return
	}

	cm := newConfigMap(r.ConfigMapName, r.Namespace)
	if _, err := clientset.CoreV1().ConfigMaps(r.Namespace).Create(cm); err != nil {
		SendResponse(c, errno.ErrCreateConfigMap, nil)
		return
	}

	sa := newServiceAccount(r.ServiceAccountName)
	if _, err := clientset.CoreV1().ServiceAccounts(r.Namespace).Create(sa); err != nil {
		SendResponse(c, errno.ErrCreateServiceAccount, nil)
		return
	}

	cr := newClusterRole(r.ClusterRoleName)
	if _, err := clientset.RbacV1().ClusterRoles().Create(cr); err != nil {
		SendResponse(c, errno.ErrCreateClusterRole, nil)
		return
	}

	crb := newClusterRoleBinding(r.ClusterRoleBindingName, r.ClusterRoleName, r.ServiceAccountName, r.Namespace)
	if _, err := clientset.RbacV1().ClusterRoleBindings().Create(crb); err != nil {
		SendResponse(c, errno.ErrCreateClusterRoleBinding, nil)
		return
	}

	deploy := newDeployment(r.Namespace, r.DeploymentName, r.PostgresqlOperatorImage, r.ServiceAccountName)
	if _, err := clientset.AppsV1().Deployments(r.Namespace).Create(deploy); err != nil {
		log.Debugf("error: %v", err)
		SendResponse(c, errno.ErrCreateDeployment, err)
		return
	}

	rsp := CreateOperatorResponse{
		ConfigMapName:          r.ConfigMapName,
		ServiceAccountName:     r.ServiceAccountName,
		ClusterRoleName:        r.ClusterRoleName,
		ClusterRoleBindingName: r.ClusterRoleBindingName,
		DeploymentName:         r.DeploymentName,
		Namespace:              r.Namespace,
	}

	SendResponse(c, nil, rsp)
}

// @Summary DeleteOperator delete pgsqloperator from the kubernetes cluster
// @Description DeleteOperator delete pgsqloperator from the kubernetes cluster
// @Tags operator
// @Accept json
// @Produce json
// @Success 200 {object} handler.Response "{"code":200,"message":"OK","data":{""}}"
// @Router /operator/pgsqloperator [delete]
func DeleteOperator(c *gin.Context) {
	log.Info("Pgsql Operator delete function called.")

	clientset, err := client.New()
	if err != nil {
		SendResponse(c, errno.ErrCreateK8sClientSet, nil)
		return
	}

	opt := &metav1.DeleteOptions{}
	clientset.CoreV1().ConfigMaps(DefaultNamespace).Delete(DefaultConfigMapName, opt)
	clientset.CoreV1().ServiceAccounts(DefaultNamespace).Delete(DefaultServiceAccountName, opt)
	clientset.RbacV1().ClusterRoleBindings().Delete(DefaultClusterRoleBindingName, opt)
	clientset.RbacV1().ClusterRoles().Delete(DefaultClusterRoleName, opt)
	clientset.AppsV1().Deployments(DefaultNamespace).Delete(DefaultDeploymentName, opt)

	SendResponse(c, errno.OK, nil)
}

func newConfigMap(configMapName, namespace string) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMapName,
			Namespace: namespace,
		},
		Data: map[string]string{
			"api_port":                     "8080",
			"aws_region":                   "eu-central-1",
			"cluster_domain":               "cluster.local",
			"cluster_history_entries":      "1000",
			"cluster_labels":               "application:spilo",
			"cluster_name_label":           "version",
			"db_hosted_zone":               "db.example.com",
			"debug_logging":                "true",
			"docker_image":                 "registry.opensource.zalan.do/acid/spilo-cdp-12:1.6-p16",
			"enable_crd_validation":        "false",
			"enable_master_load_balancer":  "false",
			"enable_replica_load_balancer": "false",
			"enable_teams_api":             "false",
			"master_dns_name_format":       "{cluster}.{team}.{hostedzone}",
			"pdb_name_format":              "postgres-{cluster}-pdb",
			"pod_deletion_wait_timeout":    "10m",
			"pod_label_wait_timeout":       "10m",
			"pod_management_policy":        "ordered_ready",
			"pod_role_label":               "spilo-role",
			"pod_service_account_name":     "zalando-postgres-operator",
			"pod_terminate_grace_period":   "5m",
			"ready_wait_interval":          "3s",
			"ready_wait_timeout":           "30s",
			"repair_period":                "5m",
			"replica_dns_name_format":      "{cluster}-repl.{team}.{hostedzone}",
			"replication_username":         "standby",
			"resource_check_interval":      "3s",
			"resource_check_timeout":       "10m",
			"resync_period":                "30m",
			"ring_log_lines":               "100",
			"secret_name_template":         "{username}.{cluster}.credentials",
			"spilo_privileged":             "false",
			"super_username":               "postgres",
			"watched_namespace":            "*",
			"workers":                      "4",
		},
	}
}

func newServiceAccount(serviceAccountName string) *corev1.ServiceAccount {
	return &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name: serviceAccountName,
		},
	}
}

func newClusterRole(clusterRoleName string) *rbacv1.ClusterRole {
	return &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: clusterRoleName,
		},
		Rules: []rbacv1.PolicyRule{
			{ // all verbs allowed for custom operator resources
				APIGroups: []string{"acid.zalan.do"},
				Resources: []string{"postgresqls", "postgresqls/status", "operatorconfigurations"},
				Verbs:     []string{"*"},
			},
			{
				// to create or get/update CRDs when starting up
				APIGroups: []string{"apiextensions.k8s.io"},
				Resources: []string{"customresourcedefinitions"},
				Verbs:     []string{"create", "get", "patch", "update"},
			},
			{
				// to read configuration from ConfigMaps
				APIGroups: []string{""},
				Resources: []string{"configmaps"},
				Verbs:     []string{"get"},
			},
			{
				// to manage endpoints which are also used by Patroni
				APIGroups: []string{""},
				Resources: []string{"endpoints"},
				// watch needed if zalando-postgres-operator account is used for pods as well
				Verbs: []string{"create", "delete", "deletecollection", "get", "list", "patch", "watch"},
			},
			{
				// to CRUD secrets for database access
				APIGroups: []string{""},
				Resources: []string{"secrets"},
				Verbs:     []string{"create", "update", "delete", "get"},
			},
			{
				// to check nodes for node readiness label
				APIGroups: []string{""},
				Resources: []string{"nodes"},
				Verbs:     []string{"get", "list", "watch"},
			},
			{
				// to read or delete existing PVCs. Creation via StatefulSet
				APIGroups: []string{""},
				Resources: []string{"persistentvolumeclaims"},
				Verbs:     []string{"delete", "get", "list"},
			},
			{
				// to read existing PVs. Creation should be done via dynamic provisioning
				APIGroups: []string{""},
				Resources: []string{"persistentvolumes"},
				// update only for resizing AWS volumes
				Verbs: []string{"get", "list", "update"},
			},
			{
				// to watch Spilo pods and do rolling updates. Creation via StatefulSet
				APIGroups: []string{""},
				Resources: []string{"pods"},
				Verbs:     []string{"delete", "get", "list", "watch", "patch"},
			},
			{
				// to resize the filesystem in Spilo pods when increasing volume size
				APIGroups: []string{""},
				Resources: []string{"pods/exec"},
				Verbs:     []string{"create"},
			},
			{
				// to CRUD services to point to Postgres cluster instances
				APIGroups: []string{""},
				Resources: []string{"services"},
				Verbs:     []string{"create", "delete", "get", "patch"},
			},
			{
				// to CRUD the StatefulSet which controls the Postgres cluster instances
				APIGroups: []string{"apps"},
				Resources: []string{"statefulsets"},
				Verbs:     []string{"create", "delete", "get", "list", "patch"},
			},
			{
				// to get namespaces operator resources can run in
				APIGroups: []string{""},
				Resources: []string{"namespaces"},
				Verbs:     []string{"get"},
			},
			{
				// to define PDBs. Update happens via delete/create
				APIGroups: []string{"policy"},
				Resources: []string{"poddisruptionbudgets"},
				Verbs:     []string{"create", "delete", "get"},
			},
			{
				// to create ServiceAccounts in each namespace the operator watches
				APIGroups: []string{""},
				Resources: []string{"serviceaccounts"},
				Verbs:     []string{"get", "create"},
			},
			{
				// to create role bindings to the operator service account
				APIGroups: []string{"rbac.authorization.k8s.io"},
				Resources: []string{"rolebindings"},
				Verbs:     []string{"get", "create"},
			},
			{
				// to CRUD cron jobs for logical backups
				APIGroups: []string{"batch"},
				Resources: []string{"cronjobs"},
				Verbs:     []string{"create", "delete", "get", "list", "patch", "update"},
			},
		},
	}
}

func newClusterRoleBinding(clusterRoleBingName, clusterRoleName, serviceAccountName, namespace string) *rbacv1.ClusterRoleBinding {
	return &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: clusterRoleBingName,
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: rbacv1.GroupName,
			Kind:     "ClusterRole",
			Name:     clusterRoleName,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      serviceAccountName,
				Namespace: namespace,
			},
		},
	}
}

func newDeployment(namespace, deploymentName, image, serviceAccountName string) *appsv1.Deployment {
	// log.Debug("init deployment basic object.")
	deployment := operator.CreateBasicDeployment(namespace, deploymentName, "name", 1)
	// log.Debug("init deployment basic object done.")
	runAsUser := int64(1000)
	runAsNonRoot := true
	readOnlyRootFilesystem := true

	rqCPULimit, _ := resource.ParseQuantity("500m")
	rqMemoryLimit, _ := resource.ParseQuantity("500Mi")
	rqCPURequest, _ := resource.ParseQuantity("100m")
	rqMemoryRequest, _ := resource.ParseQuantity("250Mi")

	containers := []corev1.Container{
		{
			Name:            deploymentName,
			ImagePullPolicy: corev1.PullIfNotPresent,
			Image:           image,
			Resources: corev1.ResourceRequirements{
				Limits:   corev1.ResourceList{corev1.ResourceCPU: rqCPULimit, corev1.ResourceMemory: rqMemoryLimit},
				Requests: corev1.ResourceList{corev1.ResourceCPU: rqCPURequest, corev1.ResourceMemory: rqMemoryRequest},
			},
			SecurityContext: &corev1.SecurityContext{
				RunAsUser:              &runAsUser,
				RunAsNonRoot:           &runAsNonRoot,
				ReadOnlyRootFilesystem: &readOnlyRootFilesystem,
			},
			Env: []corev1.EnvVar{
				{
					Name:  "CONFIG_MAP_NAME",
					Value: deploymentName,
				},
			},
		},
	}

	spec := corev1.PodSpec{
		ServiceAccountName: serviceAccountName,
		Containers:         containers,
	}

	deployment.Spec.Template.Spec = spec

	// log.Debugf("deployment object is %v", deployment)

	return deployment
}
