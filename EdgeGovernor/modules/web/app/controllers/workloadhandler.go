package controllers

import (
	"EdgeGovernor/modules/web/app/dao"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

type getDataRequest struct {
	Data string `json:"data"`
}

type getDataRespone struct {
	NodeName  []string `json:"nodeName"`
	NodeValue []string `json:"nodeValue"`
	Total     float64  `json:"total"`
}

func GetBarData(c echo.Context) error {
	node, num := dao.GetTaskNum()
	var req getDataRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, getDataRespone{})
	}
	var resp getDataRespone
	cpu, mem, disk := dao.GetBarData(node)
	switch req.Data {

	case "taskNum":
		resp.NodeName = node
		resp.NodeValue = num
	case "diskSize":
		resp.NodeName = node
		resp.NodeValue = disk
	case "memorySize":
		resp.NodeName = node
		resp.NodeValue = mem
	case "cpuSize":
		resp.NodeName = node
		resp.NodeValue = cpu
	}
	fmt.Println(resp)
	return c.JSON(http.StatusOK, resp)

}
