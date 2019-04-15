package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

const (
	RedColor   = "\033[1;31m%s\033[0m"
	GreenColor = "\033[1;32m%s\033[0m"
)

// Task : Data structure containing all the informations from a Task
type Task struct {
	id          uint16
	description string
	status      bool
}

// PrintTask : Print to the right format the task passed in parameter
func PrintTask(pTask Task) {
	if pTask.status == false {
		fmt.Printf(RedColor, "☐ "+fmt.Sprint(pTask.id)+". "+pTask.description)
	} else {
		fmt.Printf(GreenColor, "☑ "+fmt.Sprint(pTask.id)+". "+pTask.description)
	}
	fmt.Printf("\n")

}

func LoadTasks(pFilePath string) []Task {
	var tasks []Task
	var index uint16
	file, err := os.Open(pFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	index = 1
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		task := Task{id: index, description: scanner.Text()}
		tasks = append(tasks, task)
		index++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return tasks
}

func main() {
	tasks := LoadTasks("test.txt")
	for _, task := range tasks {
		PrintTask(task)
	}
}
