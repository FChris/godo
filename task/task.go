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

type TodoList struct {
	Elem []Todo
}

type Day struct {
	Date  time.Time
	Todos TodoList
}

func (t TodoList) Len() int {
	return len(t.Elem)
}

func (t TodoList) Swap(i, j int) {
	t.Elem[i], t.Elem[j] = t.Elem[j], t.Elem[i]
}

//SetTodo checks if a Todo is already in the todo list and if not adds it
func (t *TodoList) InsertTodo(td Todo) {
	for i, todo := range t.Elem {
		if todo.Description == td.Description {
			if todo.Complete != td.Complete {
				newList := append(t.Elem[:i], td)
				if len(t.Elem) - 1 > i {
					newList = append(newList, t.Elem[i+1:]...)
				}
				t.Elem = newList
				sort.Sort(t)
			}
			return
		}
	}
	newList := append(t.Elem, td)
	t.Elem = newList
	sort.Sort(t)
}

func (t *TodoList) Insert(tl TodoList) {
	for _, todo := range tl.Elem {
		t.InsertTodo(todo)
	}
}

func (t TodoList) Less(i, j int) bool {
	if t.Elem[i].Complete && !t.Elem[j].Complete {
		return true
	} else if !t.Elem[i].Complete && t.Elem[j].Complete {
		return false
	} else {
		return strings.Compare(t.Elem[i].Description, t.Elem[j].Description) == -1
	}
}

type DayList struct {
	Elem []Day
}

func (t DayList) Len() int {
	return len(t.Elem)
}

func (t DayList) Swap(i, j int) {
	t.Elem[i], t.Elem[j] = t.Elem[j], t.Elem[i]
}

func (t DayList) Less(i, j int) bool {
	return t.Elem[i].Date.After(t.Elem[j].Date)
}

func (t DayList) HasDate(date time.Time) bool {
	for _, d := range t.Elem {
		if date == d.Date {
			return true
		}
	}
	return false
}

//getDay returns the day for the give date from the DayList or a newly initialized day
//for the date if the list does not contain it yet
func (t DayList) DayByDate(date time.Time) Day {
	for _, d := range t.Elem {
		if date.YearDay() == d.Date.YearDay() && date.Year() == d.Date.Year() {
			return d
		}
	}

	return Day{date, TodoList{}}
}

//SetDay creates a new list which the original day of this date is replaced. If the list did not contain this day
//yet it is simply appended
func (t *DayList) SetDay(day Day) {
	var newList []Day
	if t.HasDate(day.Date) {
		for _, d := range t.Elem {
			if d.Date.Year() == day.Date.Year() && day.Date.YearDay() == day.Date.YearDay() {
				d.Todos.Insert(day.Todos)
			}
			newList = append(newList, d)
		}
	} else {
		newList = append(t.Elem, day)
	}

	t.Elem = newList
}
