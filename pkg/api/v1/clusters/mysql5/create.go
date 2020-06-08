package mysql5

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/klog"

	"hello-k8s/pkg/kubernetes/client"
	"hello-k8s/pkg/kubernetes/kuberesource/resource/secret"
	"hello-k8s/pkg/kubernetes/kuberesource/resource/service"
	"hello-k8s/pkg/model/common"
	"hello-k8s/pkg/storage/database"
	"hello-k8s/pkg/utils/errno"
	"hello-k8s/pkg/utils/tool"

	v1 "github.com/cuijxin/mysql-operator/pkg/apis/mysql/v1"
	"github.com/cuijxin/mysql-operator/pkg/constants"
)

// @Summary 创建MySQL V5集群
// @Description 创建MySQL V5集群
// @Tags add-on
// @Accept json
// @Produce json
// @param data body mysql5.ClusterOptions true "创建 MySQL V5集群所需参数."
// @Success 200 {object} tool.Response "{"code":0,"message":"OK","data":{""}}"
// @Router /v1/addon/mysql5/cluster/create [post]
func CreateCluster(c *gin.Context) {
	klog.Info("调用创建 MySQL V5 集群函数")

	var r ClusterOptions
	if err := c.BindJSON(&r); err != nil {
		tool.SendResponse(c, errno.InternalServerError, err)
	}
	klog.Infof("request body is [%v]", r)

	prepareClusterOptions(&r)

	isExist, err := database.DB.Exist(database.RecordOptions{
		Name:      r.Name,
		Namespace: r.Namespace,
		ClusterID: r.ClusterID,
		Type:      common.MySQLV5,
	})
	if err != nil {
		klog.Errorf("mongo数据库操作失败: ", err)
		tool.SendResponse(c, errno.InternalServerError, err)
		return
	}
	if isExist {
		klog.Info("同名记录已经存在")
		tool.SendResponse(c, errno.InternalServerError, errors.New("同名记录已经存在"))
	}

	tool.CreateNamespace(r.Namespace, client.MyClient.K8sClientset)
	klog.Info("命名空间检查执行完毕......")
	tool.CheckMySQLClusterRBAC(r.Namespace, "mysql-agent", "mysql-agent", "mysql5-operator")
	klog.Info("检查mysql rbac对象完毕......")

	// 处理用户自定义配置文件
	if r.Config != nil {
		cm := tool.NewConfigMap(r.Name+"-cnf", r.Config)
		_, err := client.MyClient.K8sClientset.CoreV1().ConfigMaps(r.Namespace).Create(context.TODO(), cm, metav1.CreateOptions{})
		if err != nil {
			tool.SendResponse(c, errno.InternalServerError, err)
			return
		}
		defer func() {
			if err != nil {
				tool.DestoryConfigMap(r.Namespace, r.Name+"-cnf")
			}
		}()
	}

	// 处理用户自定义密码
	if r.RootPassword != nil {
		s := tool.NewSecret(r.Name+"-root-user-secret", r.RootPassword)
		_, err := client.MyClient.K8sClientset.CoreV1().Secrets(r.Namespace).Create(context.TODO(), s, metav1.CreateOptions{})
		if err != nil {
			tool.SendResponse(c, errno.InternalServerError, err)
			return
		}
		defer func() {
			if err != nil {
				tool.DestorySecret(r.Namespace, r.Name+"-root-user-secret")
			}
		}()
	}

	if r.IsExport {
		svc := newMySQLNodePortService(r.Namespace, r.Name)
		_, err := client.MyClient.K8sClientset.CoreV1().Services(r.Namespace).Create(context.TODO(), svc, metav1.CreateOptions{})
		if err != nil {
			tool.SendResponse(c, errno.InternalServerError, err)
			return
		}
		defer func() {
			if err != nil {
				tool.DestoryService(r.Namespace, r.Name+"-public")
			}
		}()
	}

	mc := newMySQLCluster(r)
	_, err = client.MyClient.Mysql5Client.MysqlV1().MySQLClusters(r.Namespace).Create(context.TODO(), mc, metav1.CreateOptions{})
	if err != nil {
		tool.SendResponse(c, errno.InternalServerError, err)
		return
	}
	defer func() {
		if err != nil {
			DestoryMySQLCluster(r.Namespace, r.Name)
		}
	}()

	var nodePort int32
	err = wait.PollImmediate(100*time.Millisecond, 8*time.Second, func() (done bool, err error) {
		svcName := fmt.Sprintf("%s-public", r.Name)
		detail, err := service.GetServiceDetail(client.MyClient.K8sClientset, r.Namespace, svcName)
		if err != nil {
			return false, nil
		}
		if detail.InternalEndpoint.Ports != nil {
			for _, item := range detail.InternalEndpoint.Ports {
				nodePort = item.NodePort
				return true, nil
			}
		}
		return false, nil
	})
	svcIP := viper.GetString("constants.public_ip")
	data := newStoreData(common.MySQLV5, r, svcIP, nodePort)
	if err := database.DB.Store(data); err != nil {
		klog.Errorf("存储插件对象信息失败：%v", err)
		tool.SendResponse(c, errno.InternalServerError, err)
		return
	}

	mysqlCluster := ClusterInfo{
		Name:         r.Name,
		Namespace:    r.Namespace,
		ClusterID:    r.ClusterID,
		Host:         r.Name,
		RootUserName: "root",
		Port:         nodePort,
		Domain:       svcIP,
	}

	if r.InitDBName != nil {
		mysqlCluster.InitDBName = *r.InitDBName
	}

	if r.RootPassword != nil {
		mysqlCluster.RootPassword = r.RootPassword.SecretValue
	} else {

		err := wait.PollImmediate(100*time.Millisecond, 3*time.Minute, func() (done bool, err error) {
			secret, err := secret.GetSecretDetail(client.MyClient.K8sClientset, r.Namespace, fmt.Sprintf("%s-root-password", r.Name))
			if err != nil {
				klog.Info("重试获取密钥的方法......")
				return false, nil
			}
			mysqlCluster.RootPassword = tool.String(secret.Data["password"])
			return true, nil
		})
		if err != nil {
			tool.SendResponse(c, err, nil)
			return
		}
	}

	if err := tool.WaitForStatefulsetReady(r.Name, r.Namespace); err != nil {
		klog.Errorf("获取Statefulset对象 [%s:%s] 状态失败.", r.Namespace, r.Name, err)
		tool.SendResponse(c, errno.InternalServerError, nil)
		return
	}

	tool.SendResponse(c, errno.OK, mysqlCluster)

	return
}

