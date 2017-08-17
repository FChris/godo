# towg - Todos with Go
## A small go tool to manage todo files

[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat-square)](http://makeapullrequest.com)
[![Go Report Card](https://goreportcard.com/badge/github.com/fchris/towg)](https://goreportcard.com/report/github.com/fchris/towg)

## Background

This is work in progress, it may have some breaking changes until declared stable and probably does not provide all the functionality you might expect yet.  

The main reason for developing this project instead of using [todo.txt](http://todotxt.com/) was a matter of taste and a chance to learn the Go programming language. 

The main goals for this tool are:
* manage lists in a markdown like style
* be able to easily convert todo lists into printable pdf 
* view todo lists with text editor only
* add todos without opening or navigating through existing todo lists

Some of the features are already implemented, some of those still need some work.

Also the development and testing so far has only been done on Linux. So there might be issues on other systems. 
If you find any bugs or misbehaviour you can file a bug report or provide a fix right away. Take a look at the [contributing guidelines](https://github.com/FChris/towg/blob/master/CONTRIBUTING.md) for that.

## Installation

If you have your go environment set up already, you can simply execute

```go get github.com/fchris/towg/cmd```

and run it by executing

```go run $GOPATH/src/github.com/fchris/towg/cmd/towg.go ```

Alternatively, you can install it into your `$GOBIN` folder by navigating to 

```$GOPATH/src/github.com/fchris/towg/```

and executing 

```go install cmd/towg.go```

If your `$GOBIN` folder is part of your `$PATH` you can simply execute it as shown in [Usage](#usage),
otherwise you can execute it with `$GOBIN/towg`.

In case you haven't set up go yet you can follow the instructions provided by the [official site](https://golang.org/doc/install)

## Usage

towg works with a markdown like structure for Todo Lists.

Todos are sorted by date in the format dd.mm.yy and by completeness

An Example Todo file would look like this:
```
  # 17.07.17
  
- [X] Todo 25
- [X] Todo 27
- [ ] Todo 2
- [ ] Todo 7

  # 14.07.17

- [X] Todo 2
- [X] Todo 27
- [ ] Todo 25
- [ ] Todo 7
```

Files are given to towg with the -f flag.

If you want to print all Todos in the file you can use the print subcommand flag. For example:  
  ``` towg print -f mytodolist.todo -d -```
  
If you want to print all Todos for a given day you can specify a date with -d. For example:  
  ``` towg print -f mytodolist.todo -d 10.10.17```  
  ``` towg print -f mytodolist.todo -d yesterday```  
  ``` towg print -f mytodolist.todo -d today```  
  ``` towg print -f mytodolist.todo -d tomorrow```  

If you want to set a Todo to complete you have to use the -s flag and an id provided with -i:  
The id indicates the i-th Todo in the list you would see if you printed it. For example:  
   ```towg switch -f mytodolist.todo -n 5 ``` Switches the status of the 5th entry in the whole file.  
   ```towg switch -f mytodolist.todo -n 4 -d today``` Switches the status of the 5th entry in the list for today.  
   
Other subcommands are redate, delete and add. They work similarly the commands described above.
You can get an info text for all commands by executing ```towg <command> -h```.  

Todolists stay the same unless a status is switched (or in later versions the description changes). Therefor it is 
suggested that you simply first print the list for a given date to find out the position of your todo and then switch the status.

## Contributing

I would love to hear your feedback and input. Check out the [contributing guidelines](https://github.com/FChris/towg/blob/master/CONTRIBUTING.md) for ways to contribute.

## Further Work
* Also I am looking for a new name as towg is used already a lot in other projects.
