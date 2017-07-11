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
	"github.com/fchris/godo/task"
)

var dayList task.DayList

func main() {

	interactive := flag.Bool("interactive", false, "Trigger interactive mode")
	fileName := flag.String("f", "tasks.todo", "Filename to read from and write to")
	flag.Parse()

	fmt.Println(*fileName)

	dayList = parseFromFile(*fileName)
	fmt.Println(dayList)

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

			task := task.Todo{desc, false}

			day := dayList.GetDay(dueTime)
			day.Todos = append(day.Todos, task)

			dayList = append(dayList, day)

		case "p", "P", "print":
			sort.Sort(dayList)
			for _, day := range dayList {
				fmt.Println()
				dateString := day.Date.Format("01.02.06")
				fmt.Println(dateString)
				for _, todo := range day.Todos {
					fmt.Println(todo.Description)
				}
			}

		case "q", "Q", "quit":
			run = false
		}
	}

	saveToFile(dayList, fileName)
}

func parseFromFile(fileName string) (list task.DayList) {
	file, error := os.OpenFile(fileName, os.O_RDONLY, 0600)
	if error != nil {
		//FIXME Write proper error message. Obviously we cannot open this file
		return
		//panic(error)

	}
	defer file.Close()

	parser := parse.NewParser(file)

	list = make(task.DayList, 0, 10)
	var day *task.Day
	for {
		day, error = parser.Parse()
		if error != nil {
			return
		}
		list = append(list, *day)
	}

	return list
}

func saveToFile(list task.DayList, fileName string) {
	file, error := os.OpenFile(fileName, os.O_WRONLY, 0600)
	if error != nil {
		panic(error)

	}
	defer file.Close()

	sort.Sort(list)

	for _, taskDay := range list {
		fmt.Println()
		dateString := taskDay.Date.Format("01.02.06")
		file.WriteString("# " + dateString + "\n")

		for _, todo := range taskDay.Todos {

			fmt.Println(todo)

			if todo.Complete {
				file.WriteString("[X] " + todo.Description+ "\n")
			} else {
				file.WriteString("[ ] " + todo.Description + "\n")
			}
		}
	}
}
