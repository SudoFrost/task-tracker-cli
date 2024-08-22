package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/sudofrost/task-tracker-cli/cmd"
	"github.com/sudofrost/task-tracker-cli/persist"
	"github.com/sudofrost/task-tracker-cli/tracker"
)

func main() {
	cli := cmd.NewCLI(os.Args[1:])

	tasks, err := persist.Load[[]*tracker.Task]("tasks.json")
	if err != nil {
		panic(err)
	}
	t := tracker.Tracker{Tasks: tasks}
	defer func() {
		persist.Save("tasks.json", t.Tasks)
	}()

	cli.AddCommand(cmd.NewCommand("add", func(args []string) {
		if len(args) == 0 {
			fmt.Println("you must specify a task description")
			return
		}
		description := args[0]
		id := t.AddNewTask(description)
		fmt.Println(id)
	}))

	cli.AddCommand(cmd.NewCommand("list", func(args []string) {
		var status tracker.TaskStatus
		if len(args) > 0 {
			switch args[0] {
			case string(tracker.TaskStatusTodo):
				status = tracker.TaskStatusTodo
			case string(tracker.TaskStatusInProgress):
				status = tracker.TaskStatusInProgress
			case string(tracker.TaskStatusDone):
				status = tracker.TaskStatusDone
			default:
				fmt.Printf("invalid status: %s\n", args[0])
				return
			}
		}
		var tasks []*tracker.Task
		if status == "" {
			tasks = t.GetTasks(nil)
		} else {
			tasks = t.GetTasks(&status)
		}

		if len(tasks) == 0 {
			fmt.Println("no tasks found")
			return
		}

		var maxID uint64

		for _, task := range tasks {
			if task.ID > maxID {
				maxID = task.ID
			}
		}

		lenghOfIDColumn := len(fmt.Sprintf("%d", maxID))
		lenghOfStatusColumn := len(status)
		if lenghOfStatusColumn == 0 {
			lenghOfStatusColumn = len(string(tracker.TaskStatusInProgress))
		}

		if lenghOfIDColumn < 2 {
			lenghOfIDColumn = 2
		}

		if lenghOfStatusColumn < 6 {
			lenghOfStatusColumn = 6
		}

		fmt.Printf(" %*s | %*s | %*s | %*s | %s\n", lenghOfIDColumn, "ID", 19, "Created AT", 19, "Updated AT", lenghOfStatusColumn, "Status", "Description")
		for _, task := range tasks {
			fmt.Printf(
				" %*d | %*s | %*s | %*s | %s\n",
				lenghOfIDColumn, task.ID,
				19, time.UnixMilli(task.CreatedAt).Format(time.DateTime),
				19, time.UnixMilli(task.UpdatedAt).Format(time.DateTime),
				lenghOfStatusColumn, task.Status,
				task.Description,
			)
		}
	}))

	cli.AddCommand(cmd.NewCommand("update", func(args []string) {
		if len(args) < 2 {
			fmt.Println("you must specify a task id and a new description")
			return
		}

		id, err := strconv.ParseUint(args[0], 10, 64)
		if err != nil {
			fmt.Println("invalid task id")
			return
		}

		err = t.UpdateDescription(id, args[1])
		if err != nil {
			fmt.Println(err)
			return
		}
	}))

	cli.AddCommand(cmd.NewCommand("mark-in-progress", func(args []string) {
		if len(args) < 1 {
			fmt.Println("you must specify a task id")
			return
		}

		id, err := strconv.ParseUint(args[0], 10, 64)
		if err != nil {
			fmt.Println("invalid task id")
			return
		}

		err = t.UpdateStatus(id, tracker.TaskStatusInProgress)
		if err != nil {
			fmt.Println(err)
			return
		}
	}))

	cli.AddCommand(cmd.NewCommand("mark-done", func(args []string) {
		if len(args) < 1 {
			fmt.Println("you must specify a task id")
			return
		}

		id, err := strconv.ParseUint(args[0], 10, 64)
		if err != nil {
			fmt.Println("invalid task id")
			return
		}

		err = t.UpdateStatus(id, tracker.TaskStatusDone)
		if err != nil {
			fmt.Println(err)
			return
		}
	}))

	cli.AddCommand(cmd.NewCommand("delete", func(args []string) {
		if len(args) < 1 {
			fmt.Println("you must specify a task id")
			return
		}

		id, err := strconv.ParseUint(args[0], 10, 64)
		if err != nil {
			fmt.Println("invalid task id")
			return
		}

		err = t.DeleteTask(id)
		if err != nil {
			fmt.Println(err)
			return
		}
	}))

	cli.Run()
}
