package main

import (
	"fmt"
	"os"

	"github.com/sudofrost/task-tracker-cli/cmd"
	"github.com/sudofrost/task-tracker-cli/tracker"
	"github.com/sudofrost/task-tracker-cli/persist"
)



func main() {
  cli := cmd.NewCLI(os.Args[1:])
  
	tasks, err := persist.Load[[]*tracker.Task]("tasks.json")
	if err != nil {
		panic(err)
	}
	tracker := tracker.Tracker{Tasks: tasks}
	defer func ()  {
		persist.Save("tasks.json", tracker.Tasks)
	}()

  cli.AddCommand(cmd.NewCommand("add", func(args []string) {
		if len(args) == 0 {
			fmt.Println("you must specify a task description")
			return
		}
		description := args[0]
		id := tracker.AddNewTask(description)
		fmt.Println(id)
	}))

	cli.Run()
}
