package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/kevinkassimo/gkey/src/commands"
	"github.com/kevinkassimo/gkey/src/texts"
	//"io"
	"os"
)

func repl(handler func(s string)) {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Printf(">> ")

	for scanner.Scan() {
		handler(scanner.Text())
		fmt.Printf(">> ")
	}
}

func main() {
	shouldNewUserPtr := flag.Bool("new", false, "create new user")

	texts.Ok(">>>>>> GKEY Password Manager <<<<<<")
	fmt.Println("(Version: v0.1)")

	if len(commands.Users) <= 0 {
		texts.Warning("No users found. Creating new user...\n")
		commands.HandleNewUser()
	} else if *shouldNewUserPtr {
		commands.HandleNewUser()
	} else {
		fmt.Printf("no flags, %s\n", *shouldNewUserPtr)
	}

	for commands.HandleLogin() != true {
	}

	repl(commands.Parse)
}
