package container

import (
	"encoding/json"
	"hello-k8s/pkg/kubernetes/client"
	"hello-k8s/pkg/kubernetes/kuberesource/resource/container"
	"hello-k8s/pkg/kubernetes/kuberesource/resource/logs"
	"hello-k8s/pkg/utils/errno"
	"hello-k8s/pkg/utils/tool"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/lexkong/log"
)

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// @Summary 获取某一 Container 对象的 Logs.
// @Description 获取某一 Container 对象的 Logs.
// @Tags resource
// @Param namespace path string true "命名空间"
// @Param podId path string true "PodID"
// @Param containerId path string true "Container"
// @Success 200 {object} tool.Response "{"code":200,"message":"OK","data":{""}}"
// @Router /resource/container/logs/{namespace}/{podId}/{containerId} [get]
func GetLogs(c *gin.Context) {
	log.Debug("获取某一 Container 对象的 Logs.")

	namespace := c.Param("namespace")
	podID := c.Param("podId")
	containerID := c.Param("containerId")
	if namespace == "" || podID == "" || containerID == "" {
		tool.SendResponse(c, errno.ErrBadParam, nil)
	}

	// Init kubernetes client
	clientset, err := client.New()
	if err != nil {
		tool.SendResponse(c, errno.ErrCreateK8sClientSet, nil)
		return
	}

	// 升级 get 请求为 webSocket 协议
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Errorf(err, "升级 get 请求为 webSocket 协议失败")
		tool.SendResponse(c, errno.ErrUpGraderRequest, err)
		return
	}

	defer ws.Close()
	//读取ws中的数据
	mt, _, err := ws.ReadMessage()
	newestLogTimestamp := "newset"
	for {
		if strings.EqualFold(strings.ToLower(newestLogTimestamp), strings.ToLower("newset")) {
			podLogs, err := container.GetLogDetails(clientset, namespace, podID,
				containerID, logs.AllSelection, false)
			if err != nil {
				log.Errorf(err, "获取pod[%s:%s:%s]的日志失败!", namespace, podID, containerID)
				break
			}
			newestLogTimestamp = string(podLogs.LogLines[len(podLogs.LogLines)-1].Timestamp)

			data, err := json.Marshal(podLogs)

			err = ws.WriteMessage(mt, data)

		} else {
			selection := &logs.Selection{
				ReferencePoint: logs.LogLineId{
					LogTimestamp: logs.LogTimestamp(newestLogTimestamp),
					LineNum:      1,
				},
				OffsetFrom: 1,
				OffsetTo:   2,
			}

			podLogs, err := container.GetLogDetails(clientset, namespace, podID,
				containerID, selection, false)
			if err != nil {
				log.Errorf(err, "获取pod[%s:%s:%s]的日志失败!", namespace, podID, containerID)
				break
			}
			if len(podLogs.LogLines) == 0 {
				continue
			}
			tmpTimestamp := string(podLogs.LogLines[len(podLogs.LogLines)-1].Timestamp)

			if newestLogTimestamp == tmpTimestamp {
				continue
			} else {
				newestLogTimestamp = tmpTimestamp
			}

			data, err := json.Marshal(podLogs)

			err = ws.WriteMessage(mt, data)
		}

		// time.Sleep(time.Duration(time.Millisecond * 10))
	}
}
