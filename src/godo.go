package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
	"strings"
	"sort"
)

type Task struct {
	description string
	dueTime     time.Time
	complete    bool
}

type TaskList []Task

func (t TaskList) Len() int {
	return len(t)
}

func (t TaskList) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t TaskList) Less(i, j int) bool {
	return t[i].dueTime.Before(t[j].dueTime)
}

func main() {
	taskList := make(TaskList, 5, 10)

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

			task := Task{desc, dueTime, false}
			taskList = append(taskList, task)

		case "p", "P", "print":
			var currentDate time.Time

			sort.Sort(taskList)

			for _, task := range taskList {
				if currentDate.YearDay() != task.dueTime.YearDay() || currentDate.Year() != task.dueTime.Year() {
					fmt.Println()
					currentDate = task.dueTime
					dateString := currentDate.Format("01.02.06")
					fmt.Println(dateString)
				}
				fmt.Println(task.description)
			}

		case "q", "Q", "quit":
			run = false
		}
	}

	saveToFile(taskList);
}

func saveToFile(list TaskList) {
	file, error := os.Open(os.Args[1])
	if error != nil {
		if os.IsNotExist(error) {
			file, error = os.Create(os.Args[1])
			if error != nil {
				panic(error)
			}
		} else {
			panic(error)
		}
	}
	sort.Sort(list)
	var currentDate time.Time
	for _, task := range list {
		if currentDate.YearDay() != task.dueTime.YearDay() || currentDate.Year() != task.dueTime.Year() {
			fmt.Println()
			currentDate = task.dueTime
			dateString := currentDate.Format("01.02.06")
			file.Write([]byte(dateString))
			file.Write([]byte("\n"))
		}
		if task.complete {
			file.Write([]byte("[X]" + task.description + "\n"))
		} else {
			file.Write([]byte("[ ]" + task.description + "\n"))
		}
	}
}
