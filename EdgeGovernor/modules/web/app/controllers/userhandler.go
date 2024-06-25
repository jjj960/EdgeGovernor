package controllers

import (
	"github.com/labstack/echo/v4"
	"net/http"
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
		return c.JSON(http.StatusOK, loginResponse{Status: "success"})
	} else {
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
