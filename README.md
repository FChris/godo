[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat-square)](http://makeapullrequest.com)

A small go tool to manage my todo files

# Background

This is work in progress, it may have some breaking changes until declared stable and probably does not provide all the functionality you might expect yet.  

The main reason for developing this project instead of using [todo.txt] (http://todotxt.com/)was a matter of taste and a chance to learn the Go programming language. 

The main goals for this tool are:
* manage lists in a markdown like style
* be able to easily convert todo lists into printable pdf 
* view todo lists with text editor only
* add todos without opening or navigating through existing todo lists

Some of the features are already implemented, some of those stil need some work.

Also the development and testing so far has only been done on linux. So there might be issues on other systems. 
If you find any bugs or misbehaviour you can file a bug report or provide a fix right away. Tkae a look at the [contributing guidelines](https://github.com/FChris/godo/edit/master/CONTRIBUTING.md) for that.

# Installation

If you have your go enviroment set up already you can simply execute

```go get github.com/fchris/godo/cmd```

and run it by executing to 

```go run $GOPATH/src/github.com/fchris/godo/cmd/godo.go ```

Alternatively you can install it into your $GOBIN folder by navigating to 

```$GOPATH/src/github.com/fchris/godo/```

and executing 

```go install cmd/godo.go```

If your $GOBIN folder is part of your $PATH you can simply execute it as shown under Usage,
otherwise you can execute it with $GOBIN/godo.

In case you haven't set up go yet you can follow the instructions provided by the [official site](https://golang.org/doc/install)

# Usage

Godo works with a markdown like structure for Todo Lists.

Todos are sorted by date in the format dd.mm.yy

An Example Todo file would look like this:
```
  # 17.07.17
  
  [X] Todo 25
  [X] Todo 27
  [ ] Todo 2
  [ ] Todo 7

  # 14.07.17

  [X] Todo 2
  [X] Todo 27
  [ ] Todo 25
  [ ] Todo 7
```

Files are given to godo with the -f flag.

If you want to print all Todos in the file you can use the -p flag. For example:  
  ``` godo -f mytodolist.todo -p```
  
If you want to print all Todos for a given day you can specify a date with -d. For example:  
  ``` godo -f mytodolist.todo -p -d 10.10.17```  
  ``` godo -f mytodolist.todo -p -d yesterday```  
  ``` godo -f mytodolist.todo -p -d today```  
  ``` godo -f mytodolist.todo -p -d tomorrow```  

If you want to set a Todo to complete you have to use the -s flag and an id provided with -i:  
The id indicates the i-th Todo in the list you would see if you printed it. For example:  
   ```godo -f mytodolist.todo -i 5 -s``` Switches the status of the 5th entry in the whole file.  
   ```godo -f mytodolist.todo -i 4 -s -d today``` Switches the status of the 5th entry in the list for today.  
   
Todolists stay the same unless a status is switched (or in later versions the description changes). Therefor it is 
suggested that you simply first print the list for a given date to find out the position of your todo and then switch the status.

# Contributing

I would love to hear your feedback and input. Check out the [contributing guidelines](https://github.com/FChris/godo/edit/master/CONTRIBUTING.md) for ways to contribute.

# Further Work
The tool will be expanded to provide edit, delete and re-date functionality and the code basis will have to be cleaned up.
