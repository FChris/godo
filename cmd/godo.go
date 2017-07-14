package main

import (
	"fmt"
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

const (
	yesterday string = "yesterday"
	today     string = "today"
	tomorrow  string = "tomorrow"
)

func main() {

	fileName := flag.String("f", "tasks.todo", "Filename to read from and write to")
	add := flag.String("a", "", "Add the given text as a todo. Needs to be combined with the d flag")
	period := flag.String("d", "-", "Accepts a time which is used by print or complete or add. " +
		"Add requires a specific date while print and switch can work with time periods."+
		"E.g - for all, "+
		"01.02.06 for the ones on this date, "+
		"01.02.06-31.12.06 for all from the first to the second period " +
		"yesterday, today, tomorrow")

	switchStatus := flag.Int("s", 0, "Accepts the number of the entry of which the status is completed.")
	printDays := flag.Bool("p", false, "Prints all todos for the given time.")
	flag.Parse()

	dayList = parseFromFile(*fileName)

	if *add != "" {
		addTodoFromDesc(*add, *period)
	}

	if *printDays {
		printDayList(dayListByPeriod(*period))
	}

	if *switchStatus > 0 {
		switchTodoStatus(dayListByPeriod(*period), *switchStatus)
	}

	save(dayList, *fileName)
}

func parseFromFile(fileName string) (list task.DayList) {
	file, err := os.OpenFile(fileName, os.O_RDONLY, 0600)
	if err != nil {
		//FIXME Write proper error message. Obviously we cannot open this file
		return
	}
	defer file.Close()

	return parseData(file)
}

func parseData(r io.Reader) (list task.DayList) {
	parser := parse.NewParser(r)

	for {
		day, err := parser.Parse()
		if err != nil {
			fmt.Println(err)
			return list
		}

		if day == nil {
			break
		}

		list.SetDay(*day)
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

	for _, day := range dayList.Elem {
		sort.Sort(day.Todos)
		dateString := day.Date.Format(parse.Timeformat)
		file.WriteString("\n# " + dateString + "\n\n")

		for _, todo := range day.Todos.Elem {

			if todo.Complete {
				file.WriteString("[X] " + todo.Description)
			} else {
				file.WriteString("[ ] " + todo.Description)
			}

			file.WriteString("\n")
		}
	}
}

func addTodoFromDesc(desc string, date string) {
	var d time.Time
	var err error
	if isRelativeDayDescription(date) {
		d = dateByRelativeDayDescription(date)
	} else {
		d, err = time.Parse(parse.Timeformat, date)
		if err != nil {
			//Todo error message
			panic(err)
		}
	}

	todoString := "# " + d.Format(parse.Timeformat) +" [ ] " + desc

	parsedDayList := parseData(strings.NewReader(todoString))
	addDayList(parsedDayList)
}

func switchTodoStatus(l task.DayList, id int) {
	for i, day := range l.Elem {
		for j, todo := range day.Todos.Elem {
			if i+j+1 == id {
				todo.Complete = !todo.Complete
				day.Todos.InsertTodo(todo)
				return
			}
		}
	}
}

func addDayList(list task.DayList) {
	for _, newDay := range list.Elem {
		d := dayList.DayByDate(newDay.Date)
		d.Todos.Insert(newDay.Todos)
		dayList.SetDay(d)
	}
}

func dayListByPeriod(period string) task.DayList {
	dayDescription := strings.ToLower(period)
	var fromDate time.Time
	var toDate time.Time

	if isRelativeDayDescription(dayDescription) {
		fromDate = dateByRelativeDayDescription(dayDescription)
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

	var periodDayList task.DayList
	for _, day := range dayList.Elem {

		if inTimeSpan(fromDate, toDate, day.Date) {
			periodDayList.SetDay(day)
		}
	}

	return periodDayList
}
func dateByRelativeDayDescription(dayDescription string) time.Time {
	var date time.Time
	switch dayDescription {
	case yesterday:
		date = time.Now().AddDate(0, 0, -1)
	case today:
		date = time.Now()
	case tomorrow:
		date = time.Now().AddDate(0, 0, 1)
	}
	return date
}
func isRelativeDayDescription(dayDescription string) bool {
	return dayDescription == yesterday || dayDescription == today || dayDescription == tomorrow
}

func printDayList(list task.DayList) {
	for _, day := range list.Elem {
		fmt.Println()
		dateString := day.Date.Format(parse.Timeformat)
		fmt.Println(dateString)
		for _, todo := range day.Todos.Elem {
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
