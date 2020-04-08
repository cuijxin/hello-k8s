package model

import (
	"sync"
	"time"

	deploy "hello-k8s/pkg/kubernetes/kuberesource/resource/deployment"

	corev1 "k8s.io/api/core/v1"
)

type BaseModel struct {
	ID        uint64     `gorm:"primary_key;AUTO_INCREMENT;column:id" json:"-"`
	CreatedAt time.Time  `gorm:"column:createdAt" json:"-"`
	UpdatedAt time.Time  `gorm:"column:updatedAt" json:"-"`
	DeletedAt *time.Time `gorm:"column:deletedAt" sql:"index" json:"-"`
}

type UserInfo struct {
	ID        uint64 `json:"id"`
	UserName  string `json:"username"`
	SayHello  string `json:"sayHello"`
	Password  string `json:"password"`
	CreatedAt string `json:"createdAt"`
	UpdateAt  string `json:"updatedAt"`
}

type UserList struct {
	Lock  *sync.Mutex
	IdMap map[uint64]*UserInfo
}

// Token represents a JSON web token.
type Token struct {
	Token string `json:"token"`
}

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

// JobArgs 定义了构建一个 Job 对象时所需参数.
type JobArgs struct {
	// Specifies the maximum desired number of pods the job should
	// run at any given time. The actual number of pods running in steady state will
	// be less than this number when ((.spec.completions - .status.successful) < .spec.parallelism),
	// i.e. when the work left to do is less than max parallelism.
	// More info: https://kubernetes.io/docs/concepts/workloads/controllers/jobs-run-to-completion/
	// +optional
	Parallelism *int32 `json:"parallelism,omitempty" protobuf:"varint,1,opt,name=parallelism"`

	// Specifies the desired number of successfully finished pods the
	// job should be run with.  Setting to nil means that the success of any
	// pod signals the success of all pods, and allows parallelism to have any positive
	// value.  Setting to 1 means that parallelism is limited to 1 and the success of that
	// pod signals the success of the job.
	// More info: https://kubernetes.io/docs/concepts/workloads/controllers/jobs-run-to-completion/
	// +optional
	Completions *int32 `json:"completions,omitempty" protobuf:"varint,2,opt,name=completions"`

	// Specifies the duration in seconds relative to the startTime that the job may be active
	// before the system tries to terminate it; value must be positive integer
	// +optional
	ActiveDeadlineSeconds *int64 `json:"activeDeadlineSeconds,omitempty" protobuf:"varint,3,opt,name=activeDeadlineSeconds"`

	// PodTemplate 定义了 Job 对象管理的 Pod 对象的定义参数.
	PodTemplate PodArgs `json:"podTemplate"`
}
