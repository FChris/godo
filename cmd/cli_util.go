package cmd

import (
	"fmt"
	"github.com/fchris/towg/parse"
	"github.com/fchris/towg/task"
	"io"
	"os"
	"sort"
	"strings"
	"time"
)

const (
	yesterday       string = "yesterday"
	today           string = "today"
	tomorrow        string = "tomorrow"
	fileNameDefault string = "tasks.todo"
)

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

// addTodoFromDesc returns an updated original list with a new todo based on desc inserted into day with date or an
// error and an unchanged original in case something goes wrong
func addTodoFromDesc(original task.DayList, desc string, date string) (task.DayList, error) {
	var d time.Time
	var err error
	if isRelativeDayDescription(date) {
		d = dateByRelativeDayDescription(date)
	} else {
		d, err = time.Parse(parse.Timeformat, date)
		if err != nil {
			return original, err
		}
	}

	todoString := "# " + d.Format(parse.Timeformat) + " - [ ] " + desc

	parsedDayList, err := parseData(strings.NewReader(todoString))

	if err != nil {
		err = fmt.Errorf("Parsing from description: %s", err)
		return original, err
	}

	return addDayList(original, parsedDayList), err
}

func addDayList(original, new task.DayList) task.DayList {
	for _, newDay := range new {
		d := original.DayByDate(newDay.Date)
		d.Todos.Insert(newDay.Todos)
		original.SetDay(d)
	}

	return original
}

func switchTodoStatus(original, new task.DayList, id int) {
	for i, day := range new {
		for j, todo := range day.Todos {
			if i+j+1 == id {
				todo.Complete = !todo.Complete
				day.Todos.InsertTodo(todo)

				original.SetDay(day)
				return
			}
		}
	}
}

// deleteTodo removes the n-th todo from the new list and then updates the original list with all the remaining
// todos from new and then returns original then
//
// Note that new has to be a sublist of original
func deleteTodo(original task.DayList, subList task.DayList, n int) (new task.DayList, err error) {
	new = original
	d, ind, err := subList.TransformToDayBasedIndex(n)
	if err != nil {
		return
	}
	err = new.DeleteTodo(d, ind)
	return
}

func changeDateOfTodo(original task.DayList, new task.DayList, dayDescription string, n int) (task.DayList, error) {
	var date time.Time

	if isRelativeDayDescription(dayDescription) {
		date = dateByRelativeDayDescription(dayDescription)
	} else {
		var err error
		date, err = time.Parse(parse.Timeformat, dayDescription)
		if err != nil {
			err = fmt.Errorf("Error while parsing from date: %s", err)
			return original, err
		}
	}

	d, ind, err := new.TransformToDayBasedIndex(n - 1)
	if err != nil {
		return original, fmt.Errorf("Error while retrieving specified todo: %s", err)
	}

	todo := original.DayByDate(d).Todos[ind]
	err = original.DeleteTodo(d, ind)
	if err != nil {
		return original, fmt.Errorf("Error while deleting todo from old day: %s", err)
	}

	original.InsertTodo(date, todo)
	return original, nil
}

func dayListByPeriod(original task.DayList, period string) (task.DayList, error) {
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

	sort.Sort(original)

	var periodDayList task.DayList
	for _, day := range original {

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
