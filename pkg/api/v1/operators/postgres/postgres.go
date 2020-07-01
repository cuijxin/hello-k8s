package postgres

import (
	"context"
	"hello-k8s/pkg/kubernetes/client"
	"hello-k8s/pkg/kubernetes/component/addons"
	"hello-k8s/pkg/utils/tool"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/klog"
)

type PostgresOperator struct {
	Namespace              string
	Image                  string
	ConfigMapName          string
	ServiceAccountName     string
	ClusterRoleName        string
	ClusterRoleBindingName string
	PodClusterRoleName     string
	DeploymentName         string
	ServiceName            string
}

var _ addons.AddOn = &PostgresOperator{}

func New() *PostgresOperator {
	return &PostgresOperator{
		Namespace:              "postgres-operator",
		Image:                  "cuijx/postgres-operator:94a1a62-dirty",
		ConfigMapName:          "postgres-operator",
		ServiceAccountName:     "postgres-operator",
		ClusterRoleName:        "postgres-operator",
		ClusterRoleBindingName: "postgres-operator",
		PodClusterRoleName:     "postgres-pod",
		DeploymentName:         "postgres-operator",
		ServiceName:            "postgres-operator",
	}
}

func (o *PostgresOperator) Deploy(c *client.HelloK8SClient, options addons.AddOnOptions) (err error) {
	klog.Info("开始安装Postgres Operator......")
	tool.CreateNamespace(o.Namespace, c.K8sClientset)
	klog.Info("命名空间检查完成......")

	err = createConfigMap(o.Namespace, o.ConfigMapName)
	defer func() {
		if err != nil {
			tool.DestoryConfigMap(o.Namespace, o.ConfigMapName)
		}
	}()
	klog.Infof("Postgres Operator的ConfigMap对象创建完成......")

	err = createOperatorSA(o.Namespace, o.ServiceAccountName)
	defer func() {
		if err != nil {
			tool.DestoryServiceAccount(o.Namespace, o.ServiceAccountName)
		}
	}()
	klog.Infof("Postgres Operator的ServiceAccount对象创建完成......")

	err = createClusterRole(o.ClusterRoleName)
	defer func() {
		if err != nil {
			tool.DestoryClusterRole(o.ClusterRoleName)
		}
	}()
	klog.Infof("Postgres Operator的ClusterRole对象创建完成......")

	err = createClusterRoleBinding(o.ClusterRoleBindingName, o.ClusterRoleName, o.ServiceAccountName, o.Namespace)
	defer func() {
		if err != nil {
			tool.DestoryClusterRoleBinding(o.ClusterRoleBindingName)
		}
	}()
	klog.Infof("Postgres Operator的ClusterRoleinding对象创建完成......")

	err = createPodClusterRole(o.PodClusterRoleName)
	defer func() {
		if err != nil {
			tool.DestoryClusterRole(o.PodClusterRoleName)
		}
	}()
	klog.Infof("Postgres Operator的Pod ClusterRole对象创建完成......")

	err = createDeployment(o.DeploymentName, o.Image, o.ServiceAccountName, o.Namespace)
	defer func() {
		if err != nil {
			tool.DestoryDeployment(o.Namespace, o.DeploymentName)
		}
	}()
	klog.Infof("Postgres Operator的Deployment对象创建完成......")

	err = createService(o.ServiceName, o.Namespace)
	defer func() {
		if err != nil {
			tool.DestoryService(o.Namespace, o.ServiceName)
		}
	}()
	klog.Infof("Postgres Operator的API Service对象创建完成......")

	return nil
}

func (o *PostgresOperator) Delete() (err error) {
	klog.Info("开始删除postgres operator对象......")

	err = tool.DestoryService(o.Namespace, o.ServiceName)
	if err != nil {
		return
	}
	klog.Info("删除postgres operator的API service对象.")

	err = tool.DestoryDeployment(o.Namespace, o.DeploymentName)
	if err != nil {
		return
	}
	klog.Info("删除postgres operator的deployment对象.")

	err = tool.DestoryClusterRole(o.PodClusterRoleName)
	if err != nil {
		return
	}
	klog.Info("删除postgres operator的Pod ClusterRole对象.")

	err = tool.DestoryClusterRoleBinding(o.ClusterRoleBindingName)
	if err != nil {
		return
	}
	klog.Info("删除postgres operator的ClusterRoleBinding对象.")

	err = tool.DestoryClusterRole(o.ClusterRoleName)
	if err != nil {
		return
	}
	klog.Info("删除postgres operator的ClusterRole对象.")

	err = tool.DestoryServiceAccount(o.Namespace, o.ServiceAccountName)
	if err != nil {
		return
	}
	klog.Info("删除postgres operator的serviceaccount对象.")

	err = tool.DestoryConfigMap(o.Namespace, o.ConfigMapName)
	if err != nil {
		return
	}
	klog.Info("删除postgres operator的configmap对象.")

	klog.Info("成功删除postgres operator......")

	return
}

