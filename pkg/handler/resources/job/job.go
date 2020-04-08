package job

import (
	"hello-k8s/pkg/handler/resources/common"
)

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
	PodTemplate common.PodArgs `json:"podTemplate"`
}

// CreateJobRequest 定义了创建一个Job对象时所需参数.
type CreateJobRequest struct {
	// Name Job 对象名称.
	Name string `json:"name"`

	// Namespace 命名空间.
	Namespace string `json:"namespace"`

	// Job 对象参数.
	Job JobArgs `json:"jobTemplate"`
}

// DeleteJobRequest 定义了删除Job对象时所需参数
type DeleteJobRequest struct {
	// Job 对象名称.
	Name string `json:"name"`

	// Namespace 命名空间.
	Namespace string `json:"namespace"`
}
