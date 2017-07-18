package task

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"sort"
	"time"
)

func TestTodoList_InsertTodo(t *testing.T) {
	var todoList TodoList
	todo := Todo{Description: "Test", Complete: false}
	todoList.InsertTodo(todo)

	expectedTodoList := TodoList{todo}

	assert.Equal(
		t,
		expectedTodoList,
		todoList,
		"Actual TodoList different from expected after insert")

	todoList.InsertTodo(todo)
	assert.Equal(
		t,
		expectedTodoList,
		todoList,
		"Actual TodoList changed after re-inserting todo which was already in the list.")
}

func TestTodoList_Insert(t *testing.T) {
	var todoList TodoList
	todo1 := Todo{Description: "Test", Complete: false}
	todo2 := Todo{Description: "Test1", Complete: true}
	insertableList := TodoList{todo1, todo2}
	todoList.Insert(insertableList)

	expectedTodoList := TodoList{todo1, todo2}

	sort.Sort(expectedTodoList)
	sort.Sort(todoList)
	assert.Equal(
		t,
		expectedTodoList,
		todoList,
		"Actual TodoList different from expected after insert")
}

func TestDayList_HasDate(t *testing.T) {
	testDate1, err := time.Parse("02.01.06", "01.01.20")
	assert.Equal(
		t,
		nil,
		err,
		"Error for parsing date is not nil")

	todo1 := Todo{Description: "Test", Complete: false}
	todo2 := Todo{Description: "Test1", Complete: true}
	todoList := TodoList{todo1, todo2}

	day := Day{testDate1, todoList}

	dayList := DayList{day}

	assert.True(
		t,
		dayList.HasDate(testDate1),
		"DayList contains a day for the give date but HasDate returns false.")
}

func TestDayList_SetDay(t *testing.T) {
	testDate1, err := time.Parse("02.01.06", "01.01.20")
	assert.Equal(t, nil, err, "Error for parsing date is not nil")

	todo1 := Todo{Description: "Test", Complete: false}
	todo2 := Todo{Description: "Test1", Complete: true}
	todoList := TodoList{todo1, todo2}

	day := Day{testDate1, todoList}

	var dayList DayList
	dayList.SetDay(day)

	assert.True(
		t,
		dayList.HasDate(testDate1),
		"DayList does not contain day after SetDay has been called.")

	assert.Equal(
		t,
		DayList{day},
		dayList,
		"DayList does not contain day after SetDay has been called.")

	dayList.SetDay(day)
	assert.NotEqual(
		t,
		DayList{day, day},
		dayList,
		"DayList contains same day twice after calling SetDay multiple times for the same day")

	todo3 := Todo{Description: "Test3", Complete: true}
	day.Todos.InsertTodo(todo3)
	dayList.SetDay(day)
	assert.Equal(
		t,
		DayList{day},
		dayList,
		"DayList does not contain updated day after inserting a new Todo into a day.")

}
