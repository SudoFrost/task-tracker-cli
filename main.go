package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/sudofrost/task-tracker-cli/cmd"
	"github.com/sudofrost/task-tracker-cli/persist"
	"github.com/sudofrost/task-tracker-cli/tracker"
)

func PrintAsTable(columns []string, rows []map[string]string, seprator string) {
	columnsLengths := make([]int, len(columns))

	for i, column := range columns {
		columnsLengths[i] = len(column)
	}

	for _, row := range rows {
		for i, column := range columns {
			if len(row[column]) > columnsLengths[i] {
				columnsLengths[i] = len(row[column])
			}
		}
	}

	for i, column := range columns {
		fmt.Printf("%s%s", strings.Repeat(" ", columnsLengths[i]-len(column)), column)
		if i < len(columns)-1 {
			fmt.Print(seprator)
		}
	}
	fmt.Println()

	for _, row := range rows {
		for i, column := range columns {
			fmt.Printf("%s%s", strings.Repeat(" ", columnsLengths[i]-len(row[column])), row[column])
			if i < len(columns)-1 {
				fmt.Print(seprator)
			}
		}
		fmt.Println()
	}
}

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

		columns := []string{"ID", "Created AT", "Updated AT", "Status", "Description"}
		data := make([]map[string]string, 0)

		for _, task := range tasks {
			data = append(data, map[string]string{
				"ID":          fmt.Sprintf("%d", task.ID),
				"Created AT":  time.UnixMilli(task.CreatedAt).Format(time.DateTime),
				"Updated AT":  time.UnixMilli(task.UpdatedAt).Format(time.DateTime),
				"Status":      string(task.Status),
				"Description": task.Description,
			})
		}

		PrintAsTable(columns, data, " | ")
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
