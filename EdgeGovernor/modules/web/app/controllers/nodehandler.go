package controllers

import (
	"EdgeGovernor/modules/web/app/dao"
	"EdgeGovernor/pkg/constants"
	"EdgeGovernor/pkg/utils"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func GetNodeName(c echo.Context) error {

	nodeName := utils.NodeTables.GetAllNodeName()

	response := map[string]interface{}{
		"node": nodeName,
	}

	fmt.Println(response)
	return c.JSON(http.StatusOK, response)
}

func GetNodeMsg(c echo.Context) error {
	currentPage, _ := strconv.Atoi(c.FormValue("currentPage"))
	pageSize, _ := strconv.Atoi(c.FormValue("pageSize"))
	fmt.Println(currentPage, pageSize)
	start := (currentPage - 1) * pageSize
	end := currentPage * pageSize
	if end > constants.NodeCount {
		end = constants.NodeCount
	}

	nodeMessages := dao.GetNodesMsg(start, end)

	response := map[string]interface{}{
		"nodeMessage": nodeMessages,
		"total":       constants.NodeCount,
	}
	fmt.Println(response)
	return c.JSON(http.StatusOK, response)
}

func SearchNode(c echo.Context) error {
	nodeName := c.FormValue("nodeName")
	fmt.Println(nodeName)
	nodeMessages := dao.SearchNodeMsg(nodeName)

	response := map[string]interface{}{
		"nodeMessage": nodeMessages,
	}
	fmt.Println(response)
	return c.JSON(http.StatusOK, response)
}
