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

const yesterday string = "yesterday"
const today string = "today"
const tomorrow string = "tomorrow"

func main() {

	interactive := flag.Bool("i", false, "Trigger interactive mode")
	fileName := flag.String("f", "tasks.todo", "Filename to read from and write to")
	add := flag.String("a", "", "Add Todo given in the format '# <Date> [ ] <Task Description>")
	period := flag.String("d", "-", "Accepts a time which is used by print or complete." +
		"E.g - for all, "+
		"01.02.06 for the ones on this date, "+
		"01.02.06-31.12.06 for all from the first to the second period ")

	switchStatus := flag.Int("s", 0, "Accepts the number of the entry of which the status is completed.")
	printDays := flag.Bool("p", false, "Prints all todos for the given time.")
	flag.Parse()

	fmt.Println(*interactive)
	fmt.Println(*fileName)
	fmt.Println(*add)
	fmt.Println(*period)
	fmt.Println(*printDays)
	fmt.Println(*switchStatus)


	dayList = parseFromFile(*fileName)

	if *add != "" {
		parsedDayList := parseData(strings.NewReader(*add))
		addDayList(parsedDayList)
	}

	if *printDays {
		printDayList(dayListByPeriod(*period))
	}

	if *switchStatus > 0 {
		switchTodoStatus(dayListByPeriod(*period), *switchStatus)
	}

	if *interactive {
		shell()
	}

	save(dayList, *fileName)
}

func parseFromFile(fileName string) (list task.DayList) {
	file, err := os.OpenFile(fileName, os.O_RDONLY, 0600)
	if err != nil {
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
		day, err := parser.Parse()
		if err != nil {
			fmt.Println(err)
			return list
		}

		if day == nil {
			break
		}

		list = append(list, *day)
	}

	return list
}

func save(dayList task.DayList, fileName string) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)

	}
	defer file.Close()

	sort.Sort(dayList)

	for _, day := range dayList {
		sort.Sort(day.Todos)
		dateString := day.Date.Format(parse.Timeformat)
		file.WriteString("\n# " + dateString + "\n\n")

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

func switchTodoStatus(l task.DayList, id int) {
	for i, day := range l {
		for j, todo := range day.Todos {
			if i+j+1 == id {
				fmt.Println(todo)
				todo.Complete = !todo.Complete
				todos := day.Todos.InsertTodo(todo)
				dayList = dayList.SetDay(task.Day{Date: day.Date, Todos:todos})
				fmt.Println(dayList)
				return
			}
		}
	}
}

func addDayList(list task.DayList) {
	for _, newDay := range list {
		d := dayList.DayByDate(newDay.Date)
		d.Todos = d.Todos.Insert(newDay.Todos)

		if !dayList.HasDate(d.Date) {
			dayList = append(dayList, d)
			sort.Sort(dayList)
		} else {
			dayList = dayList.SetDay(d)
		}
	}
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
			fmt.Println("What time period do you want to print?")
			period, err := reader.ReadString('\n')
			period = strings.Trim(period, "\n")
			if err != nil {
				fmt.Println("Could not read time period")
				continue
			}

			printDayList(dayListByPeriod(period))
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
		dueTime, err = time.Parse(parse.Timeformat, dateString)
		if err != nil {
			log.Fatal("Could not parse date")
			return err
		}
	}

	todo := task.Todo{Description: desc, Complete:false}

	day := dayList.DayByDate(dueTime)
	day.Todos = day.Todos.InsertTodo(todo)
	insertList := task.DayList{day}
	addDayList(insertList)

	fmt.Println(dayList)

	return nil
}

func dayListByPeriod(period string) task.DayList {
	dayDescription := strings.ToLower(period)
	var fromDate time.Time
	var toDate time.Time

	if dayDescription == yesterday || dayDescription == today || dayDescription == tomorrow {
		switch dayDescription {
		case yesterday:
			fromDate = time.Now().AddDate(0, 0, -1)
		case today:
			fromDate = time.Now()
		case tomorrow:
			fromDate = time.Now().AddDate(0, 0, 1)
		}

		toDate = fromDate
	} else if strings.IndexRune(period, '-') >= 0 {
		var err error
		timeFrame := strings.Split(period, "-")
		timeFrame = deleteEmpty(timeFrame)

		if len(timeFrame) == 0 {
			toDate = time.Now().AddDate(100, 0, 0)
		} else if len(timeFrame) > 0 {
			fromDate, err = time.Parse(parse.Timeformat, timeFrame[0])
			if err != nil {
				panic(err)
			}
		} else if len(timeFrame) > 1 {
			toDate, err = time.Parse(parse.Timeformat, timeFrame[1])
			if err != nil {
				panic(err)
			}
		}
	}

	fromDate = ignoreTime(fromDate)
	toDate = ignoreTime(toDate)

	sort.Sort(dayList)

	fmt.Println(dayList)

	var periodDayList task.DayList
	for _, day := range dayList {

		if inTimeSpan(fromDate, toDate, day.Date) {
			periodDayList = append(periodDayList, day)
		}
	}

	return periodDayList
}

func printDayList(list task.DayList) {
	for _, day := range list {
		fmt.Println()
		dateString := day.Date.Format(parse.Timeformat)
		fmt.Println(dateString)
		for _, todo := range day.Todos {
			if todo.Complete {
				fmt.Print("[X] ")
			} else {
				fmt.Print("[ ] ")
			}
			fmt.Println(todo.Description)
		}
	}
}

func inTimeSpan(from, to, check time.Time) bool {
	return (check.After(from) && check.Before(to)) || check == to || check == from
}

func ignoreTime(date time.Time) time.Time {
	var res time.Time
	res = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)

	return res
}

func deleteEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}
