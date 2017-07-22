package cmd

import (
	"fmt"
	"github.com/urfave/cli"
	"os"
	"sort"
)

// RunCLI executes the Command Line Interface for GoDo
func RunCLI(messages chan string) {
	app := cli.NewApp()
	app.Name = "GoDo"
	app.Usage = "A small go tool to manage todo files"

	app.Commands = []cli.Command{
		printCommand(), addCommand(), switchStatusCommand(), deleteCommand(), redateCommand(),
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	app.Run(os.Args)

	messages <- "Finished"
}

func printCommand() cli.Command {
	return cli.Command{

		Name:  "print",
		Usage: "print all tasks for a time period",
		Flags: []cli.Flag{fileFlag(), dateFlag()},
		Action: func(c *cli.Context) error {
			fileName := c.String("file")
			if fileName == "" {
				fileName = fileNameDefault
			}
			list, err := parseFromFile(fileName)
			if err != nil {
				return err
			}
			date := c.String("date")
			if date == "" {
				date = today
			}
			periodList, err := dayListByPeriod(list, date)
			if err != nil {
				fmt.Println(err)
				return err
			}
			printDayList(periodList)
			return nil
		},
	}
}

func addCommand() cli.Command {
	return cli.Command{
		Name:  "add",
		Usage: "add given text as a todo",
		Flags: []cli.Flag{
			fileFlag(),
			dateFlag(),
			cli.StringFlag{
				Name:  "text, t",
				Usage: "text describing the task",
			},
		},
		Action: func(c *cli.Context) error {
			fileName := c.String("file")
			if fileName == "" {
				fileName = fileNameDefault
			}
			list, err := parseFromFile(fileName)
			if err != nil {
				fmt.Println(err)
				return err
			}

			date := c.String("date")
			if date == "" {
				date = today
			}
			text := c.String("text")
			list, err = addTodoFromDesc(list, text, date)
			if err != nil {
				fmt.Println(err)
				return err
			}
			save(list, fileName)
			if err != nil {
				fmt.Println(err)
				return err
			}
			return nil
		},
	}
}

func switchStatusCommand() cli.Command {
	return cli.Command{
		Name:  "switch",
		Usage: "switches the status of the n-th todo in the list of todos for the given date",
		Flags: []cli.Flag{
			fileFlag(),
			dateFlag(),
			cli.IntFlag{
				Name:  "number, n",
				Usage: "number of the todo of which the status has to be switched",
			},
		},
		Action: func(c *cli.Context) error {
			fileName := c.String("file")
			if fileName == "" {
				fileName = fileNameDefault
			}
			listByFile, err := parseFromFile(fileName)
			if err != nil {
				fmt.Println(err)
				return err
			}

			date := c.String("date")
			if date == "" {
				date = today
			}
			listByPeriod, err := dayListByPeriod(listByFile, date)
			if err != nil {
				fmt.Println(err)
				return err
			}

			number := c.Int("number")
			switchTodoStatus(listByFile, listByPeriod, number)
			save(listByFile, fileName)
			if err != nil {
				fmt.Println(err)
				return err
			}
			return nil
		},
	}
}

func deleteCommand() cli.Command {
	return cli.Command{
		Name:  "delete",
		Usage: "switches the status of the n-th todo in the list of todos for the given date",
		Flags: []cli.Flag{
			fileFlag(),
			dateFlag(),
			cli.IntFlag{
				Name:  "number, n",
				Usage: "number of the todo of which the status has to be switched",
			},
		},
		Action: func(c *cli.Context) error {
			fileName := c.String("file")
			if fileName == "" {
				fileName = fileNameDefault
			}
			listByFile, err := parseFromFile(fileName)
			if err != nil {
				fmt.Println(err)
				return err
			}

			date := c.String("date")
			if date == "" {
				date = today
			}
			listByPeriod, err := dayListByPeriod(listByFile, date)
			if err != nil {
				fmt.Println(err)
				return err
			}

			number := c.Int("number")
			listByFile = deleteTodo(listByFile, listByPeriod, number)
			save(listByFile, fileName)
			if err != nil {
				fmt.Println(err)
				return err
			}
			return nil
		},
	}
}

func redateCommand() cli.Command {
	return cli.Command{
		Name:  "redate",
		Usage: "redates the n-th todo in the list of given todos by date to the new date",
		Flags: []cli.Flag{
			fileFlag(),
			dateFlag(),
			cli.IntFlag{
				Name:  "number, n",
				Usage: "number of the todo which will be redated",
			},
			cli.StringFlag{
				Name:  "newdate",
				Usage: "new date for the todo. Allows dates as 'dd.mm.yy',or as 'yesterday', 'today', 'tomorrow'",
			},
		},
		Action: func(c *cli.Context) error {
			fileName := c.String("file")
			if fileName == "" {
				fileName = fileNameDefault
			}
			listByFile, err := parseFromFile(fileName)
			if err != nil {
				fmt.Println(err)
				return err
			}

			date := c.String("date")
			if date == "" {
				date = today
			}
			listByPeriod, err := dayListByPeriod(listByFile, date)
			if err != nil {
				fmt.Println(err)
				return err
			}

			number := c.Int("number")
			newDate := c.String("newdate")
			if newDate == "" {
				newDate = today
			}
			listByFile, err = changeDateOfTodo(listByFile, listByPeriod, newDate, number)
			if err != nil {
				fmt.Println(err)
				return err
			}
			save(listByFile, fileName)
			return nil
		},
	}
}

func fileFlag() cli.Flag {
	return cli.StringFlag{
		Name:  "file, f",
		Usage: "Load tasks from file",
	}
}

func dateFlag() cli.Flag {
	return cli.StringFlag{
		Name: "date, d",
		Usage: "date or time for the command. Allows dates as 'dd.mm.yy', 'dd.mm.yy-dd.mm.yy' " +
			"or as 'yesterday', 'today', 'tomorrow', or '-' for all days",
	}
}
