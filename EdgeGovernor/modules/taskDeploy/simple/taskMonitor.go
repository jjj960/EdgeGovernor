package simple

import (
	"EdgeGovernor/pkg/models"
	"EdgeGovernor/pkg/utils"
	"fmt"
	"time"
)

type TaskModule struct {
	taskMonitorChannel   chan bool
	taskOperationChannel chan bool
	done                 chan struct{} // 用于通知 goroutine 退出的通道
}

func NewTaskModule() *TaskModule {
	return &TaskModule{
		taskMonitorChannel:   utils.TaskMonitorChannel,
		taskOperationChannel: utils.PopTaskChannel,
		done:                 make(chan struct{}), // 初始化通知退出的通道
	}
}

var count int64
var second int64
var counting bool

func (tm *TaskModule) Start() {

	go func() {
		for {
			select {
			case start := <-utils.TaskMonitorChannel:
				if start {
					if !counting {
						counting = true
					}
				} else {
					if counting {
						count = 0
						counting = false
						fmt.Println("There are no more tasks in the task queue")
					}
				}
			case <-time.After(time.Second):

			default:
				if counting {
					count++
					if count == second {
						utils.PopTaskChannel <- true
					}
				}
				time.Sleep(time.Second)
			}
		}
	}()

	go func() {
		for {
			select {
			case signal := <-utils.PopTaskChannel:
				if signal {
					task1, err := PopTask()
					if err == nil {
						err = taskPublish(&task1)
						id := utils.SnowFlake.Generate()
						utils.AlarmMsgChannel <- models.Msg{
							ID:           id.Int64(),
							GenerateTime: id.Time(),
							Tpye:         "Task deployment failed",
							Detail:       []byte(err.Error()),
							Status:       false,
						}
					}
					for {
						nextTask, err := PeekFirstTask(&Pq)
						if err != nil {
							break
						}
						if nextTask.PublishTime == task1.PublishTime {
							poppedTask, _ := PopTask()
							err = taskPublish(&poppedTask)
							id := utils.SnowFlake.Generate()
							utils.AlarmMsgChannel <- models.Msg{
								ID:           id.Int64(),
								GenerateTime: id.Time(),
								Tpye:         "Task deployment failed",
								Detail:       []byte(err.Error()),
								Status:       false,
							}
						} else {
							break
						}
					}
				}
			}
		}
	}()

	<-tm.done
}

func (tm *TaskModule) Close() {
	close(tm.done)
}

func SetTiming() {
	nextTask, _ := PeekFirstTask(&Pq)
	fmt.Println(nextTask)
	duration := nextTask.PublishTime.Sub(time.Now())
	second = int64(duration.Seconds())
	count = 0
}
