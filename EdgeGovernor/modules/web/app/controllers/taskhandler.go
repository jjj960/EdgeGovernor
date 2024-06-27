package controllers

import (
	"EdgeGovernor/modules/taskDeploy/simple"
	"EdgeGovernor/modules/web/app/dao"
	"EdgeGovernor/pkg/constants"
	"EdgeGovernor/pkg/database/duckdb"
	"EdgeGovernor/pkg/k8s"
	"EdgeGovernor/pkg/models"
	"EdgeGovernor/pkg/utils"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"strconv"
	"time"
)

func AddTask(c echo.Context) error {
	task := new(models.Task)
	if err := c.Bind(task); err != nil {
		return err
	}
	fmt.Println(task)
	fmt.Println(task.PublishTime)

	err := taskProcess(task)
	if err != nil {
		return err
	}
	// 返回成功响应
	return c.JSON(http.StatusOK, map[string]string{"message": "Task added successfully"})
}

func taskProcess(task *models.Task) error {
	if task.PublishTime.IsZero() {
		task.PublishTime = time.Now()
		err := simple.AddTask(task, false)
		if err != nil {
			return err
		}
	} else {
		isPassed, err := isExecutionTimePassed(task.PublishTime)
		if err != nil {
			log.Printf("Error parsing execution time: %v", err)
			return errors.New("Error parsing execution time")
		}

		if isPassed {
			err = simple.AddTask(task, true)
			if err != nil {
				return err
			}
		} else {
			log.Printf("Refuse adding task: %s, execution time: %s has passed", task.Name, task.PublishTime)
			return errors.New("Refuse adding task, execution time has passed")
		}
	}

	return nil
}

func isExecutionTimePassed(executeTime time.Time) (bool, error) {
	location, err := time.LoadLocation("Asia/Shanghai") // 指定时区为亚洲/上海
	if err != nil {
		fmt.Printf("Failed to load time zone: %v", err)
		return false, err
	}

	currentTime := time.Now().In(location) // 获取当前时间并设置时区为上海

	return executeTime.After(currentTime), nil
}

func SearchTask(c echo.Context) error {
	nodeName := c.FormValue("nodeName")
	tasks := dao.GetNodeTasks(nodeName)
	currentPage, _ := strconv.Atoi(c.FormValue("currentPage"))
	pageSize, _ := strconv.Atoi(c.FormValue("pageSize"))

	start := (currentPage - 1) * pageSize
	end := currentPage * pageSize
	if end > len(tasks) {
		end = len(tasks)
	}

	taskMessage := tasks[start:end]

	response := map[string]interface{}{
		"taskMessage": taskMessage,
		"total":       len(tasks),
	}

	return c.JSON(http.StatusOK, response)
}

func DeleteTask(c echo.Context) error {
	taskName := c.FormValue("taskName")
	fmt.Println(taskName)
	dao.DeleteTask(taskName)
	return c.String(http.StatusOK, "成功")
}

func MigrationTask(c echo.Context) error {
	taskName := c.FormValue("taskName")
	sourceNode := c.FormValue("From")
	targetNode := c.FormValue("To")

	if dao.CheckTaskNode(taskName, sourceNode) { //检查该节点是否存在这个任务
		dao.UpdateTaskMsg(taskName, targetNode)
		_, err := k8s.UpdatePodNode(taskName, targetNode)
		if err != nil {
			id, _ := utils.GetID()
			logEntry := models.OperationLog{
				ID:            id,
				NodeName:      constants.Hostname,
				NodeIP:        constants.IP,
				OperationType: "task migration",
				Description:   fmt.Sprintf("The user with IP address %s is attempting to migrate task %s, and the result is %s", c.RealIP(), taskName, "fail"),
				Result:        true,
				CreatedAt:     time.Now(),
			}
			duckdb.InsertOperationLog(logEntry)
			return err
		}
	} else {
		log.Println("The task")
		return errors.New("The task does not meet the migration criteria.")
	}

	id, _ := utils.GetID()
	logEntry := models.OperationLog{
		ID:            id,
		NodeName:      constants.Hostname,
		NodeIP:        constants.IP,
		OperationType: "task migration",
		Description:   fmt.Sprintf("The user with IP address %s is attempting to migrate task %s, and the result is %s", c.RealIP(), taskName, "success"),
		Result:        true,
		CreatedAt:     time.Now(),
	}
	duckdb.InsertOperationLog(logEntry)

	return c.String(http.StatusOK, "成功")
}

func GetTaskNum(c echo.Context) error {
	nodeName, nodeValue := dao.GetTaskNum()

	response := map[string]interface{}{
		"nodeName":  nodeName,
		"nodeValue": nodeValue,
	}
	fmt.Println(response)
	return c.JSON(http.StatusOK, response)
}
