// Copyright 2017 The Kubernetes Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package deployment

import (
	"fmt"
	"io"
	"log"
	"path"
	"strconv"
	"strings"

	"hello-k8s/pkg/kubernetes/kuberesource/errors"

	"github.com/spf13/viper"
	apps "k8s.io/api/apps/v1"
	api "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/rand"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	client "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	// DescriptionAnnotationKey is annotation key for a description.
	DescriptionAnnotationKey = "description"
)

// AppDeploymentSpec is a specification for an app deployment.
type AppDeploymentSpec struct {
	// Name of the application.
	Name string `json:"name"`

	// Docker image path for the application.
	ContainerImage string `json:"containerImage"`

	// The name of an image pull secret in case of a private docker repository.
	ImagePullSecret *string `json:"imagePullSecret"`

	// Command that is executed instead of container entrypoint, if specified.
	ContainerCommand *string `json:"containerCommand"`

	// Arguments for the specified container command or container entrypoint (if command is not
	// specified here).
	ContainerCommandArgs *string `json:"containerCommandArgs"`

	// Number of replicas of the image to maintain.
	Replicas int32 `json:"replicas"`

	// Port mappings for the service that is created. The service is created if there is at least
	// one port mapping.
	PortMappings []PortMapping `json:"portMappings"`

	// List of user-defined environment variables.
	Variables []EnvironmentVariable `json:"variables"`

	// List of user-defined configmap variables.
	ConfigMaps []ConfigVariable `json:"configmaps"`

	// List of user-defined PersistentVolumeClaim variables.
	PersistentVolumeClaims []PersistentVolumeClaimVariable `json:"pvcs"`

	// Whether the created service is external.
	IsExternal bool `json:"isExternal"`

	// Whether the created service is Loadbalancer type
	IsLoadBalancer bool `json:"isLoadBalancer"`

	// Description of the deployment.
	Description *string `json:"description"`

	// Target namespace of the application.
	Namespace string `json:"namespace"`

	// // Optional memory requirement for the container.
	// MemoryRequirement *resource.Quantity `json:"memoryRequirement"`
	//
	// // Optional CPU requirement for the container.
	// CpuRequirement *resource.Quantity `json:"cpuRequirement"`

	// Optional memory requirement for the container.
	MemoryRequirement float64 `json:"memoryRequirement"`

	// Optional CPU requirement for the container.
	CpuRequirement float64 `json:"cpuRequirement"`

	// Labels that will be defined on Pods/RCs/Services
	Labels []Label `json:"labels"`

	// Whether to run the container as privileged user (essentially equivalent to root on the host).
	RunAsPrivileged bool `json:"runAsPrivileged"`
}

// AppDeploymentFromFileSpec is a specification for deployment from file
type AppDeploymentFromFileSpec struct {
	// Name of the file
	Name string `json:"name"`

	// Namespace that object should be deployed in
	Namespace string `json:"namespace"`

	// File content
	Content string `json:"content"`

	// Whether validate content before creation or not
	Validate bool `json:"validate"`
}

// AppDeploymentFromFileResponse is a specification for deployment from file
type AppDeploymentFromFileResponse struct {
	// Name of the file
	Name string `json:"name"`

	// File content
	Content string `json:"content"`

	// Error after create resource
	Error string `json:"error"`
}

// PortMapping is a specification of port mapping for an application deployment.
type PortMapping struct {
	// Port that will be exposed on the service.
	Port int32 `json:"port"`

	// Docker image path for the application.
	TargetPort int32 `json:"targetPort"`

	// IP protocol for the mapping, e.g., "TCP" or "UDP".
	Protocol api.Protocol `json:"protocol"`
}

// ConfigVariable 代表一个服务运行时所需的配置文件参数
type ConfigVariable struct {
	// Name 一个ConfigMap对象的名称，必须是与应用同命名空间下的一个可用ConfigMap对象的名称.
	Name string `json:"name"`

	// MountPath 配置文件挂载路径，代表这个服务如果想要成功运行，需要到那个路径下去获取这个配置文件.
	MountPath string `json:"mountPath"`

	// ReadOnly
	ReadOnly bool `json:"readOnly"`
}

// PersistentVolumeClaimVariable 代表一个服务运行时需要的持久化存储参数.
type PersistentVolumeClaimVariable struct {
	// Name 一个 PersistentVolumeClaim 对象的名称，必须是与应用同命名空间下的一个可用 PersistentVolumeClaim 对象的名称.
	Name string `json:"name"`

	// MountPath 持久化存储挂载路径.
	MountPath string `json:"mountPath"`

	// ReadOnly
	ReadOnly bool `json:"readOnly"`
}

