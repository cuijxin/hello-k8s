package tool

import (
	"context"
	"encoding/base64"
	"hello-k8s/pkg/kubernetes/client"
	deploy "hello-k8s/pkg/kubernetes/kuberesource/resource/deployment"
	"hello-k8s/pkg/kubernetes/kuberesource/resource/statefulset"
	"hello-k8s/pkg/model"
	"hello-k8s/pkg/model/common"
	"hello-k8s/pkg/utils/errno"
	"net/http"
	"path"
	"strconv"
	"time"
	"unsafe"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/teris-io/shortid"
	"github.com/unknwon/com"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/kube"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func SendResponse(c *gin.Context, err error, data interface{}) {
	code, message := errno.DecodeErr(err)

	// always return http.StatusOK
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

func ActionConfigInit(namespace string) (*action.Configuration, error) {
	actionConfig := new(action.Configuration)
	clientConfig := kube.GetConfig(settings.KubeConfig, settings.KubeContext)
}

func CreateNamespace(namespace string, clientset kubernetes.Interface) {
	_, err := clientset.CoreV1().Namespaces().Get(context.TODO(), namespace, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		ns := &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: namespace,
			},
		}
		clientset.CoreV1().Namespaces().Create(context.TODO(), ns, metav1.CreateOptions{})
	}
}

func NewServiceAccount(serviceAccountName string) *corev1.ServiceAccount {
	return &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name: serviceAccountName,
		},
	}
}

func CreateBasicDeployment(namespace, name, labelKey string, replicas int32) *appsv1.Deployment {
	count := int32(replicas)
	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: appsv1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				labelKey: name,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &count,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					labelKey: name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						labelKey: name,
					},
				},
			},
		},
	}
}

func CreatePodSpec(name string, podSpecArgs model.PodArgs) *corev1.PodSpec {
	containerSpec := corev1.Container{
		Name:  name,
		Image: podSpecArgs.ContainerImage,
		Resources: corev1.ResourceRequirements{
			Requests: make(map[corev1.ResourceName]resource.Quantity),
		},
		Env: ConvertEnvVarsSpec(podSpecArgs.Variables),
	}

	if len(podSpecArgs.ConfigMaps) > 0 {
		for _, configMapObj := range podSpecArgs.ConfigMaps {
			configFileName := path.Base(configMapObj.MountPath)
			volumeMount := corev1.VolumeMount{
				Name:      configMapObj.Name,
				MountPath: configMapObj.MountPath,
				SubPath:   configFileName,
				ReadOnly:  configMapObj.ReadOnly,
			}
			containerSpec.VolumeMounts = append(containerSpec.VolumeMounts, volumeMount)
		}
	}

	if len(podSpecArgs.PersistentVolumeClaims) > 0 {
		for _, pvc := range podSpecArgs.PersistentVolumeClaims {
			volumeMount := corev1.VolumeMount{
				Name:      pvc.Name,
				MountPath: pvc.MountPath,
				ReadOnly:  pvc.ReadOnly,
			}
			containerSpec.VolumeMounts = append(containerSpec.VolumeMounts, volumeMount)
		}
	}

	if podSpecArgs.ContainerCommand != nil {
		for _, cmd := range podSpecArgs.ContainerCommand {
			containerSpec.Command = append(containerSpec.Command, cmd)
		}
	}

	if podSpecArgs.ContainerCommandArgs != nil {
		for _, arg := range podSpecArgs.ContainerCommandArgs {
			containerSpec.Args = append(containerSpec.Args, arg)
		}
	}

	if podSpecArgs.CpuRequirement > 0 {
		capacity := strconv.FormatFloat(podSpecArgs.CpuRequirement, 'f', 5, 32)
		request, _ := resource.ParseQuantity(capacity)
		containerSpec.Resources.Requests[corev1.ResourceCPU] = request
	}
	if podSpecArgs.MemoryRequirement > 0 {
		in := strconv.FormatFloat(podSpecArgs.MemoryRequirement, 'f', 5, 32)
		capacity := in + viper.GetString("constants.storage_unit")
		request, _ := resource.ParseQuantity(capacity)
		containerSpec.Resources.Requests[corev1.ResourceMemory] = request
	}

	podSpec := corev1.PodSpec{
		Containers:    []corev1.Container{containerSpec},
		RestartPolicy: podSpecArgs.RestartPolicy,
	}

	if len(podSpecArgs.ConfigMaps) > 0 {
		for _, configMapObj := range podSpecArgs.ConfigMaps {
			configmap := corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: configMapObj.Name,
				},
			}
			volume := corev1.Volume{
				Name: configMapObj.Name,
				VolumeSource: corev1.VolumeSource{
					ConfigMap: &configmap,
				},
			}
			podSpec.Volumes = append(podSpec.Volumes, volume)
		}
	}

	if len(podSpecArgs.PersistentVolumeClaims) > 0 {
		for _, pvc := range podSpecArgs.PersistentVolumeClaims {
			volume := corev1.Volume{
				Name: pvc.Name,
				VolumeSource: corev1.VolumeSource{
					PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
						ClaimName: pvc.Name,
					},
				},
			}
			podSpec.Volumes = append(podSpec.Volumes, volume)
		}
	}

	return &podSpec
}

