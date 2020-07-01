package postgres

import (
	"context"
	"fmt"
	"hello-k8s/pkg/kubernetes/client"
	"hello-k8s/pkg/model/common"
	"hello-k8s/pkg/storage/database"
	"hello-k8s/pkg/utils/errno"
	"hello-k8s/pkg/utils/tool"

	"github.com/gin-gonic/gin"
	"k8s.io/klog"

	acidv1 "github.com/cuijxin/postgres-operator-atom/pkg/apis/acid.zalan.do/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// @Summary 创建 Postgresql 集群
// @Description 创建 Postgresql 集群
// @Tags add-on
// @Accept json
// @Produce json
// @param data body postgres.CreateClusterOptions true "创建 Postgresql 集群参数"
// @Success 200 {object} tool.Response "{"code":200,"message":"OK","data":{""}}"
// @Router /v1/addon/postgresql/cluster/create [post]
func CreateCluster(c *gin.Context) {
	klog.Info("调用创建 Postgresql 集群的函数")

	var r CreateClusterOptions
	if err := c.BindJSON(&r); err != nil {
		tool.SendResponse(c, errno.ErrBind, err)
		return
	}
	klog.Infof("request body is [%v]", r)

	isExist, err := database.DB.Exist(database.RecordOptions{
		Name:      r.Name,
		Namespace: r.Namespace,
		ClusterID: r.ClusterID,
		Type:      common.PostgreSQL,
	})
	if err != nil {
		tool.SendResponse(c, errno.InternalServerError, err)
		return
	}
	if isExist {
		tool.SendResponse(c, errno.InternalServerError, err)
		return
	}

	tool.CreateNamespace(r.Namespace, client.MyClient.K8sClientset)
	klog.Info("命名空间检查完成......")
	pgc := newPostgresCluster(r)
	_, err = client.MyClient.PgClient.AcidV1().Postgresqls(r.Namespace).Create(context.TODO(), pgc, metav1.CreateOptions{})
	if err != nil {
		tool.SendResponse(c, errno.InternalServerError, err)
		return
	}

	tool.SendResponse(c, errno.OK, nil)
}

func newPostgresCluster(r CreateClusterOptions) *acidv1.Postgresql {
	spec := acidv1.PostgresSpec{
		Resources: acidv1.Resources{
			ResourceRequests: acidv1.ResourceDescription{
				CPU:    "100m",
				Memory: "250Mi",
			},
		},
		Users:     initUsers(r.Users),
		Databases: r.Databases,
	}

	if r.TeamID != nil {
		spec.TeamID = *r.TeamID
	} else {
		spec.TeamID = "atom"
	}

	if r.Replicas != nil {
		spec.NumberOfInstances = *r.Replicas
	} else {
		spec.NumberOfInstances = 2
	}

	if r.Volume != nil {
		volume := acidv1.Volume{
			Size:         r.Volume.Size,
			StorageClass: r.Volume.StorageClass,
		}
		spec.Volume = volume
	} else {
		volume := acidv1.Volume{
			Size: "1Gi",
		}
		spec.Volume = volume
	}

	if r.PgParams != nil && r.PgParams.PgVersion != nil {
		param := acidv1.PostgresqlParam{
			PgVersion: *r.PgParams.PgVersion,
		}
		spec.PostgresqlParam = param
	} else {
		param := acidv1.PostgresqlParam{
			PgVersion: "11",
		}
		spec.PostgresqlParam = param
	}

	pgsql := acidv1.Postgresql{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-%s", spec.TeamID, r.Name),
			Namespace: r.Namespace,
		},
	}

	pgsql.Spec = spec

	return &pgsql
}

func initUsers(users map[string]UserFlags) map[string]acidv1.UserFlags {
	result := map[string]acidv1.UserFlags{}
	for key, value := range users {
		var tmp acidv1.UserFlags
		for _, str := range value {
			tmp = append(tmp, str)
		}
		result[key] = tmp
	}
	return result
}

func newStoreData(addonType common.AppType, r CreateClusterOptions, serviceDomain string) *common.AtomApplication {
	data := &common.AtomApplication{
		Name:      r.Name,
		Namespace: r.Namespace,
		ClusterID: r.ClusterID,
		Type:      addonType,
	}

	var secretName, teamID string
	if r.TeamID != nil {
		secretName = fmt.Sprintf("postgres.%s-%s.credentials", r.TeamID, r.Name)
		teamID = *r.TeamID
	} else {
		secretName = fmt.Sprintf("postgres.atom-%s.credentials", r.Name)
		teamID = "atom"
	}

	pgsqlData := common.PgSQLData{
		TeamID:                &teamID,
		DefaultRootSecretName: &secretName,
		IngressRouteName:      &r.Name,
		ServiceDomain:         &serviceDomain,
	}

	data.PgSQLAddon = &pgsqlData

	return data
}
