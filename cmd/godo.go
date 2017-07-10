package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
	"strings"
	"sort"
	"flag"
	"github.com/fchris/godo/parse"
)

type TaskList []parse.Task

func (t TaskList) Len() int {
	return len(t)
}

func (t TaskList) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t TaskList) Less(i, j int) bool {
	return t[i].DueTime.Before(t[j].DueTime)
}

var taskList TaskList

func main() {

	interactive := flag.Bool("interactive", false, "Trigger interactive mode")
	fileName := flag.String("f", "", "Filename to read from and write to")
	flag.Parse()

	fmt.Println(*fileName)

	taskList = parseFromFile(*fileName)
	fmt.Println(taskList)

	if *interactive {
		shell(*fileName)
	}
}

func shell(fileName string) {
	run := true
	for run {
		fmt.Println("Enter action:")

		reader := bufio.NewReader(os.Stdin)
		action, err := reader.ReadString('\n')

		if err != nil {
			log.Fatal("Could not read action")
			return
		}

		action = strings.Trim(action, "\n")

		switch action {
		case "a", "A", "add":
			fmt.Println("What is the description:")
			desc, err := reader.ReadString('\n')
			desc = strings.Trim(desc, "\n")
			if err != nil {
				log.Fatal("Could not read description")
				return
			}

			fmt.Println("When is it due? Hit enter for today:")
			dateString, err := reader.ReadString('\n')
			dateString = strings.Trim(dateString, "\n")
			if err != nil {
				log.Fatal("Could not read due date")
				continue
			}

			dueTime := time.Now()
			if len(dateString) != 0 {
				dueTime, err = time.Parse("01.02.06", dateString)
				if err != nil {
					log.Fatal("Could not parse date")
				}
			}

			task := parse.Task{desc, dueTime, false}
			taskList = append(taskList, task)

		case "p", "P", "print":
			var currentDate time.Time

			sort.Sort(taskList)

			for _, task := range taskList {
				if currentDate.YearDay() != task.DueTime.YearDay() || currentDate.Year() != task.DueTime.Year() {
					fmt.Println()
					currentDate = task.DueTime
					dateString := currentDate.Format("01.02.06")
					fmt.Println(dateString)
				}
				fmt.Println(task.Description)
			}

		case "q", "Q", "quit":
			run = false
		}
	}

	saveToFile(taskList, fileName)
}

func parseFromFile(fileName string) (list TaskList) {
	file, error := os.OpenFile(fileName, os.O_RDONLY, 0600)
	if error != nil {
		panic(error)

	}
	defer file.Close()

	parser := parse.NewParser(file)

	list = make(TaskList, 0, 10)
	var task *parse.Task
	for {
		task, error = parser.Parse()
		if error != nil {
			return
		}
		list = append(list, *task)
	}

	return list
}

func saveToFile(list TaskList, fileName string) {
	file, error := os.OpenFile(fileName, os.O_WRONLY, 0600)
	if error != nil {
		panic(error)

	}
	defer file.Close()

	sort.Sort(list)
	for _, task := range list {
		fmt.Println()
		dateString := task.DueTime.Format("01.02.06")

		fmt.Println(task)

		if task.Complete {
			file.WriteString("[X] " + task.Description + " {" + dateString + "}" + "\n")
		} else {
			file.WriteString("[ ] " + task.Description + " {" + dateString + "}" + "\n")
		}
	}
}
