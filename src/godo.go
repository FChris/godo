package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
	"strings"
)

type task struct {
	task    string
	dueTime time.Time
}

func main() {

	run := true
	for run {
		fmt.Println("Enter action:")

		reader := bufio.NewReader(os.Stdin)
		action, err := reader.ReadString('\n')

		if err != nil {
			log.Fatal("Could not read action")
			return
		}

		action = strings.Trim(action, "\n")

		switch action {
		case "a", "A", "add":
			fmt.Println("What is the task:")
			desc, err := reader.ReadString('\n')
			desc = strings.Trim(desc, "\n")
			if err != nil {
				log.Fatal("Could not read task")
				return
			}

			fmt.Println("When is it due? Hit enter for today:")
			dateString, err := reader.ReadString('\n')
			dateString = strings.Trim(dateString, "\n")
			if err != nil {
				log.Fatal("Could not read due date")
				return
			}

			dueTime := time.Now()
			if len(dateString) != 0 {
				dueTime, err = time.Parse("01.02.06", dateString)
				if err != nil {
					log.Fatal("Could not parse date")
				}
			}

			task := task{desc, dueTime}
			fmt.Println(task)
		case "q", "Q", "quit":
			run = false
		}
	}
}