func createConfigMap(namespace, name string) (err error) {
	cm := newConfigMap(name)
	if _, err := client.MyClient.K8sClientset.CoreV1().ConfigMaps(namespace).Create(context.TODO(), cm, metav1.CreateOptions{}); err != nil {
		klog.Errorf("创建Postgres Operator的ConfigMap对象[%s:%s]失败.", namespace, name, err)
		return err
	}

	return
}

func newConfigMap(name string) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Data: map[string]string{
			"api_port":                     "8080",
			"aws_region":                   "eu-central-1",
			"cluster_domain":               "cluster.local",
			"cluster_history_entries":      "1000",
			"cluster_labels":               "application:spilo",
			"cluster_name_label":           "cluster-name",
			"connection_pooler_image":      "cuijx/pgbouncer:master-8",
			"db_hosted_zone":               "db.example.com",
			"debug_logging":                "true",
			"docker_image":                 "cuijx/spilo-12:1.6-p3",
			"enable_crd_validation":        "false",
			"enable_master_load_balancer":  "false",
			"enable_replica_load_balancer": "false",
			"enable_teams_api":             "false",
			"logical_backup_docker_image":  "cuijx/logical-backup",
			"logical_backup_s3_bucket":     "my-bucket-url",
			"logical_backup_s3_sse":        "AES256",
			"logical_backup_schedule":      "30 00 * * *",
			"master_dns_name_format":       "{cluster}.{team}.{hostedzone}",
			"pdb_name_format":              "postgres-{cluster}-pdb",
			"pod_deletion_wait_timeout":    "10m",
			"pod_label_wait_timeout":       "10m",
			"pod_management_policy":        "ordered_ready",
			"pod_role_label":               "spilo-role",
			"pod_service_account_name":     "postgres-pod",
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
			"workers":                      "8",
		},
	}
}

func createOperatorSA(namespace, name string) (err error) {
	account := tool.NewServiceAccount(name)
	if _, err := client.MyClient.K8sClientset.CoreV1().ServiceAccounts(namespace).Create(context.TODO(), account, metav1.CreateOptions{}); err != nil {
		klog.Errorf("创建Postgres Operator的ServiceAccount对象[%s:%s]失败.", namespace, name, err)
		return err
	}
	return
}

func createClusterRole(name string) (err error) {
	cr := newClusterRole(name)
	if _, err := client.MyClient.K8sClientset.RbacV1().ClusterRoles().Create(context.TODO(), cr, metav1.CreateOptions{}); err != nil {
		klog.Errorf("创建Postgres Operator的ClusterRole对象[%s]失败.", name, err)
		return err
	}
	return
}

func newClusterRole(name string) *rbacv1.ClusterRole {
	return &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Rules: []rbacv1.PolicyRule{
			{
				// all verbs allowed for custom operator resources
				APIGroups: []string{"acid.zalan.do"},
				Resources: []string{"postgresqls", "postgresqls/status", "operatorconfigurations"},
				Verbs:     []string{"create", "delete", "deletecollection", "get", "list", "patch", "update", "watch"},
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
				// to send events to the CRs
				APIGroups: []string{""},
				Resources: []string{"events"},
				Verbs:     []string{"create", "get", "list", "patch", "update", "watch"},
			},
			{
				// to manage endpoints which are also used by Patroni
				APIGroups: []string{""},
				Resources: []string{"endpoints"},
				Verbs:     []string{"create", "delete", "deletecollection", "get", "list", "patch", "update", "watch"},
			},
			{
				// to CRUD secrets for database access
				APIGroups: []string{""},
				Resources: []string{"secrets"},
				Verbs:     []string{"create", "get", "delete", "update"},
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
				Verbs:     []string{"get", "list", "update"},
			},
			{
				// to watch Spilo pods and do rolling updates. Creation via StatefulSet
				APIGroups: []string{""},
				Resources: []string{"pods"},
				Verbs:     []string{"delete", "get", "list", "patch", "update", "watch"},
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
				Verbs:     []string{"create", "delete", "get", "patch", "update"},
			},
			{
				// to CRUD the StatefulSet which controls the Postgres cluster instances
				APIGroups: []string{"apps"},
				Resources: []string{"statefulsets", "deployments"},
				Verbs:     []string{"create", "delete", "get", "list", "patch"},
			},
			{
				// to CRUD cron jobs for logical backups
				APIGroups: []string{""},
				Resources: []string{"cronjobs"},
				Verbs:     []string{"create", "delete", "get", "patch", "update"},
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
				// to create role bindings to the postgres-pod service account
				APIGroups: []string{"rbac.authorization.k8s.io"},
				Resources: []string{"rolebindings"},
				Verbs:     []string{"get", "create"},
			},
			{
				// to grant privilege to run privileged pods
				APIGroups:     []string{"extensions"},
				Resources:     []string{"podsecuritypolicies"},
				ResourceNames: []string{"privileged"},
				Verbs:         []string{"use"},
			},
		},
	}
}

