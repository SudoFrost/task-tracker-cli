package tracker

import "time"

type TaskStatus string

const (
	TaskStatusTodo TaskStatus = "todo"
	TaskStatusInProgress TaskStatus = "in-progress"
	TaskStatusDone TaskStatus = "done"
)

type Task struct {
	ID uint64
	Description string
	Status TaskStatus
	CreatedAt int64
	UpdatedAt int64
}

type Tracker struct {
	Tasks []*Task
}

func (T *Tracker) NewId() uint64 {
 var id uint64 = 0

 for _, task := range T.Tasks {
	 if id < task.ID {
		 id = task.ID
	 }
 }

 return id + 1
}

func (t *Tracker) AddNewTask(description string) uint64{
	now := time.Now().UTC().UnixMilli()
	task := &Task{
		ID: t.NewId(),
		Description: description,
		Status: TaskStatusTodo,
		CreatedAt: now,
		UpdatedAt: now,
	}
	t.Tasks = append(t.Tasks, task)
	return task.ID
}
