package controllers

//
//import (
//	"EdgeGovernor/modules/web/app/dao"
//	"EdgeGovernor/pkg/constants"
//	"github.com/labstack/echo/v4"
//	"net/http"
//	"strconv"
//)
//
//func GetMirrorMsg(c echo.Context) error {
//
//	currentPage, _ := strconv.Atoi(c.FormValue("currentPage"))
//	pageSize, _ := strconv.Atoi(c.FormValue("pageSize"))
//
//	start := (currentPage - 1) * pageSize
//	end := currentPage * pageSize
//	if end > constants.MirrorCount {
//		end = constants.MirrorCount
//	}
//
//	mirrorMessages := dao.GetMirrorsMsg(start, end)
//
//	response := map[string]interface{}{
//		"mirrorMessages": mirrorMessages,
//		"total":          constants.MirrorCount,
//	}
//
//	return c.JSON(http.StatusOK, response)
//}
//
//func GetMirrorsName(c echo.Context) error {
//	mirrorNames := dao.GetMirrorsName()
//	response := map[string]interface{}{
//		"mirrorNames": mirrorNames,
//	}
//
//	return c.JSON(http.StatusOK, response)
//}
//
//func AddMirror(c echo.Context) error {
//	mirrorName := c.FormValue("mirrorRes")
//	ip := c.FormValue("mirrorIP")
//	port := c.FormValue("mirrorPort")
//	username := c.FormValue("loginName")
//	password := c.FormValue("password")
//	status := c.FormValue("use")
//	detail := c.FormValue("remark")
//	if dao.CheckMirrorNameUniqueness(mirrorName) {
//		return c.String(http.StatusOK, "fail")
//	} else {
//		dao.AddMirror(mirrorName, ip, port, username, password, status, detail)
//		return c.String(http.StatusOK, "success")
//	}
//}
//
//func ChangStatus(c echo.Context) error {
//	mirrorName := c.FormValue("mirrorRes")
//	dao.ChangeMirrorStatus(mirrorName)
//	return c.String(http.StatusOK, "成功")
//}
//
//func DeleteMirror(c echo.Context) error {
//	mirrorName := c.FormValue("mirrorRes")
//	if dao.CheckMirrorNameisExist(mirrorName) {
//		dao.DeleteMirror(mirrorName)
//		return c.String(http.StatusOK, "成功")
//	} else {
//		return c.String(http.StatusOK, "失败")
//	}
//}
//
//func SearchMirror(c echo.Context) error {
//	mirrorName := c.FormValue("mirrorRes")
//
//	mirror := dao.SearchMirror(mirrorName)
//
//	response := map[string]interface{}{
//		"mirrorMessage": mirror,
//	}
//
//	return c.JSON(http.StatusOK, response)
//}