func createClusterRoleBinding(clusterRoleBindingName, clusterRoleName, serviceAccountName, namespace string) (err error) {
	crb := newClusterRoleBinding(clusterRoleBindingName, clusterRoleName, serviceAccountName, namespace)
	if _, err := client.MyClient.K8sClientset.RbacV1().ClusterRoleBindings().Create(context.TODO(), crb, metav1.CreateOptions{}); err != nil {
		klog.Errorf("创建Postgres Operator的ClusterRoleBinding对象[%s]失败.", clusterRoleBindingName, err)
		return err
	}
	return
}

func newClusterRoleBinding(clusterRoleBindingName, clusterRoleName, serviceAccountName, namespace string) *rbacv1.ClusterRoleBinding {
	return &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: clusterRoleBindingName,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      serviceAccountName,
				Namespace: namespace,
			},
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     clusterRoleName,
		},
	}
}

func createPodClusterRole(podClusterRoleName string) (err error) {
	cr := newPodClusterRole(podClusterRoleName)
	if _, err := client.MyClient.K8sClientset.RbacV1().ClusterRoles().Create(context.TODO(), cr, metav1.CreateOptions{}); err != nil {
		klog.Errorf("创建Postgres Operator的Pod ClusterRole对象[%s]失败.", podClusterRoleName, err)
		return err
	}
	return
}

func newPodClusterRole(podClusterRoleName string) *rbacv1.ClusterRole {
	return &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: podClusterRoleName,
		},
		Rules: []rbacv1.PolicyRule{
			{
				// Patroni needs to watch and manage endpoints
				APIGroups: []string{""},
				Resources: []string{"endpoints"},
				Verbs:     []string{"create", "delete", "deletecollection", "get", "list", "patch", "update", "watch"},
			},
			{
				// Patroni needs to watch pods
				APIGroups: []string{""},
				Resources: []string{"pods"},
				Verbs:     []string{"get", "list", "patch", "update", "watch"},
			},
			{
				// to let Patroni create a headless service
				APIGroups: []string{""},
				Resources: []string{"services"},
				Verbs:     []string{"create"},
			},
			{
				// to run privileged pods
				APIGroups:     []string{"extensions"},
				Resources:     []string{"podsecuritypolicies"},
				ResourceNames: []string{"privileged"},
				Verbs:         []string{"use"},
			},
		},
	}
}

func createDeployment(deploymentName, image, serviceAccountName, namespace string) (err error) {
	d := newDeployment(namespace, image, deploymentName, serviceAccountName)
	if _, err = client.MyClient.K8sClientset.AppsV1().Deployments(namespace).Create(context.TODO(), d, metav1.CreateOptions{}); err != nil {
		return
	}
	return
}

func newDeployment(namespace, image, deploymentName, serviceAccountName string) *appsv1.Deployment {
	d := tool.CreateBasicDeployment(namespace, deploymentName, "name", 1)

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
					Value: "postgres-operator",
				},
			},
		},
	}

	spec := corev1.PodSpec{
		ServiceAccountName: serviceAccountName,
		Containers:         containers,
	}

	d.Spec.Template.Spec = spec

	return d
}

func createService(name, namespace string) (err error) {
	s := newService(name)
	_, err = client.MyClient.K8sClientset.CoreV1().Services(namespace).Create(context.TODO(), s, metav1.CreateOptions{})
	if err != nil {
		return
	}
	return
}

func newService(name string) *corev1.Service {
	servicePort := corev1.ServicePort{Port: 8080, Protocol: corev1.ProtocolTCP, TargetPort: intstr.IntOrString{IntVal: 8080}}
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{servicePort},
			Selector: map[string]string{
				"name": "postgres-operator",
			},
			Type: corev1.ServiceTypeClusterIP,
		},
	}
}
