package common

import (
	deploy "hello-k8s/pkg/kubernetes/kuberesource/resource/deployment"

	corev1 "k8s.io/api/core/v1"
)

type PodArgs struct {
	// Docker image path for the application.
	ContainerImage string `json:"containerImage"`

	// Command that is executed instead of container entrypoint, if specified.
	ContainerCommand []string `json:"containerCommand"`

	// Arguments for the specified container command or container entrypoint (if command is not
	// specified here).
	ContainerCommandArgs []string `json:"containerCommandArgs"`

	// List of user-defined environment variables.
	Variables []deploy.EnvironmentVariable `json:"variables"`

	// Optional memory requirement for the container.
	MemoryRequirement float64 `json:"memoryRequirement"`

	// Optional CPU requirement for the container.
	CpuRequirement float64 `json:"cpuRequirement"`

	// Restart policy for all containers within the pod.
	// One of Always, OnFailure, Never.
	RestartPolicy corev1.RestartPolicy `json:"restartPolicy"`

	// Labels that will be defined on Pods/RCs/Services
	Labels []deploy.Label `json:"labels"`

	// List of user-defined configmap variables.
	ConfigMaps []deploy.ConfigVariable `json:"configmaps"`

	// List of user-defined PersistentVolumeClaim variables.
	PersistentVolumeClaims []deploy.PersistentVolumeClaimVariable `json:"pvcs"`
}
