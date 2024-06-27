package routers

import (
	"EdgeGovernor/modules/web/app/controllers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type WebServer struct {
	echo *echo.Echo
}

func NewWebServer() *WebServer {
	e := echo.New()
	e.Use(middleware.CORS())

	webServer := &WebServer{
		echo: e,
	}

	webServer.registerRoutes()

	return webServer
}

func (ws *WebServer) registerRoutes() {
	e := ws.echo

	e.POST("/Login", controllers.Login)
	e.GET("/check-connection", controllers.CheckConnection)

	e.GET("/getnodeName", controllers.GetNodeName)
	e.POST("/getnodeMessage", controllers.GetNodeMsg)
	e.POST("/searchNode", controllers.SearchNode)

	/*单一任务相关操作*/
	e.POST("/addTask", controllers.AddTask)
	e.POST("/searchTask", controllers.SearchTask)
	e.POST("/deleteTask", controllers.DeleteTask)
	e.POST("/getTaskNum", controllers.GetTaskNum)
	e.POST("/postTransfer", controllers.MigrationTask)

	/*工作流任务相关操作*/
	e.POST("/addWorkflow", controllers.AddWorkflow)
	e.POST("/searchWorkflow", controllers.SearchWorkflow)
	e.POST("/deleteWorkflow", controllers.DeleteWorkflow)
	e.POST("/getWorkflowNum", controllers.GetWorkflowNum)

	/*算法相关操作*/
	e.POST("/addAlgorithm", controllers.AddAlgorithm)
	e.GET("/getAlgorithmName", controllers.GetAlgorithmName)
	e.GET("/getAlgorithm", controllers.GetAlgorithm)
	e.POST("/searchAlgorithm", controllers.SearchAlgorithm)

	/*集群负载相关操作*/
	e.POST("/getbarData", controllers.GetBarData)

	/*文件操作*/
	e.POST("/upload", controllers.FileUpLoad)
	e.POST("/removeFile", controllers.RemoveFile)
	e.POST("/workLoadPublish", controllers.WorkLoadPublish)

	/*日志操作*/
	e.GET("/getTaskLog", controllers.GetTaskLog)
	e.POST("/deleteTaskLog", controllers.DeleteTaskLog)
	e.GET("/deleteAllLog", controllers.DeleteAllLog)
	e.GET("/SSE_Connect/:client_id", controllers.SSE_Connect)
}

func (ws *WebServer) Start(port string) error {
	return ws.echo.Start(":" + port)
}

func (ws *WebServer) Close() error {
	return ws.echo.Close()
}
