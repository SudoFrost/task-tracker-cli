package tracker

import (
	"fmt"
	"time"
)

type TaskStatus string

const (
	TaskStatusTodo       TaskStatus = "todo"
	TaskStatusInProgress TaskStatus = "in-progress"
	TaskStatusDone       TaskStatus = "done"
)

type Task struct {
	ID          uint64     `json:"id"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"status"`
	CreatedAt   int64      `json:"createdAt"`
	UpdatedAt   int64      `json:"updatedAt"`
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

func (t *Tracker) AddNewTask(description string) uint64 {
	now := time.Now().UTC().UnixMilli()
	task := &Task{
		ID:          t.NewId(),
		Description: description,
		Status:      TaskStatusTodo,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	t.Tasks = append(t.Tasks, task)
	return task.ID
}

func (t *Tracker) GetTask(id uint64) (*Task, error) {
	for _, task := range t.Tasks {
		if task.ID == id {
			return task, nil
		}
	}
	return nil, fmt.Errorf("task with id %d not found", id)
}

func (t *Tracker) UpdateStatus(id uint64, status TaskStatus) error {
	task, err := t.GetTask(id)
	if err != nil {
		return err
	}
	task.Status = status
	task.UpdatedAt = time.Now().UTC().UnixMilli()
	return nil
}

func (t *Tracker) UpdateDescription(id uint64, description string) error {
	task, err := t.GetTask(id)
	if err != nil {
		return err
	}
	task.Description = description
	task.UpdatedAt = time.Now().UTC().UnixMilli()
	return nil
}

func (t *Tracker) DeleteTask(id uint64) error {
	for i, task := range t.Tasks {
		if task.ID == id {
			t.Tasks = append(t.Tasks[:i], t.Tasks[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("task with id %d not found", id)
}

func (t *Tracker) GetTasks(status *TaskStatus) []*Task {
	tasks := make([]*Task, 0)

	for _, task := range t.Tasks {
		if status == nil || *status == task.Status {
			tasks = append(tasks, task)
		}
	}

	return tasks
}
