package main

import (
	"fmt"
	"os"

	"github.com/sudofrost/task-tracker-cli/cmd"
	"github.com/sudofrost/task-tracker-cli/tracker"
)

func main() {
  cli := cmd.NewCLI(os.Args[1:])

	tracker := &tracker.Tracker{}

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
