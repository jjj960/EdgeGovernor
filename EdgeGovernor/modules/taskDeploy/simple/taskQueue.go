package simple

import (
	"EdgeGovernor/pkg/models"
	"EdgeGovernor/pkg/sec"
	"EdgeGovernor/pkg/utils"
	"container/heap"
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

// TaskQueue 是任务队列
type TaskQueue []models.Task

var (
	Pq    TaskQueue
	mutex sync.Mutex
)

func (Pq TaskQueue) Len() int { return len(Pq) }

func (Pq TaskQueue) Less(i, j int) bool {
	if Pq[i].PublishTime.Equal(Pq[j].PublishTime) {
		return Pq[i].Priority < Pq[j].Priority
	}
	return Pq[i].PublishTime.Before(Pq[j].PublishTime)
}

func (Pq TaskQueue) Swap(i, j int) { Pq[i], Pq[j] = Pq[j], Pq[i] }

func (Pq *TaskQueue) Push(x interface{}) {
	item := x.(models.Task)
	*Pq = append(*Pq, item)
}

func (Pq *TaskQueue) Pop() interface{} {
	old := *Pq
	n := len(old)
	item := old[n-1]
	*Pq = old[0 : n-1]
	return item
}

// peekFirstTask 获取第一个任务的信息但不弹出队列
func PeekFirstTask(Pq *TaskQueue) (models.Task, error) {
	mutex.Lock()
	defer mutex.Unlock()

	if Pq.Len() == 0 {
		return models.Task{}, errors.New("The task queue is empty")
	}

	return (*Pq)[0], nil
}

// 添加任务
func PushTask(task *models.Task) error {
	mutex.Lock()

	if task.PublishTime.After(time.Now()) {
		heap.Push(&Pq, task)
		mutex.Unlock()

		if Pq.Len() == 1 {
			SetTiming()
			utils.TaskMonitorChannel <- true
		}
	} else {
		return errors.New("The task's publication time has passed and cannot be added.")
	}

	return nil
}

func PopTask() (models.Task, error) {
	mutex.Lock()

	if Pq.Len() == 0 {
		return models.Task{}, errors.New("The task queue is empty")
	}
	task := heap.Pop(&Pq).(models.Task)

	mutex.Unlock()

	if Pq.Len() == 0 {
		utils.TaskMonitorChannel <- false
	} else {
		SetTiming()
	}

	return task, nil
}

// deleteTask 删除指定任务，保持队列顺序不变
func DeleteTask(Pq *TaskQueue, taskName string) error {
	mutex.Lock()
	defer mutex.Unlock()

	for i, task := range *Pq {
		if task.Name == taskName {
			*Pq = append((*Pq)[:i], (*Pq)[i+1:]...) // 从队列中删除指定任务
			if len(*Pq) == 0 {
				utils.TaskMonitorChannel <- false
			} else {
				SetTiming()
			}

			return nil
		}
	}

	return errors.New("Task not found in queue")
}

func MoveTaskToFront(Pq *TaskQueue, taskName string) error {
	mutex.Lock()
	defer mutex.Unlock()

	// 检查队列是否为空
	if len(*Pq) == 0 {
		return errors.New("Task queue is empty")
	}

	// 获取第一个任务的部署时间
	firstTask := (*Pq)[0]
	for _, task := range *Pq {
		if task.Name == taskName {
			// 检查部署时间是否在第一个任务之后
			if task.PublishTime.After(firstTask.PublishTime) {
				return errors.New("Task deployment time is after the first task")
			}

			// 移动任务到队列最前面
			taskIndex := 0
			for i, t := range *Pq {
				if t.Name == taskName {
					taskIndex = i
					break
				}
			}

			task := (*Pq)[taskIndex]
			copy((*Pq)[1:], (*Pq)[0:taskIndex])
			(*Pq)[0] = task

			SetTiming()

			return nil
		}
	}

	return errors.New("Task not found in queue")
}

func ModifyTask(taskName string, newParam *models.Task) error {
	err := DeleteTask(&Pq, taskName)
	if err != nil {
		return errors.New(err.Error())
	}

	err = PushTask(newParam)
	if err != nil {
		return errors.New(err.Error())
	}

	return nil
}

func ClearQueue() { //清除队列
	mutex.Lock()
	defer mutex.Unlock()

	Pq = make(TaskQueue, 0)
}

func GetTotalRequestResource() (int64, int64, int64, int64) {
	mutex.Lock()
	defer mutex.Unlock()

	var CPU, Mem, Net, Disk int64

	for _, task := range Pq {
		CPU = CPU + task.RequestCPU
		Mem = Mem + task.RequestMem
		Net = Net + task.RequestNet
		Disk = Disk + task.RequestDisk
	}

	return CPU, Mem, Net, Disk
}

func BackupTaskQueue() error {
	jsonData, err := utils.Jsoniter.Marshal(Pq)
	if err != nil {
		return fmt.Errorf("Error marshalling tasks to JSON: %s", err)
	}

	encryptQueue := sec.Safer.Encrypt(jsonData) //加密map

	_, err = utils.ETCDCli.Put(context.Background(), "/menet/backup/taskQueue", string(encryptQueue))
	if err != nil {
		return fmt.Errorf("failed to put queue in etcd: %w", err)
	}

	return nil
}

func GetBackupTaskQueue() error {
	data, err := utils.ETCDCli.Get(context.Background(), "/menet/backup/taskQueue")
	if err != nil {
		return fmt.Errorf("failed to get map in etcd: %s", err)
	}

	for _, kv := range data.Kvs {
		decryptMap := sec.Safer.Decrypt(kv.Value)

		var newPq TaskQueue
		err = utils.Jsoniter.Unmarshal([]byte(decryptMap), &newPq)
		if err != nil {
			return fmt.Errorf("Error unmarshalling TaskQueue data: %s", err)
		}
		Pq = newPq
	}

	return nil
}
