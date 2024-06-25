package controllers

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func FileUpLoad(c echo.Context) error {
	file, err := c.FormFile("upload")
	if err != nil {
		return c.String(http.StatusBadRequest, "No upload selected for uploading")
	}
	// 先打开文件源
	src, err := file.Open()
	if err != nil {
		return c.String(http.StatusInternalServerError, "open upload source error")
	}
	defer src.Close()

	// 保存文件的目录
	dir, _ := os.Getwd()
	dir = filepath.Dir(dir)
	savePath := filepath.Join(dir, "upload")

	// 确保保存目录存在
	if _, err := os.Stat(savePath); os.IsNotExist(err) {
		err = os.Mkdir(savePath, 0777)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to create upload directory")
		}
	}

	// 文件名
	filename := filepath.Base(file.Filename)
	// 完整路径
	file_path := filepath.Join(savePath, filename)

	// 打开文件以保存上传的文件
	out, err := os.Create(file_path)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to create upload")
	}
	defer out.Close()

	//复制文件内容
	_, err = io.Copy(out, src)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to copy upload")
	}

	// 返回一个JSON响应
	return c.JSON(http.StatusOK, map[string]interface{}{
		"url":  file_path, // 这里需要根据实际情况替换为你的服务器URL
		"name": filename,
	})
}

func RemoveFile(c echo.Context) error {
	fileName := c.FormValue("fileName")

	if fileName == "" {
		return c.String(http.StatusBadRequest, "文件名不能为空")
	}

	// 保存文件的目录
	dir, _ := os.Getwd()
	dir = filepath.Dir(dir)
	savePath := filepath.Join(dir, "upload")

	// 完整路径
	filePath := filepath.Join(savePath, fileName)

	err := os.Remove(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return c.String(http.StatusNotFound, fmt.Sprintf("文件 %s 不存在", filePath))
		}
		if os.IsPermission(err) {
			return c.String(http.StatusForbidden, fmt.Sprintf("没有足够的权限删除文件 %s", filePath))
		}
		return c.String(http.StatusInternalServerError, fmt.Sprintf("删除文件 %s 时发生错误: %v", filePath, err))
	}

	return c.String(http.StatusOK, "文件删除成功")
}

func WorkLoadPublish(c echo.Context) error {
	taskName := c.FormValue("taskName")
	taskNum := c.FormValue("taskNum")
	tpValue := c.FormValue("tpValue")
	taskpostTime1 := c.FormValue("taskpostTime1")
	taskpostTime2 := c.FormValue("taskpostTime2")
	pmValue := c.FormValue("pmValue")
	textArea := c.FormValue("textArea")
	fileList := c.FormValue("fileList")
	date := c.FormValue("date")
	fmt.Println(taskName, taskNum, tpValue, taskpostTime1, taskpostTime2, pmValue, textArea, fileList, date)
	return c.String(http.StatusOK, "发布成功")
}
