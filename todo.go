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

var todoDirPath string
var todoFilePath string
var todoArchivePath string

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
					// if part is a tag
					if len(parts[index]) > 0 && parts[index][0] == '+' {
						task.tags = append(task.tags, parts[index])
					}
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
	file, err := os.OpenFile(todoFilePath, os.O_WRONLY|os.O_CREATE, 0666)
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
	file, err := os.OpenFile(todoFilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// TODO: Add parsing of line passed in parameter
	writer := bufio.NewWriter(file)
	writer.WriteString(pLineToParse + "\n")
	writer.Flush()
}

// RemoveTask : Remove a task from the list
func RemoveTask(pTaskID uint16) {
	tasks := LoadTasks(todoFilePath)
	tasks = append(tasks[:pTaskID-1], tasks[pTaskID:]...)
	SaveTasks(tasks, todoFilePath)
}

// SetTaskStatus : Mark as done the task with the id passed in paramater
func SetTaskStatus(pID uint16, pNewState bool) {
	tasks := LoadTasks(todoFilePath)
	if pID > uint16(len(tasks)) || pID <= 0 {
		log.Println("ValidateTask : index out of range")
	} else {
		tasks[pID-1].status = pNewState
	}

	SaveTasks(tasks, todoFilePath)
}

// ArchiveTasks : remove all tasks done and store them in the archive file
func ArchiveTasks() {
	tasks := LoadTasks(todoFilePath)

	// Open / Create archive file
	file, err := os.OpenFile(todoArchivePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
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
	SaveTasks(tasks, todoFilePath)
}

// ListTasks : Display a list of all the tasks
func ListTasks() {
	fmt.Println("Task List :")
	tasks := LoadTasks(todoFilePath)
	for _, task := range tasks {
		PrintTask(task)
	}
}

// ListTasksFiltered : Display all task filtered by any tag(s) given
func ListTasksFiltered(searchTags []string) {
	fmt.Println("Task List :")
	tasks := LoadTasks(todoFilePath)
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

// Init : Initialize configuration / directory and files used by the software
func Init() {
	// TODO: Support for custom configuration

	// Get home directory environnement variable
	homeDir := os.Getenv("HOME")

	// Create .todo directory if it doesn't exists
	if _, err := os.Stat(homeDir + "/.todo"); os.IsNotExist(err) {
		fmt.Println("INFO : Directory " + homeDir + "/.todo does not exists")
		fmt.Println("INFO : Creating " + homeDir + "/.todo directory")
		os.Mkdir(homeDir+"/.todo", os.ModePerm)
	}

	todoDirPath = homeDir + "/.todo"
	todoFilePath = todoDirPath + "/todo.txt"
	todoArchivePath = todoDirPath + "/archive.txt"
}

// PrintHelp : print help / usage message
func PrintHelp() {
	fmt.Printf("\nUsage : todo <OPTION> <TASK> \n\n")
	fmt.Printf("Options : \n\n")
	fmt.Println("	todo add <TASK>		Add the task <TASK> to the list")
	fmt.Println("	todo do <IDTASK>	Check the <IDTASK> task")
	fmt.Println("	todo undo <IDTASK>	Uncheck the <IDTASK> task")
	fmt.Println("	todo ls			Display all tasks")
	fmt.Println("	todo ls <TAG>		Display all tasks containing the <TAG>")
}

func main() {
	Init()
	args := os.Args[1:]
	// TODO: Move args to function parsing outside of main
	if len(args) > 0 {
		if args[0] == "ls" {
			if len(args) > 1 {
				ListTasksFiltered(args[1:])
			} else {
				ListTasks()
			}
		} else if args[0] == "archive" {
			ArchiveTasks()
		} else if args[0] == "help" {
			PrintHelp()
		} else if args[0] == "add" && len(args) > 1 {
			AddTask(todoFilePath, args[1])
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
		} else if (args[0] == "delete") && len(args) > 1 {
			id, casterr := strconv.ParseUint(args[1], 10, 16)
			if casterr != nil {
				log.Fatal(casterr)
			} else {
				RemoveTask(uint16(id))
			}
		} else {
			fmt.Println("ERROR : missing / wrong arguments")
			PrintHelp()
		}
	}
}