func ConvertEnvVarsSpec(variables []deploy.EnvironmentVariable) []corev1.EnvVar {
	var result []corev1.EnvVar
	for _, variable := range variables {
		result = append(result, corev1.EnvVar{Name: variable.Name, Value: variable.Value})
	}
	return result
}

// Converts array of labels to map[string]string
func GetLabelsMap(labels []deploy.Label) map[string]string {
	result := make(map[string]string)

	for _, label := range labels {
		result[label.Key] = label.Value
	}

	return result
}

// String Convert []byte object to string.
func String(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// InitCustomResourceDefinition 初始化一个CustomResourceDefinition对象的定义.
func InitCustomResourceDefinition(name, kind, singular, plural, group, version string) *v1beta1.CustomResourceDefinition {
	return &v1beta1.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: v1beta1.CustomResourceDefinitionSpec{
			Group:   group,
			Version: version,
			Scope:   v1beta1.NamespaceScoped,
			Names: v1beta1.CustomResourceDefinitionNames{
				Kind:     kind,
				Singular: singular,
				Plural:   plural,
			},
		},
	}
}

// InitClusterRoleBinding 初始化一个ClusterRoleBinding对象.
func InitClusterRoleBinding(subjectKind, RoleRefKind, name, serviceaccountName, clusterRoleName, namespace string) *rbacv1.ClusterRoleBinding {
	return &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      subjectKind,
				Name:      serviceaccountName,
				Namespace: namespace,
			},
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: rbacv1.GroupName,
			Kind:     RoleRefKind,
			Name:     clusterRoleName,
		},
	}
}

// CheckMySQLClusterRBAC 检查MySQL集群RBAC对象的函数.
func CheckMySQLClusterRBAC(namespace, serviceaccountName, roleName, clusterRoleName string) error {
	klog.Info("check mysql serviceaccount object.")
	_, err := client.MyClient.K8sClientset.CoreV1().ServiceAccounts(namespace).Get(context.TODO(), serviceaccountName, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		klog.Info("create mysql agent serviceaccount object.")
		sa := NewServiceAccount(serviceaccountName)
		if _, err = client.MyClient.K8sClientset.CoreV1().ServiceAccounts(namespace).Create(context.TODO(), sa, metav1.CreateOptions{}); err != nil {
			return err
		}
	}

	klog.Info("check mysql rolebinding object.")
	_, err = client.MyClient.K8sClientset.RbacV1().RoleBindings(namespace).Get(context.TODO(), roleName, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		klog.Info("create mysql agent rolebinding object.")
		rb := &rbacv1.RoleBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name:      roleName,
				Namespace: namespace,
			},
			RoleRef: rbacv1.RoleRef{
				APIGroup: rbacv1.GroupName,
				Kind:     "ClusterRole",
				Name:     clusterRoleName,
			},
		}
		rb.Subjects = append(rb.Subjects, rbacv1.Subject{
			Kind:      "ServiceAccount",
			Name:      serviceaccountName,
			Namespace: namespace,
		})
		if _, err = client.MyClient.K8sClientset.RbacV1().RoleBindings(namespace).Create(context.TODO(), rb, metav1.CreateOptions{}); err != nil {
			return err
		}
	}
	return nil
}

func NewVolumeClaimTemplate(name string, r *common.DataVolumeArg) *corev1.PersistentVolumeClaim {
	template := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
	for _, accessMode := range r.AccessModes {
		var mode corev1.PersistentVolumeAccessMode
		switch accessMode {
		case "ReadWriteOnce":
			mode = corev1.ReadWriteOnce
		case "ReadOnlyMany":
			mode = corev1.ReadOnlyMany
		case "ReadWriteMany":
			mode = corev1.ReadWriteMany
		default:
			mode = corev1.ReadWriteMany
		}
		template.Spec.AccessModes = append(template.Spec.AccessModes, mode)
	}

	template.Spec.StorageClassName = r.StorageClassName

	in := strconv.FormatFloat(r.Capacity, 'f', 5, 32)
	capacity := in + viper.GetString("constants.storage_unit")
	request, _ := resource.ParseQuantity(capacity)
	template.Spec.Resources = corev1.ResourceRequirements{
		Requests: corev1.ResourceList{
			corev1.ResourceStorage: request,
		},
	}
	return template
}

