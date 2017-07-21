package main

import (
	"flag"
	"fmt"
	"github.com/fchris/godo/parse"
	"github.com/fchris/godo/task"
	"io"
	"os"
	"sort"
	"strings"
	"time"
)

var dayList task.DayList

const (
	yesterday string = "yesterday"
	today     string = "today"
	tomorrow  string = "tomorrow"
)

func main() {

	fileName := flag.String("f", "tasks.todo", "Filename to read from and write to")
	period := flag.String("d", "-", "Accepts a time which is used by print or complete or add. "+
		"Add requires a specific date while print and switch can work with time periods."+
		"E.g - for all, "+
		"01.02.06 for the ones on this date, "+
		"01.02.06-31.12.06 for all from the first to the second period "+
		"yesterday, today, tomorrow")

	add := flag.String("add", "", "Add the given text as a todo. Needs to be combined with the d flag")
	switchStatus := flag.Int("switch", 0, "Accepts the number of the entry of which the status is completed.")
	printDays := flag.Bool("print", false, "Prints all todos for the given time.")
	del := flag.Int("delete", 0, "Accepts the number of the entry which shall be deleted")
	flag.Parse()

	dayList, err := parseFromFile(*fileName)
	if err != nil {
		fmt.Println(err)
		return
	}

	if *add != "" {
		err := addTodoFromDesc(*add, *period)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	if *printDays {
		list, err := dayListByPeriod(dayList, *period)
		if err != nil {
			fmt.Println(err)
			return
		}
		printDayList(list)
	}

	if *switchStatus > 0 {
		list, err := dayListByPeriod(dayList, *period)
		if err != nil {
			fmt.Println(err)
			return
		}
		switchTodoStatus(list, *switchStatus)
	}

	if *del > 0 {
		list, err := dayListByPeriod(dayList, *period)
		if err != nil {
			fmt.Println(err)
			return
		}
		deleteTodo(list, *del)
	}

	err = save(dayList, *fileName)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func parseFromFile(fileName string) (list task.DayList, err error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY, 0600)
	if err != nil {
		err = fmt.Errorf("Error while opening file: %s", err)
		return list, err

	}
	defer file.Close()

	list, err = parseData(file)
	if err != nil {
		err = fmt.Errorf("Parsing from file: %s", err)
	}
	return list, err
}

func parseData(r io.Reader) (list task.DayList, err error) {
	parser := parse.NewParser(r)

	for {
		day, e := parser.Parse()
		if e != nil {
			err = fmt.Errorf("Error while parsing data: %s", e)
			return
		}

		if day.Date.IsZero() {
			return
		}

		list.SetDay(day)
	}

	return
}

func save(dayList task.DayList, fileName string) error {
	err := os.Remove("." + fileName + ".bak")
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("Error while deleting old backup: %s", err)
	}

	err = os.Rename(fileName, "."+fileName+".bak")
	if err != nil {
		return fmt.Errorf("Backing up existing todo list: %s", err)
	}

	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return fmt.Errorf("Error while opening file for writing : %s", err)
	}
	defer file.Close()

	sort.Sort(dayList)

	for _, day := range dayList {
		sort.Sort(day.Todos)
		dateString := day.Date.Format(parse.Timeformat)
		file.WriteString("\n# " + dateString + "\n\n")

		for _, todo := range day.Todos {
			file.WriteString(todo.String())
			file.WriteString("  \n")
		}
	}

	return nil
}

func addTodoFromDesc(desc string, date string) error {
	var d time.Time
	var err error
	if isRelativeDayDescription(date) {
		d = dateByRelativeDayDescription(date)
	} else {
		d, err = time.Parse(parse.Timeformat, date)
		if err != nil {
			return err
		}
	}

	todoString := "# " + d.Format(parse.Timeformat) + " [ ] " + desc

	parsedDayList, err := parseData(strings.NewReader(todoString))

	if err != nil {
		err = fmt.Errorf("Parsing from description: %s", err)
		return err
	}

	addDayList(parsedDayList)
	return nil
}

func switchTodoStatus(l task.DayList, id int) {
	for i, day := range l {
		for j, todo := range day.Todos {
			if i+j+1 == id {
				todo.Complete = !todo.Complete
				day.Todos.InsertTodo(todo)

				dayList.SetDay(day)
				return
			}
		}
	}
}

func deleteTodo(l task.DayList, id int) {
	var newDay task.Day
	for i, day := range l {
		for j := range day.Todos {
			if i+j+1 == id {
				newDay.Date = day.Date
			}
		}
	}

	//Todo find more efficient way to do this instead of running through the whole list twice
	for i, day := range l {
		for j, todo := range day.Todos {
			if i+j+1 != id && day.Date == newDay.Date {
				newDay.Todos.InsertTodo(todo)
			}
		}
	}

	dayList.SetDay(newDay)
}

func addDayList(list task.DayList) {
	for _, newDay := range list {
		d := dayList.DayByDate(newDay.Date)
		d.Todos.Insert(newDay.Todos)
		dayList.SetDay(d)
	}
}

func dayListByPeriod(dayList task.DayList, period string) (task.DayList, error) {
	dayDescription := strings.ToLower(period)
	var fromDate time.Time
	var toDate time.Time
	var err error

	if isRelativeDayDescription(dayDescription) {
		fromDate = dateByRelativeDayDescription(dayDescription)
		toDate = fromDate
	} else if strings.IndexRune(period, '-') >= 0 {
		timeFrame := strings.Split(period, "-")
		timeFrame = deleteEmpty(timeFrame)

		if len(timeFrame) == 0 {
			toDate = time.Now().AddDate(100, 0, 0)
		} else if len(timeFrame) > 1 {
			toDate, err = time.Parse(parse.Timeformat, timeFrame[1])
			if err != nil {
				err = fmt.Errorf("Error while parsing to date: %s", err)
				return task.DayList{}, err
			}
		} else if len(timeFrame) > 0 {
			fromDate, err = time.Parse(parse.Timeformat, timeFrame[0])
			if err != nil {
				err = fmt.Errorf("Error while parsing from date: %s", err)
				return task.DayList{}, err
			}
		}
	} else {
		fromDate, err = time.Parse(parse.Timeformat, period)
		if err != nil {
			err = fmt.Errorf("Error while parsing from date: %s", err)
			return task.DayList{}, err
		}
		toDate = fromDate
	}

	fromDate = ignoreTime(fromDate)
	toDate = ignoreTime(toDate)

	sort.Sort(dayList)

	var periodDayList task.DayList
	fmt.Println(dayList)
	for _, day := range dayList {

		if inTimeSpan(fromDate, toDate, day.Date) {
			periodDayList.SetDay(day)
		}
	}

	return periodDayList, nil
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
	for _, day := range list {
		fmt.Println()
		dateString := day.Date.Format(parse.Timeformat)
		fmt.Println(dateString)
		for _, todo := range day.Todos {
			fmt.Println(todo)
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
