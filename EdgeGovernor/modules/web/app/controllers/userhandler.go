package controllers

import (
	"EdgeGovernor/pkg/constants"
	"EdgeGovernor/pkg/database/duckdb"
	"EdgeGovernor/pkg/models"
	"EdgeGovernor/pkg/utils"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

type loginRequest struct {
	Password string `json:"password"`
}

type loginResponse struct {
	Status string `json:"status"`
}

func Login(c echo.Context) error {
	var req loginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, loginResponse{Status: "fail"})
	}

	if req.Password == "123456" {
		id, _ := utils.GetID()
		logEntry := models.OperationLog{
			ID:            id,
			NodeName:      constants.Hostname,
			NodeIP:        constants.IP,
			OperationType: "user login",
			Description:   fmt.Sprintf("User with IP address %s is attempting to connect to EdgeGovernor, and the result is %s", c.RealIP(), "success"),
			Result:        true,
			CreatedAt:     time.Now(),
		}
		duckdb.InsertOperationLog(logEntry)
		return c.JSON(http.StatusOK, loginResponse{Status: "success"})
	} else {
		id, _ := utils.GetID()
		logEntry := models.OperationLog{
			ID:            id,
			NodeName:      constants.Hostname,
			NodeIP:        constants.IP,
			OperationType: "user login",
			Description:   fmt.Sprintf("User with IP address %s is attempting to connect to EdgeGovernor, and the result is %s", c.RealIP(), "fail"),
			Result:        true,
			CreatedAt:     time.Now(),
		}
		duckdb.InsertOperationLog(logEntry)
		return c.JSON(http.StatusServiceUnavailable, loginResponse{Status: "fail"})
	}
}

type checkConnectionResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func CheckConnection(c echo.Context) error {
	var response checkConnectionResponse
	response = checkConnectionResponse{Status: "connected", Message: "Connection is normal"}
	return c.JSON(http.StatusOK, response)
}
