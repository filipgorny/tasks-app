package main

import (
	"github.com/filipgorny/org-tool/app"
	"github.com/filipgorny/org-tool/console"
)

func main() {
	c := console.NewConsole()

	tasksScreen := app.NewTaskScreen()

	c.RegisterScreen("tasks", &tasksScreen)

	c.LoadScreen("tasks")
	c.Run()
}
