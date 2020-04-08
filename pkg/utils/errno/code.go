package errno

var (
	// Common errors
	OK                  = &Errno{Code: 0, Message: "OK"}
	InternalServerError = &Errno{Code: 100001, Message: "Internal server error"}
	ErrBind             = &Errno{Code: 100002, Message: "Error occurred while binding the request body to the struct."}
	ErrBadParam         = &Errno{Code: 100003, Message: "Bad Parameters."}
	ErrValidation       = &Errno{Code: 100004, Message: "Validation failed."}
	ErrDatabase         = &Errno{Code: 100005, Message: "Database error."}
	ErrToken            = &Errno{Code: 100006, Message: "Error occurred while signing the JSON web token."}

	// user errors
	ErrEncrypt           = &Errno{Code: 100101, Message: "加密用户密码时发生错误！"}
	ErrUserNotFound      = &Errno{Code: 100102, Message: "The user was not found."}
	ErrTokenInvalid      = &Errno{Code: 100103, Message: "Token 无效！"}
	ErrPasswordIncorrect = &Errno{Code: 100104, Message: "用户密码无效！"}

	ErrBadK8sConfig         = &Errno{Code: 200001, Message: "Kubernetes config err."}
	ErrCreateK8sClientSet   = &Errno{Code: 200002, Message: "Kubernete clientset init err."}
	ErrCreatePgsqlClientSet = &Errno{Code: 200003, Message: "Pgsql clientset init err."}
	ErrCreateApiClientSet   = &Errno{Code: 200004, Message: "Kubernetes api clientset init err."}
	ErrCreateRedisClientSet = &Errno{Code: 200005, Message: "Redis clientset init err."}
	ErrCreateMySQLClientSet = &Errno{Code: 200006, Message: "创建MySQL Clientset 对象失败！"}
	ErrUpGraderRequest      = &Errno{Code: 200020, Message: "升级get请求为websocket协议失败."}

	ErrCreateServiceAccount     = &Errno{Code: 200102, Message: "Create serviceaccount failed."}
	ErrCreateClusterRole        = &Errno{Code: 200103, Message: "Create clustrrole failed."}
	ErrCreateClusterRoleBinding = &Errno{Code: 200104, Message: "Crate clusterrolebinding failed."}
	ErrCreateDeployment         = &Errno{Code: 200105, Message: "Create deployment failed."}

	// Postgresql cluster.
	ErrCreatePostgresCluster  = &Errno{Code: 200201, Message: "Create postgres cluster failed."}
	ErrDeletePostgresCluster  = &Errno{Code: 200202, Message: "Delete postgres cluster failed."}
	ErrGetPostgresCluster     = &Errno{Code: 200203, Message: "Get postgres cluster failed."}
	ErrGetPostgresClusterList = &Errno{Code: 200204, Message: "Get postgres cluster list failed."}

	// RedisFailover cluster.
	ErrCreateRedisFailoverCluster  = &Errno{Code: 200211, Message: "Create redis failover cluster failed."}
	ErrDeleteRedisFailoverCluster  = &Errno{Code: 200212, Message: "Delete redis failover cluster failed."}
	ErrGetRedisFailoverCluster     = &Errno{Code: 200213, Message: "Get redis cluster failed."}
	ErrGetRedisFailoverClusterList = &Errno{Code: 200214, Message: "Get redis cluster list failed."}

	// MySQL cluster.
	ErrCreateMySQLOperator = &Errno{Code: 200220, Message: "创建MySQL Operator组件失败！"}
	ErrCreateMySQLCluster  = &Errno{Code: 200221, Message: "创建MySQL集群失败！"}
	ErrDeleteMySQLCluster  = &Errno{Code: 200222, Message: "删除MySQL集群失败！"}
	ErrGetMySQLCluster     = &Errno{Code: 200223, Message: "获取MySQL集群信息失败！"}
	ErrGetMySQLClusterList = &Errno{Code: 200224, Message: "获取MySQL集群列表失败！"}
	ErrMySQLRBACCheck      = &Errno{Code: 200225, Message: "创建MySQL集群RBAC对象失败！"}

	// Traefik
	ErrCreateTraefikAddon = &Errno{Code: 200230, Message: "安装Traefik插件失败！"}

	ErrCreatePersistentVolumeClaim  = &Errno{Code: 200401, Message: "Create persistent volume claim failed."}
	ErrDeletePersistentVolumeClaim  = &Errno{Code: 200402, Message: "Delete persistent volume claim failed."}
	ErrGetPersistentVolumeClaim     = &Errno{Code: 200403, Message: "Get persistent volume claim failed."}
	ErrGetPersistentVolumeClaimList = &Errno{Code: 200404, Message: "Get persistent volume claim list failed."}

	ErrCreateJob      = &Errno{Code: 200411, Message: "Create job failed."}
	ErrDeleteJob      = &Errno{Code: 200412, Message: "Delete job failed."}
	ErrGetJob         = &Errno{Code: 200413, Message: "Get job failed."}
	ErrGetJobList     = &Errno{Code: 200414, Message: "Get job list failed."}
	ErrGetJobPodsList = &Errno{Code: 200415, Message: "Get job pods list failed."}

	ErrDeleteDeployment      = &Errno{Code: 200422, Message: "Delete deployment failed."}
	ErrGetDeployment         = &Errno{Code: 200423, Message: "Get deployment failed."}
	ErrGetDeploymentList     = &Errno{Code: 200424, Message: "Get deployment list failed."}
	ErrGetDeploymentPodsList = &Errno{Code: 200425, Message: "Get deployment pods list failed."}

	ErrDeleteService      = &Errno{Code: 200432, Message: "Delete service failed."}
	ErrGetService         = &Errno{Code: 200433, Message: "Get service failed."}
	ErrGetServiceList     = &Errno{Code: 200434, Message: "Get service list failed."}
	ErrGetServicePodsList = &Errno{Code: 200435, Message: "Get service pods list failed."}

	ErrGetStorageClass     = &Errno{Code: 200443, Message: "Get storageclass failed."}
	ErrGetStorageClassList = &Errno{Code: 200444, Message: "Get storageclass list failed."}

	ErrCreateCronJob  = &Errno{Code: 200451, Message: "Create cron job failed."}
	ErrGetCronJob     = &Errno{Code: 200452, Message: "Get cron job failed."}
	ErrGetCronJobList = &Errno{Code: 200453, Message: "Get cron job list failed."}
	ErrDeleteCronJob  = &Errno{Code: 200454, Message: "Delete cron job failed."}

	ErrCreateSecret  = &Errno{Code: 200461, Message: "Create secret failed."}
	ErrGetSecret     = &Errno{Code: 200462, Message: "Get secret failed."}
	ErrGetSecretList = &Errno{Code: 200463, Message: "Get secret list failed."}
	ErrDeleteSecret  = &Errno{Code: 200464, Message: "Delete secret failed."}

	ErrCreateConfigMap    = &Errno{Code: 200471, Message: "Create configmap failed."}
	ErrGetConfigMapDetail = &Errno{Code: 200472, Message: "Get configmap failed."}
	ErrGetConfigMapList   = &Errno{Code: 200473, Message: "Get configmap list failed."}
	ErrDeleteConfigMap    = &Errno{Code: 200474, Message: "Delete configmap failed."}

	ErrGetPodDetail     = &Errno{Code: 200482, Message: "Get pod detail failed."}
	ErrGetPodList       = &Errno{Code: 200483, Message: "Get pod list failed."}
	ErrGetPodContainers = &Errno{Code: 200484, Message: "Get pod containers failed."}
	ErrGetPodLogs       = &Errno{Code: 200485, Message: "Get pod logs failed."}

	ErrCreateCloneCodeJob = &Errno{Code: 201010, Message: "Create clone code job failed."}

	ErrCreateBuildImageJob = &Errno{Code: 201020, Message: "Create build image pod failed."}

	ErrDeployAtomService     = &Errno{Code: 201030, Message: "Create atom service failed."}
	ErrScaleDeployment       = &Errno{Code: 201031, Message: "Scale deployment pods count failed."}
	ErrUpdateDeploymentImage = &Errno{Code: 201032, Message: "Update deployment image failed."}
)
