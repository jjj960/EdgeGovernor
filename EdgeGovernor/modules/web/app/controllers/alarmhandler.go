package controllers

import (
	"EdgeGovernor/pkg/models"
	"EdgeGovernor/pkg/utils"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

func GetTaskLog(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"taskLogs":  models.Msgs,
		"pageSizes": len(models.Msgs),
	})
}

func DeleteTaskLog(c echo.Context) error {
	sids := c.FormValue("sids")
	// 分割字符串为切片
	sidSlice := strings.Split(sids, ",")
	newTaskLogs := make([]models.Msg, 0, len(models.Msgs))
	for _, taskLog := range models.Msgs {
		if !contains(sidSlice, taskLog.ID) {
			newTaskLogs = append(newTaskLogs, taskLog)
		}
	}
	models.Msgs = newTaskLogs
	return c.String(http.StatusOK, "Delete successful")
}

func DeleteAllLog(c echo.Context) error {
	models.Msgs = []models.Msg{}
	return c.String(http.StatusOK, "Delete all successful")
}

func contains(s []string, e int64) bool {
	for _, a := range s {
		if i, _ := utils.StringtoInt64(a); i == e {
			return true
		}
	}
	return false
}

func SSE_Connect(c echo.Context) error {
	clientID := c.Param("client_id")
	// 调用generateSSE函数开始SSE通信
	return generateSSE(clientID, c)
}

// generateSSE 生成SSE消息
func generateSSE(clientID string, c echo.Context) error {
	// 设置响应头信息
	c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
	c.Response().Header().Set(echo.HeaderCacheControl, "no-cache")
	c.Response().Header().Set("Connection", "keep-alive")

	for {
		select {
		case msg := <-utils.AlarmMsgChannel:
			// 从 AlarmMsgChannel 接收消息
			jsonData, err := utils.Jsoniter.Marshal(msg)
			if err != nil {
				c.Logger().Error(err)
				continue
			}
			// 构造SSE消息并发送给客户端
			_, err = fmt.Fprintf(c.Response().Writer, "data: %s\n\n", jsonData)
			if err != nil {
				c.Logger().Error(err)
				continue
			}
			// 刷新响应，确保数据发送到客户端
			c.Response().Flush()
			models.Msgs = append(models.Msgs, msg)
		case <-c.Request().Context().Done():
			// 客户端断开连接，退出循环
			return nil
		}
	}
}
