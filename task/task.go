package task

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

var errOob = fmt.Errorf("index out of bounds")

// Todo is the base type for all tasks we want to save
type Todo struct {
	Description string
	Complete    bool
}

func (t Todo) String() string {
	if t.Complete {
		return "- [x] " + t.Description
	}
	return "- [ ] " + t.Description
}

// TodoList is a simple list of Todos
type TodoList []Todo

// Day is composition of a date and a TodoList.
// It is supposed to contain all todos for the given date.
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

//InsertTodo checks if a Todo is already in the todo list and if not adds it
//In case the Todo is already in the list but has a different Complete Status, the todo will be overwritten
func (t *TodoList) InsertTodo(td Todo) {
	for i, todo := range *t {
		if todo.Description == td.Description {
			if todo.Complete != td.Complete {
				newList := append((*t)[:i], td)
				if len(*t)-1 > i {
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

// Insert inserts all Todos from the parameter TodoList into the list on which the method is called
func (t *TodoList) Insert(tl TodoList) {
	for _, todo := range tl {
		t.InsertTodo(todo)
	}
}

// DayList is a simple list of Days
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

// HasDate returns true if the DayList contains a date with the given date
func (t DayList) HasDate(date time.Time) bool {
	for _, d := range t {
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
	var newList DayList
	if t.HasDate(day.Date) {
		for _, d := range *t {
			if d.Date.Year() == day.Date.Year() && day.Date.YearDay() == d.Date.YearDay() {
				newList = append(newList, day)
			} else {
				newList = append(newList, d)
			}
		}
	} else {
		newList = append(*t, day)
	}

	*t = newList
}

// DeleteDay removes a day completely from the DayList
func (t *DayList) DeleteDay(date time.Time) {
	var newList []Day
	if t.HasDate(date) {
		for _, d := range *t {
			if !(d.Date.Year() == date.Year() && date.YearDay() == d.Date.YearDay()) {
				newList = append(newList, d)
			}
		}
	}
	*t = newList
}

// InsertTodo inserts the todo into the day of this DayList corresponding to the given date
func (t *DayList) InsertTodo(date time.Time, todo Todo) {
	day := t.DayByDate(date)
	day.Todos.InsertTodo(todo)
	t.SetDay(day)
	sort.Sort(t)
}

// DeleteTodo delets the todo from the day of this DayList corresponding to the given date
func (t *DayList) DeleteTodo(date time.Time, ind int) error {
	day := t.DayByDate(date)
	if ind < 0 || ind >= day.Todos.Len() {
		return errOob
	}
	day.Todos = append(day.Todos[:ind], day.Todos[ind+1:]...)
	t.SetDay(day)
	return nil
}

// TransformToDayBasedIndex transforms the given index that corresponds to this whole DayList into an index
// which corresponds the the day in the DayList for the given date
func (t *DayList) TransformToDayBasedIndex(ind int) (date time.Time, transInd int, err error) {
	if ind < 0 {
		err = errOob
		return
	}
	transInd = ind
	for _, d := range *t {
		if transInd >= d.Todos.Len() {
			transInd -= d.Todos.Len()
		} else {
			date = d.Date
			return
		}
	}
	err = errOob
	return
}
