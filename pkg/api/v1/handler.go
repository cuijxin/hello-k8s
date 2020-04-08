package v1

import (
	"hello-k8s/pkg/handler/resources/common"
	deploy "hello-k8s/pkg/kubernetes/kuberesource/resource/deployment"
	"hello-k8s/pkg/utils/errno"
	"net/http"
	"path"
	"strconv"
	"unsafe"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"github.com/spf13/viper"

	appsv1 "k8s.io/api/apps/v1"
	api "k8s.io/api/core/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
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

func CreateNamespace(namespace string, clientset kubernetes.Interface) {
	_, err := clientset.CoreV1().Namespaces().Get(namespace, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		ns := &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: namespace,
			},
		}
		clientset.CoreV1().Namespaces().Create(ns)
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

func CreatePodSpec(name string, podSpecArgs common.PodArgs) *corev1.PodSpec {
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
			volumeMount := api.VolumeMount{
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
			volumeMount := api.VolumeMount{
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
			configmap := api.ConfigMapVolumeSource{
				LocalObjectReference: api.LocalObjectReference{
					Name: configMapObj.Name,
				},
			}
			volume := api.Volume{
				Name: configMapObj.Name,
				VolumeSource: api.VolumeSource{
					ConfigMap: &configmap,
				},
			}
			podSpec.Volumes = append(podSpec.Volumes, volume)
		}
	}

	if len(podSpecArgs.PersistentVolumeClaims) > 0 {
		for _, pvc := range podSpecArgs.PersistentVolumeClaims {
			volume := api.Volume{
				Name: pvc.Name,
				VolumeSource: api.VolumeSource{
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
func CheckMySQLClusterRBAC(namespace, serviceaccountName, roleName, clusterRoleName string, clientset kubernetes.Interface) error {
	log.Debugf("check mysql serviceaccount object.")
	_, err := clientset.CoreV1().ServiceAccounts(namespace).Get(serviceaccountName, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		log.Debugf("create mysql agent serviceaccount object.")
		sa := NewServiceAccount(serviceaccountName)
		if _, err = clientset.CoreV1().ServiceAccounts(namespace).Create(sa); err != nil {
			return err
		}
	}

	log.Debugf("check mysql rolebinding object.")
	_, err = clientset.RbacV1().RoleBindings(namespace).Get(roleName, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		log.Debugf("create mysql agent rolebinding object.")
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
		if _, err = clientset.RbacV1().RoleBindings(namespace).Create(rb); err != nil {
			return err
		}
	}
	return nil
}
