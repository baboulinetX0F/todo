package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
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
	tags        []string
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

			// If task is marked as done (first part = X)
			if strings.Compare(parts[0], "X") == 0 {
				task.status = true
				for index := 1; index < len(parts); index++ {
					if index > 1 {
						task.description += " "
					}
					task.description += parts[index]
				}
			} else {
				task.status = false
				for index := 0; index < len(parts); index++ {
					// if part is a tag
					if len(parts[index]) > 0 && parts[index][0] == '+' {
						task.tags = append(task.tags, parts[index])
					}
					if index > 0 {
						task.description += " "
					}
					task.description += parts[index]
				}
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

// SaveTasks : Write the pFilePath file with the contents of the Tasks Array
func SaveTasks(pTasks []Task, pFilePath string) {
	file, err := os.OpenFile("todo.txt", os.O_WRONLY, 0666)
	file.Truncate(0)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, task := range pTasks {
		if task.status {
			writer.WriteString("X ")
		}
		writer.WriteString(task.description + "\n")
	}
	writer.Flush()
}

// AddTask : Add the task passed in parameter to the todo file
func AddTask(pFilePath string, pLineToParse string) {
	file, err := os.OpenFile("todo.txt", os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// TODO: Add parsing of line passed in parameter
	writer := bufio.NewWriter(file)
	writer.WriteString(pLineToParse + "\n")
	writer.Flush()
}

// SetTaskStatus : Mark as done the task with the id passed in paramater
func SetTaskStatus(pID uint16, pNewState bool) {
	tasks := LoadTasks("todo.txt")
	if pID > uint16(len(tasks)) || pID <= 0 {
		log.Println("ValidateTask : index out of range")
	} else {
		tasks[pID-1].status = pNewState
	}

	SaveTasks(tasks, "todo.txt")
}

// ArchiveTasks : remove all tasks done and store them in the archive file
func ArchiveTasks() {
	tasks := LoadTasks("todo.txt")

	// Open / Create archive file
	file, err := os.OpenFile("archive.txt", os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for i := len(tasks) - 1; i >= 0; i-- {
		if tasks[i].status {
			writer.WriteString("X " + tasks[i].description + "\n")
			tasks = append(tasks[:i], tasks[i+1:]...)
		}
	}
	writer.Flush()
	SaveTasks(tasks, "todo.txt")
}

// ListTasks : Display a list of all the tasks
func ListTasks() {
	fmt.Println("Task List :")
	tasks := LoadTasks("todo.txt")
	for _, task := range tasks {
		PrintTask(task)
	}
}

func ListTasksFiltered(searchTags []string) {
	fmt.Println("Task List :")
	tasks := LoadTasks("todo.txt")
	for _, task := range tasks {
		displayTask := true
		for filterIdx := range searchTags {
			tagFound := false
			for i := range task.tags {
				if task.tags[i] == searchTags[filterIdx] {
					tagFound = true
					break
				}
			}
			if !tagFound {
				displayTask = false
				break
			}
		}
		if displayTask {
			PrintTask(task)
		}
	}
}

func main() {
	args := os.Args[1:]
	if len(args) > 0 {
		if args[0] == "ls" {
			if len(args) > 1 {
				ListTasksFiltered(args[1:])
			} else {
				ListTasks()
			}
		} else if args[0] == "archive" {
			ArchiveTasks()
		} else if args[0] == "add" && len(args) > 1 {
			AddTask("todo.txt", args[1])
		} else if (args[0] == "do" || args[0] == "undo") && len(args) > 1 {
			id, casterr := strconv.ParseUint(args[1], 10, 16)
			if casterr != nil {
				log.Fatal(casterr)
			} else {
				if args[0] == "do" {
					SetTaskStatus(uint16(id), true)
				} else {
					SetTaskStatus(uint16(id), false)
				}
			}
		}
	}
}
