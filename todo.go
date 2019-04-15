package main

import (
	"fmt"
)

const (
	RedColor = "\033[1;31m%s\033[0m"
)

// Task : Data structure containing all the informations from a Task
type Task struct {
	id          uint16
	description string
	status      bool
}

// PrintTask : Print to the right format the task passed in parameter
func PrintTask(pTask Task) {
	fmt.Printf(RedColor, fmt.Sprint(pTask.id)+". "+pTask.description)
}

func main() {
	t := Task{id: 1, description: "Je suis une tache.", status: false}
	PrintTask(t)
}
