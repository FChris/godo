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
	"io"
)

var dayList task.DayList

func main() {

	interactive := flag.Bool("interactive", false, "Trigger interactive mode")
	fileName := flag.String("f", "tasks.todo", "Filename to read from and write to")
	add := flag.String("a", "", "Add Todo given in the format '# <Date> [ ] <Task Description>")
	complete := flag.String("c", "", "Action which shall be executed")
	flag.Parse()

	fmt.Println(*fileName)

	dayList = parseFromFile(*fileName)
	fmt.Println(dayList)

	if *add != "" {
		parseData(strings.NewReader(*add))
	}

	if *complete != "" {
		//TODO implement
	}

	if *interactive {
		shell()
	}

	save(dayList, *fileName)
}

func shell() {
	run := true
	for run {
		fmt.Println("Enter action:")

		reader := bufio.NewReader(os.Stdin)
		action, err := reader.ReadString('\n')

		if err != nil {
			log.Fatal("Could not read action")
			continue
		}
		action = strings.Trim(action, "\n")

		switch action {
		case "a", "A", "add":
			err = addTask(*reader)
			if err != nil {
				continue
			}
		case "p", "P", "print":
			printDayList()
		case "q", "Q", "quit":
			run = false
		}
	}
}

func addTask(reader bufio.Reader) error {
	fmt.Println("What is the description:")
	desc, err := reader.ReadString('\n')
	desc = strings.Trim(desc, "\n")
	if err != nil {
		log.Fatal("Could not read description")
		return err
	}

	fmt.Println("When is it due? Hit enter for today:")
	dateString, err := reader.ReadString('\n')
	dateString = strings.Trim(dateString, "\n")
	if err != nil {
		log.Fatal("Could not read due date")
		return err
	}

	dueTime := time.Now()
	if len(dateString) != 0 {
		dueTime, err = time.Parse("01.02.06", dateString)
		if err != nil {
			log.Fatal("Could not parse date")
			return err
		}
	}

	task := task.Todo{desc, false}

	day := dayList.DayByDate(dueTime)
	fmt.Println(day)
	day.Todos = append(day.Todos, task)
	fmt.Println(day)
	dayList = dayList.SetDay(day)

	return nil
}

func printDayList() {
	sort.Sort(dayList)
	for _, day := range dayList {
		sort.Sort(day.Todos)
		fmt.Println()
		dateString := day.Date.Format("01.02.06")
		fmt.Println(dateString)
		for _, todo := range day.Todos {
			fmt.Println(todo.Description)
		}
	}
}

func parseFromFile(fileName string) (list task.DayList) {
	file, error := os.OpenFile(fileName, os.O_RDONLY, 0600)
	if error != nil {
		//FIXME Write proper error message. Obviously we cannot open this file
		return
		//panic(error)
	}
	defer file.Close()

	return parseData(file)
}

func parseData(r io.Reader) (list task.DayList) {
	parser := parse.NewParser(r)

	list = make(task.DayList, 0, 10)
	for {
		day, error := parser.Parse()
		if error != nil {
			return
		}
		list = append(list, *day)
	}

	return list
}

func save(dayList task.DayList, fileName string) {
	file, error := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0600)
	if error != nil {
		panic(error)

	}
	defer file.Close()

	sort.Sort(dayList)

	for _, day := range dayList {
		sort.Sort(day.Todos)
		dateString := day.Date.Format("01.02.06")
		file.WriteString("# " + dateString + "\n")

		for _, todo := range day.Todos {

			if todo.Complete {
				file.WriteString("[X] " + todo.Description)
			} else {
				file.WriteString("[ ] " + todo.Description)
			}

			file.WriteString("\n")
		}
	}
}
