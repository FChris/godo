package task

import (
	"time"
	"strings"
	"sort"
)

type Todo struct {
	Description string
	Complete    bool
}

type TodoList []Todo

type Day struct {
	Date  time.Time
	Todos TodoList
}

func (t TodoList) Len() int {
	return len(t)
}

func (t TodoList) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

//SetTodo checks if a Todo is already in the todo list and if not adds it
func (t TodoList) InsertTodo(td Todo) {
	for _, todo := range t {
		if strings.Compare(strings.ToLower(td.Description), strings.ToLower(todo.Description)) == 0{
			return
		}
	}

	newList := append(t, td)
	sort.Sort(newList)
	t = newList
}

func (t TodoList) Insert (tl TodoList) {
	for _, todo := range tl {
		t.InsertTodo(todo)
	}
}

func (t TodoList) Less(i, j int) bool {
	if t[i].Complete && !t[j].Complete {
		return true
	} else if !t[i].Complete && t[j].Complete {
		return false
	} else {
		return strings.Compare(t[i].Description, t[j].Description) == -1
	}
}

type DayList []Day

func (t DayList) Len() int {
	return len(t)
}

func (t DayList) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t DayList) Less(i, j int) bool {
	return t[i].Date.Before(t[j].Date)
}

func (t DayList) hasDate(date time.Time) bool {
	for _, d := range t {
		if date == d.Date {
			return true
		}
	}
	return false
}

//getDay returns the day for the give date from the DayList or a newly initialized day
//for the date if the list does not contain it yet
func (t DayList) DayByDate(date time.Time) Day {
	for _, d := range t {
		if date.YearDay() == d.Date.YearDay() && date.Year() == d.Date.Year() {
			return d
		}
	}

	return Day{date, make([]Todo, 0, 1)}
}

//SetDay creates a new list which the original day of this date is replaced. If the list did not contain this day
//yet it is simply appended
func (t DayList) SetDay(day Day) DayList {
	//TODO check if we can modify this method so we change it inplace
	for i, d := range t {
		if d.Date.Year() == day.Date.Year() && day.Date.YearDay() == day.Date.YearDay() {

			newList := append(t[:i], day)
			if t.Len() > i+1 {
				newList = append(newList, t[i+1:]...)
			}
			return newList
		}
	}

	return append(t, day)
}
