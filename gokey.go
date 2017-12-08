package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/kevinkassimo/gokey/src/commands"
	"github.com/kevinkassimo/gokey/src/texts"
	//"io"
	"os"
)

const (
	VERSION = "0.1.1"
)

func printWelcome() {
	texts.Ok(">>>>>> GOKEY Password Manager <<<<<<\n")
	fmt.Printf("(Version: v%s)\n", VERSION)
}

func repl(handler func(s string)) {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Printf(">> ")

	for scanner.Scan() {
		handler(scanner.Text())
		fmt.Printf(">> ")
	}
}

func main() {
	// Flags
	shouldNewUserPtr := flag.Bool("n", false, "create new user")
	shouldShowUsagePtr := flag.Bool("h", false, "show usage")
	flag.Parse()

	printWelcome()

	if *shouldShowUsagePtr {
		commands.PrintHelp()
		os.Exit(0)
	}

	if len(commands.Users) <= 0 {
		texts.Warning("No users found. Creating new user...\n")
		commands.HandleNewUser()
	} else if *shouldNewUserPtr {
		fmt.Printf("Creating new user...\n")
		commands.HandleNewUser()
	}

	for commands.HandleLogin() != true {
	}

	repl(commands.Parse)
}
