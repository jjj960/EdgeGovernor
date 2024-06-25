package simple

//import (
//	"log"
//)
//
//type PreDeployMonitor struct {
//	taskList []string
//	preChan  chan string
//	done     chan struct{} // 新增一个用于通知 goroutine 退出的通道
//}
//
//func NewPreDeployMonitor(preChan chan string) *PreDeployMonitor {
//	return &PreDeployMonitor{
//		taskList: []string{},
//		preChan:  preChan,
//		done:     make(chan struct{}), // 初始化通知退出的通道
//	}
//}
//
//func (pdm *PreDeployMonitor) Start() {
//	// 启动一个 goroutine 监听任务触发通道
//	go func() {
//		for {
//			select {
//			case taskName := <-pdm.preChan:
//				log.Println("Task pre deploy:", taskName)
//				pdm.taskList = append(pdm.taskList, taskName)
//				Predeploy(pdm.taskList)
//				delete(tasks, taskName) // 执行完成后删除任务
//			case <-pdm.done: // 当 done 通道关闭时退出 goroutine
//				return
//			}
//		}
//	}()
//	// 主 goroutine 阻塞保持程序运行
//	<-pdm.done
//}
//
//func (pdm *PreDeployMonitor) Close() {
//	close(pdm.done) // 通知 goroutine 退出
//}