// EnvironmentVariable represents a named variable accessible for containers.
type EnvironmentVariable struct {
	// Name of the variable. Must be a C_IDENTIFIER.
	Name string `json:"name"`

	// Value of the variable, as defined in Kubernetes core API.
	Value string `json:"value"`
}

// Label is a structure representing label assignable to Pod/RC/Service
type Label struct {
	// Label key
	Key string `json:"key"`

	// Label value
	Value string `json:"value"`
}

// Protocols is a structure representing supported protocol types for a service
type Protocols struct {
	// Array containing supported protocol types e.g., ["TCP", "UDP"]
	Protocols []api.Protocol `json:"protocols"`
}

// DeployApp deploys an app based on the given configuration. The app is deployed using the given
// client. App deployment consists of a deployment and an optional service. Both of them
// share common labels.
func DeployApp(spec *AppDeploymentSpec, client client.Interface) error {
	log.Printf("Deploying %s application into %s namespace", spec.Name, spec.Namespace)

	annotations := map[string]string{}
	if spec.Description != nil {
		annotations[DescriptionAnnotationKey] = *spec.Description
	}
	labels := getLabelsMap(spec.Labels)
	objectMeta := metaV1.ObjectMeta{
		Annotations: annotations,
		Name:        spec.Name,
		Labels:      labels,
	}

	containerSpec := api.Container{
		Name:  spec.Name,
		Image: spec.ContainerImage,
		SecurityContext: &api.SecurityContext{
			Privileged: &spec.RunAsPrivileged,
		},
		Resources: api.ResourceRequirements{
			Requests: make(map[api.ResourceName]resource.Quantity),
		},
		Env: convertEnvVarsSpec(spec.Variables),
	}

	if len(spec.ConfigMaps) > 0 {
		for _, configMapObj := range spec.ConfigMaps {
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

	if len(spec.PersistentVolumeClaims) > 0 {
		for _, pvc := range spec.PersistentVolumeClaims {
			volumeMount := api.VolumeMount{
				Name:      pvc.Name,
				MountPath: pvc.MountPath,
				ReadOnly:  pvc.ReadOnly,
			}
			containerSpec.VolumeMounts = append(containerSpec.VolumeMounts, volumeMount)
		}
	}

	if spec.ContainerCommand != nil {
		containerSpec.Command = []string{*spec.ContainerCommand}
	}
	if spec.ContainerCommandArgs != nil {
		containerSpec.Args = []string{*spec.ContainerCommandArgs}
	}

	// if spec.CpuRequirement != nil {
	// 	containerSpec.Resources.Requests[api.ResourceCPU] = *spec.CpuRequirement
	// }
	// if spec.MemoryRequirement != nil {
	// 	containerSpec.Resources.Requests[api.ResourceMemory] = *spec.MemoryRequirement
	// }

	if spec.CpuRequirement > 0 {
		capacity := strconv.FormatFloat(spec.CpuRequirement, 'f', 5, 32)
		request, _ := resource.ParseQuantity(capacity)
		containerSpec.Resources.Requests[api.ResourceCPU] = request
	}
	if spec.MemoryRequirement > 0 {
		in := strconv.FormatFloat(spec.MemoryRequirement, 'f', 5, 32)
		capacity := in + viper.GetString("constants.storage_unit")
		request, _ := resource.ParseQuantity(capacity)
		containerSpec.Resources.Requests[api.ResourceMemory] = request
	}

	podSpec := api.PodSpec{
		Containers: []api.Container{containerSpec},
	}
	if spec.ImagePullSecret != nil {
		podSpec.ImagePullSecrets = []api.LocalObjectReference{{Name: *spec.ImagePullSecret}}
	}

	if len(spec.ConfigMaps) > 0 {
		for _, configMapObj := range spec.ConfigMaps {
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

	if len(spec.PersistentVolumeClaims) > 0 {
		for _, pvc := range spec.PersistentVolumeClaims {
			volume := api.Volume{
				Name: pvc.Name,
				VolumeSource: api.VolumeSource{
					PersistentVolumeClaim: &api.PersistentVolumeClaimVolumeSource{
						ClaimName: pvc.Name,
					},
				},
			}
			podSpec.Volumes = append(podSpec.Volumes, volume)
		}
	}

	podTemplate := api.PodTemplateSpec{
		ObjectMeta: objectMeta,
		Spec:       podSpec,
	}

	deployment := &apps.Deployment{
		ObjectMeta: objectMeta,
		Spec: apps.DeploymentSpec{
			Replicas: &spec.Replicas,
			Template: podTemplate,
			Selector: &metaV1.LabelSelector{
				// Quoting from https://kubernetes.io/docs/concepts/workloads/controllers/deployment/#selector:
				// In API version apps/v1beta2, .spec.selector and .metadata.labels no longer default to
				// .spec.template.metadata.labels if not set. So they must be set explicitly.
				// Also note that .spec.selector is immutable after creation of the Deployment in apps/v1beta2.
				MatchLabels: labels,
			},
		},
	}
	_, err := client.AppsV1().Deployments(spec.Namespace).Create(deployment)

	if err != nil {
		return err
	}

	if len(spec.PortMappings) > 0 {
		service := &api.Service{
			ObjectMeta: objectMeta,
			Spec: api.ServiceSpec{
				Selector: labels,
			},
		}

		if spec.IsExternal {
			if spec.IsLoadBalancer {
				service.Spec.Type = api.ServiceTypeLoadBalancer
			} else {
				service.Spec.Type = api.ServiceTypeNodePort
			}
		} else {
			service.Spec.Type = api.ServiceTypeClusterIP
		}

		for _, portMapping := range spec.PortMappings {
			servicePort :=
				api.ServicePort{
					Protocol: portMapping.Protocol,
					Port:     portMapping.Port,
					Name:     generatePortMappingName(portMapping),
					TargetPort: intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: portMapping.TargetPort,
					},
				}
			service.Spec.Ports = append(service.Spec.Ports, servicePort)
		}

		_, err = client.CoreV1().Services(spec.Namespace).Create(service)
		return err
	}

	return nil
}

// GetAvailableProtocols returns list of available protocols. Currently it is TCP and UDP.
func GetAvailableProtocols() *Protocols {
	return &Protocols{Protocols: []api.Protocol{api.ProtocolTCP, api.ProtocolUDP}}
}

func convertEnvVarsSpec(variables []EnvironmentVariable) []api.EnvVar {
	var result []api.EnvVar
	for _, variable := range variables {
		result = append(result, api.EnvVar{Name: variable.Name, Value: variable.Value})
	}
	return result
}

func generatePortMappingName(portMapping PortMapping) string {
	return generateName(fmt.Sprintf("%s-%d-%d-", strings.ToLower(string(portMapping.Protocol)),
		portMapping.Port, portMapping.TargetPort))
}

func generateName(base string) string {
	maxNameLength := 63
	randomLength := 5
	maxGeneratedNameLength := maxNameLength - randomLength
	if len(base) > maxGeneratedNameLength {
		base = base[:maxGeneratedNameLength]
	}
	return fmt.Sprintf("%s%s", base, rand.String(randomLength))
}

// Converts array of labels to map[string]string
func getLabelsMap(labels []Label) map[string]string {
	result := make(map[string]string)

	for _, label := range labels {
		result[label.Key] = label.Value
	}

	return result
}

// DeployAppFromFile deploys an app based on the given yaml or json file.
func DeployAppFromFile(cfg *rest.Config, spec *AppDeploymentFromFileSpec) (bool, error) {
	reader := strings.NewReader(spec.Content)
	log.Printf("Namespace for deploy from file: %s\n", spec.Namespace)
	d := yaml.NewYAMLOrJSONDecoder(reader, 4096)
	for {
		data := unstructured.Unstructured{}
		if err := d.Decode(&data); err != nil {
			if err == io.EOF {
				return true, nil
			}
			return false, err
		}

		version := data.GetAPIVersion()
		kind := data.GetKind()

		gv, err := schema.ParseGroupVersion(version)
		if err != nil {
			gv = schema.GroupVersion{Version: version}
		}

		discoveryClient, err := discovery.NewDiscoveryClientForConfig(cfg)
		if err != nil {
			return false, err
		}

		apiResourceList, err := discoveryClient.ServerResourcesForGroupVersion(version)
		if err != nil {
			return false, err
		}
		apiResources := apiResourceList.APIResources
		var resource *metaV1.APIResource
		for _, apiResource := range apiResources {
			if apiResource.Kind == kind && !strings.Contains(apiResource.Name, "/") {
				resource = &apiResource
				break
			}
		}
		if resource == nil {
			return false, fmt.Errorf("unknown resource kind: %s", kind)
		}

		dynamicClient, err := dynamic.NewForConfig(cfg)
		if err != nil {
			return false, err
		}

		groupVersionResource := schema.GroupVersionResource{Group: gv.Group, Version: gv.Version, Resource: resource.Name}

		if strings.Compare(spec.Namespace, "_all") == 0 {
			_, err = dynamicClient.Resource(groupVersionResource).Namespace(data.GetNamespace()).Create(&data, metaV1.CreateOptions{})
		} else {
			_, err = dynamicClient.Resource(groupVersionResource).Namespace(spec.Namespace).Create(&data, metaV1.CreateOptions{})
		}

		if err != nil {
			return false, errors.LocalizeError(err)
		}
	}
}
