package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

// ANSI Strings for colored output
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

// LoadTasks : Return an array of Tasks from the file passed in parameter
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
		parts := strings.Split(scanner.Text(), " ")
		if len(parts) == 1 {
			task := Task{id: index, description: parts[0], status: false}
			tasks = append(tasks, task)
		} else if len(parts) > 1 {
			task := Task{id: index}
			if strings.Compare(parts[0], "X") == 0 {
				task.status = true
			} else {
				task.status = false
			}
			for index := 1; index < len(parts); index++ {
				task.description += " " + parts[index]
			}
			tasks = append(tasks, task)
		}
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