func WaitForStatefulsetReady(name, namespace string) (err error) {
	err = wait.PollImmediate(100*time.Millisecond, 10*time.Minute, func() (done bool, err error) {
		detail, err := statefulset.GetStatefulSetDetail(client.MyClient.K8sClientset, nil, namespace, name)
		if err != nil {
			klog.Errorf("获取Statefulset对象[%s:%s]详情失败: %v", namespace, name, err)
			return false, err
		}
		if detail.Pods.Running > 0 && *detail.Pods.Desired == detail.Pods.Running {
			return true, nil
		}
		return false, nil
	})
	if err != nil {
		return err
	}
	return nil
}

func GetPage(c *gin.Context) int {
	result := 0
	page, _ := com.StrTo(c.Query("page")).Int()
	if page > 0 {
		result = (page - 1) * viper.GetInt("app.page_size")
	}
	return result
}

func GenShortId() (string, error) {
	return shortid.Generate()
}

func GetReqID(c *gin.Context) string {
	v, ok := c.Get("X-Request-Id")
	if !ok {
		return ""
	}
	if requestId, ok := v.(string); ok {
		return requestId
	}
	return ""
}

func NewConfigMap(name string, r *common.ConfigMapArg) *corev1.ConfigMap {
	klog.Info("创建用户自定义配置文件对应的 ConfigMap 对象.")
	configmap := corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}

	if len(r.Items) > 0 {
		tmp := make(map[string]string)
		for _, item := range r.Items {
			d, _ := base64.StdEncoding.DecodeString(item.Value)
			str := String(d)
			tmp[item.Key] = str
		}
		configmap.Data = tmp
	}

	return &configmap
}

func NewSecret(name string, r *common.CustomRootPasswordArg) *corev1.Secret {
	secret := corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}

	if len(r.SecretValue) > 0 {
		tmp := make(map[string]string)
		tmp["password"] = r.SecretValue

		secret.StringData = tmp
	}

	secret.Type = "Opaque"

	return &secret
}

func DestoryConfigMap(namespace, name string) (err error) {
	err = client.MyClient.K8sClientset.CoreV1().ConfigMaps(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		klog.Errorf("删除configmap对象[%s:%s]失败.", namespace, name, err)
		return
	}
	return
}

func DestorySecret(namespace, name string) (err error) {
	err = client.MyClient.K8sClientset.CoreV1().Secrets(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		klog.Errorf("删除secret对象[%s:%s]失败.", namespace, name, err)
		return
	}
	return
}

func DestoryService(namespace, name string) (err error) {
	err = client.MyClient.K8sClientset.CoreV1().Services(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		klog.Errorf("删除service对象[%s:%s]失败.", namespace, name, err)
		return
	}
	return
}

func DestoryPVC(namespace, name string) (err error) {
	err = client.MyClient.K8sClientset.CoreV1().PersistentVolumeClaims(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		klog.Errorf("删除pvc对象[%s:%s]失败.", namespace, name, err)
		return
	}
	return
}

func DestoryServiceAccount(namespace, name string) (err error) {
	err = client.MyClient.K8sClientset.CoreV1().ServiceAccounts(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		klog.Errorf("删除serviceAccount对象[%s:%s]失败.", namespace, name, err)
		return
	}
	return
}

func DestoryClusterRole(name string) (err error) {
	err = client.MyClient.K8sClientset.RbacV1().ClusterRoles().Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		klog.Errorf("删除clusterRole对象[%s]失败.", name, err)
		return
	}
	return
}

func DestoryClusterRoleBinding(name string) (err error) {
	err = client.MyClient.K8sClientset.RbacV1().ClusterRoleBindings().Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		klog.Errorf("删除clusterRoleBinding对象[%s]失败.", name, err)
		return
	}
	return nil
}

func DestoryDeployment(namespace, name string) (err error) {
	err = client.MyClient.K8sClientset.AppsV1().Deployments(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		klog.Errorf("删除deployment对象[%s:%s]失败.", namespace, name, err)
		return
	}
	return
}
