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

func (t TodoList) Less(i, j int) bool {
	if t[i].Complete && !t[j].Complete {
		return true
	} else if !t[i].Complete && t[j].Complete {
		return false
	} else {
		return strings.Compare(t[i].Description, t[j].Description) == -1
	}
}


//SetTodo checks if a Todo is already in the todo list and if not adds it
func (t *TodoList) InsertTodo(td Todo) {
	for i, todo := range *t {
		if todo.Description == td.Description {
			if todo.Complete != td.Complete {
				newList := append((*t)[:i], td)
				if len(*t) - 1 > i {
					newList = append(newList, (*t)[i+1:]...)
				}
				*t = newList
				sort.Sort(t)
			}
			return
		}
	}
	newList := append(*t, td)
	*t = newList
	sort.Sort(t)
}

func (t *TodoList) Insert(tl TodoList) {
	for _, todo := range tl {
		t.InsertTodo(todo)
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
	return t[i].Date.After(t[j].Date)
}

func (t DayList) HasDate(date time.Time) bool {
	for _, d := range t{
		if date == d.Date {
			return true
		}
	}
	return false
}

//DayByDate returns a copy of the day for the given date from the DayList or a newly initialized day
//for the date if the list does not contain it yet
func (t DayList) DayByDate(date time.Time) Day {
	for _, d := range t {
		if date.YearDay() == d.Date.YearDay() && date.Year() == d.Date.Year() {
			return d
		}
	}

	return Day{date, TodoList{}}
}

// SetDay inserts a given day into the DayList. If the day is already in the list it is overwritten otherwise
// it is simply appended
func (t *DayList) SetDay(day Day) {
	var newList []Day
	if t.HasDate(day.Date) {
		for _, d := range *t {
			if d.Date.Year() == day.Date.Year() && day.Date.YearDay() == d.Date.YearDay() {
				d.Todos.Insert(day.Todos)
			}
			newList = append(newList, d)
		}
	} else {
		newList = append(*t, day)
	}

	*t= newList
}
