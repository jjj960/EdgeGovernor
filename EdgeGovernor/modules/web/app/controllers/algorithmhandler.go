package controllers

import (
	"EdgeGovernor/modules/web/app/dao"
	"github.com/labstack/echo/v4"
	"net/http"
)

func AddAlgorithm(c echo.Context) error {
	name := c.FormValue("name")
	url := c.FormValue("url")
	use := c.FormValue("use")
	detail := c.FormValue("detail")

	dao.AddAlgorithms(name, url, use, detail)
	return c.JSON(http.StatusOK, map[string]string{"status": "success"})
}

func GetAlgorithm(c echo.Context) error {
	result := dao.GetAlgorithms()
	response := map[string]interface{}{
		"algorithms": result,
		"pageSizes":  len(result),
	}
	return c.JSON(http.StatusOK, response)
}

func SearchAlgorithm(c echo.Context) error {
	alName := c.FormValue("algorithmName")
	result := dao.GetAlgorithmMsg(alName)

	response := map[string]map[string]string{
		"algorithmMessage": {
			"name":   result.Name,
			"mirror": result.Status,
			"URL":    result.URL,
			"detail": result.Type,
		},
	}

	return c.JSON(http.StatusOK, response)
}

func GetAlgorithmName(c echo.Context) error {
	names := dao.GetAlgorithmNames()

	response := map[string]interface{}{
		"algorithms": names,
	}

	return c.JSON(http.StatusOK, response)
}
