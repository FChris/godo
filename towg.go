package main

import (
	"github.com/FChris/towg/cmd"
)

func main() {
	messages := make(chan string)
	go cmd.RunCLI(messages)

	//Wait for a message to show the CLI has finished
	<-messages

}
