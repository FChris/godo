package parse

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"github.com/fchris/godo/task"
	"time"
	"testing"
)

func TestParseString(t *testing.T) {
	parser := NewParser(strings.NewReader("# 01.01.20 [ ] Test String"))
	day, err := parser.Parse()

	assert.Equal(t, nil, err, "Error is not nil")

	testDate, err := time.Parse(Timeformat, "01.01.20")

	assert.Equal(t, nil, err, "Error for parsing date is not nil")

	testTodoList := task.TodoList{Elem: []task.Todo{task.Todo{Description: "Test String", Complete: false}}}
	testDay := task.Day{testDate, testTodoList}

	assert.Equal(t, testDay, day, "Test Day does not equal actual parsed day")
}

func TestParseMultilineString(t *testing.T) {
	parser := NewParser(strings.NewReader("# 01.01.20 [ ] Test String\n" +
		"[ ] Test String2 \n" +
		"# 01.02.20 [ ] Test String2\n"))

	day, err := parser.Parse()
	assert.Equal(t, nil, err, "Error is not nil")
	testDate1, err := time.Parse(Timeformat, "01.01.20")
	assert.Equal(t, nil, err, "Error for parsing date is not nil")
	testTodoList := task.TodoList{Elem: []task.Todo{task.Todo{Description: "Test String", Complete: false},
													task.Todo{Description: "Test String2", Complete:false}}}
	testDay := task.Day{Date: testDate1, Todos: testTodoList}

	assert.Equal(t, testDay, day, "Test Day does not equal actual parsed day")
	day, err = parser.Parse()
	assert.Equal(t, nil, err, "Error is not nil")
	testDate2, err := time.Parse(Timeformat, "01.02.20")
	assert.Equal(t, nil, err, "Error for parsing date is not nil")
	testTodoList2 := task.TodoList{Elem: []task.Todo{{Description: "Test String2", Complete: false}}}
	testDay2 := task.Day{Date: testDate2, Todos: testTodoList2}
	assert.Equal(t, testDay2, day, "Test Day 2 does not equal actual parsed day")

}
