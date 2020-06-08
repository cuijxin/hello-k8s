package mysql5

import (
	"context"
	"hello-k8s/pkg/kubernetes/client"
	"hello-k8s/pkg/model/common"
	"hello-k8s/pkg/storage/database"
	"hello-k8s/pkg/utils/errno"
	"hello-k8s/pkg/utils/tool"

	"github.com/gin-gonic/gin"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"
)

// @Summary 删除指定的MySQL V5集群
// @Description 删除指定的MySQL V5集群
// @Tags add-on
// @Accept json
// @produce json
// @param data body mysql5.DeleteClusterOptions true "删除参数"
// @Success 200 {object} tool.Response "{"code":200,"message":"OK","data":{""}}"
// @Router /v1/addon/mysql5/cluster/delete [delete]
func DeleteCluster(c *gin.Context) {
	klog.Info("调用删除MySQL集群的函数.")

	var r DeleteClusterOptions
	if err := c.BindJSON(&r); err != nil {
		tool.SendResponse(c, errno.ErrBind, err)
		return
	}

	klog.Infof("request body is: [name:%s],[namespace:%s],[clusterId:%s]", r.Name, r.Namespace, r.ClusterID)

	mysqlData, err := database.DB.Get(database.RecordOptions{
		Name:      r.Name,
		Namespace: r.Namespace,
		ClusterID: r.ClusterID,
		Type:      common.MySQLV5,
	})
	if err != nil {
		klog.Errorf("获取存储的mysql插件信息失败：%v", err)
		tool.SendResponse(c, errno.InternalServerError, err)
		return
	}

	if err := client.MyClient.Mysql5Client.MysqlV1().MySQLClusters(r.Namespace).Delete(context.TODO(), r.Name, metav1.DeleteOptions{}); err != nil {
		klog.Errorf("Delete MySQL Cluster object failed.", err)
		tool.SendResponse(c, errno.InternalServerError, err)
		return
	}

	// 处理衍生资源对象的删除
	if mysqlData.MySQLAddon.ConfigMapName != nil {
		name := *mysqlData.MySQLAddon.ConfigMapName
		klog.Infof("Delete custom configmap [%s]", name)
		if err := client.MyClient.K8sClientset.CoreV1().ConfigMaps(r.Namespace).Delete(context.TODO(), name, metav1.DeleteOptions{}); err != nil {
			klog.Errorf("删除用户自定义configmap对象 [%s] failed.", name, err)
			tool.SendResponse(c, errno.InternalServerError, err)
			return
		}
	}
	if mysqlData.MySQLAddon.RootPasswordSecretName != nil {
		name := *mysqlData.MySQLAddon.RootPasswordSecretName
		klog.Infof("Delete root user secret [%s]", name)
		if err := client.MyClient.K8sClientset.CoreV1().Secrets(r.Namespace).Delete(context.TODO(), name, metav1.DeleteOptions{}); err != nil {
			klog.Errorf("删除用户自定义secret对象 [%s] falied.", name, err)
			tool.SendResponse(c, errno.ErrDeleteSecret, err)
			return
		}
	}
	if mysqlData.MySQLAddon.DataVolumeName != nil {
		for _, item := range mysqlData.MySQLAddon.DataVolumeName {
			klog.Infof("删除用户自定义pvc对象 [%s]", item)
			if err := tool.DestoryPVC(r.Namespace, item); err != nil {
				klog.Infof("删除用户自定义pvc对象 [%s] 失败.", item, err)
				tool.SendResponse(c, errno.InternalServerError, err)
				return
			}
		}
	}
	if mysqlData.MySQLAddon.BackupVolumeName != nil {
		for _, item := range mysqlData.MySQLAddon.BackupVolumeName {
			klog.Infof("删除用户自定义pvc对象 [%s]", item)
			if err := tool.DestoryPVC(r.Namespace, item); err != nil {
				klog.Infof("删除用户自定义pvc对象 [%s] 失败.", item, err)
				tool.SendResponse(c, errno.InternalServerError, err)
				return
			}
		}
	}

	if mysqlData.MySQLAddon.NodePortServiceName != nil {
		name := *mysqlData.MySQLAddon.NodePortServiceName
		klog.Infof("删除mysql集群创建的NodePort类型的service对象 [%s]", name)
		if err := client.MyClient.K8sClientset.CoreV1().Services(r.Namespace).Delete(context.TODO(), name, metav1.DeleteOptions{}); err != nil {
			klog.Errorf("删除mysql集群创建的NodePort类型的service对象 [%s] 失败.", name, err)
			tool.SendResponse(c, errno.InternalServerError, err)
			return
		}
	}

	if err := database.DB.Delete(database.RecordOptions{
		Name:      r.Name,
		Namespace: r.Namespace,
		ClusterID: r.ClusterID,
		Type:      common.MySQLV5,
	}); err != nil {
		klog.Errorf("删除MySQL集群失败: %v", err)
		tool.SendResponse(c, errno.InternalServerError, err)
		return
	}

	tool.SendResponse(c, errno.OK, nil)

	return
}