func DestoryMySQLCluster(namespace, name string) {
	client.MyClient.Mysql5Client.MysqlV1().MySQLClusters(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
}

func prepareClusterOptions(r *ClusterOptions) {
	var replica int32 = 1
	if r.Members == nil {
		klog.Info("用户没有设置members字段,设置为默认值1")
		r.Members = &replica
	} else if *r.Members <= 0 && *r.Members > 10 {
		klog.Info("用户设置的members字段超出有效范围,设置为默认值1")
		r.Members = &replica
	}
}

func newMySQLNodePortService(namespace, name string) *corev1.Service {
	mysqlPort := corev1.ServicePort{Port: 3306}
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Labels:    map[string]string{constants.MySQLClusterLabel: name},
			Name:      fmt.Sprintf("%s-public", name),
			Namespace: namespace,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{mysqlPort},
			Selector: map[string]string{
				constants.MySQLClusterLabel: name,
			},
			Type: corev1.ServiceTypeNodePort,
		},
	}
	return svc
}

func newMySQLCluster(r ClusterOptions) *v1.MySQLCluster {
	mysqlCluster := v1.MySQLCluster{
		ObjectMeta: metav1.ObjectMeta{
			Name: r.Name,
		},
	}

	spec := v1.MySQLClusterSpec{}
	if r.Members != nil {
		spec.Replicas = *r.Members
	}

	if r.Config != nil {
		config := &corev1.LocalObjectReference{
			Name: r.Name + "-cnf",
		}
		spec.ConfigRef = config
	}

	if r.RootPassword != nil {
		secret := &corev1.LocalObjectReference{
			Name: r.Name + "-root-user-secret",
		}
		spec.SecretRef = secret
	}

	if r.DataVolume != nil {
		template := tool.NewVolumeClaimTemplate("data", r.DataVolume)
		spec.VolumeClaimTemplate = template
	}

	if r.BackupVolume != nil {
		template := tool.NewVolumeClaimTemplate("backup", r.BackupVolume)
		spec.BackupVolumeClaimTemplate = template
	}

	if r.InitDBName != nil {
		spec.InitDBName = r.InitDBName
	}

	mysqlCluster.Spec = spec

	return &mysqlCluster
}

func newStoreData(appType common.AppType, r ClusterOptions, serviceDomain string, port int32) *common.AtomApplication {
	data := &common.AtomApplication{
		Name:      r.Name,
		Namespace: r.Namespace,
		ClusterID: r.ClusterID,
		Type:      appType,
	}

	nodePortServiceName := fmt.Sprintf("%s-public", r.Name)
	mysqlData := common.MySQLAddOnData{
		StatefulsetName: &r.Name,
		ServiceName:     &r.Name,
		// IngressRouteName:    &r.Name,
		NodePortServiceName: &nodePortServiceName,
		Port:                &port,
	}

	mysqlData.ServiceDomain = &serviceDomain

	if r.Members != nil {
		mysqlData.Members = r.Members
	} else {
		var defaultReplicas int32 = 1
		mysqlData.Members = &defaultReplicas
	}

	if r.Config != nil {
		cm := fmt.Sprintf("%s-cnf", r.Name)
		mysqlData.ConfigMapName = &cm
	}
	if r.RootPassword != nil {
		s := fmt.Sprintf("%s-root-user-secret", r.Name)
		mysqlData.RootPasswordSecretName = &s
	}
	if r.DataVolume != nil {
		var count, i int32
		count = *mysqlData.Members
		volumeNames := make([]string, count)
		for i = 0; i < count; i++ {
			volumeNames[i] = fmt.Sprintf("data-%s-%d", r.Name, i)
		}
		mysqlData.DataVolumeName = volumeNames
	}
	if r.BackupVolume != nil {
		var count, i int32
		count = *mysqlData.Members
		volumeNames := make([]string, count)
		for i = 0; i < count; i++ {
			volumeNames[i] = fmt.Sprintf("backup-%s-%d", r.Name, i)
		}
		mysqlData.BackupVolumeName = volumeNames
	}

	data.MySQLAddon = &mysqlData

	return data
}
