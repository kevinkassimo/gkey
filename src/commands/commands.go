// Package  provides ...
package commands

import (
	"bytes"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/kevinkassimo/gkey/src/confirm"
	"github.com/kevinkassimo/gkey/src/entry"
	"github.com/kevinkassimo/gkey/src/files"
	"github.com/kevinkassimo/gkey/src/texts"
	"os"
)

const (
	USER    = "login"
	EXIT    = "exit"
	ADD     = "add"
	REMOVE  = "del"
	CLEAR   = "clear"
	DESTROY = "destroy"
	NEW     = "new"
	GET     = "get"
	SHOW    = "show"
	LIST    = "list"
	COPY    = "copy"
	WHO     = "who"
	HELP    = "help"
)

// GLOB
var DataCache entry.UserEntry
var Users []string

func init() {
	Users = files.ScanAllUsers()
}

func Parse(s string) {
	args, ok := GetArgs(s)
	if !ok {
		// handle empty input
		return
	}

	CommandDispatcher(args)
}

func GetArgs(s string) ([]string, bool) {
	sArr := texts.SplitByFirstSpace(s) // Only split out the first one

	ok := len(sArr) > 0

	return sArr, ok
}

func CommandDispatcher(args []string) {
	switch args[0] {
	case USER:
		for HandleLogin() != true {
			// empty
		}
	case EXIT:
		os.Exit(0)
	case ADD:
		HandleAdd(args[1:])
	case REMOVE:
		HandleRemove(args[1:])
	case CLEAR:
		HandleClear()
	case DESTROY:
		HandleDestroy()
	case NEW:
		HandleNewUser()
	case GET:
		HandleLookup(args[1:])
	case SHOW:
		HandleShowDetail(args[1:])
	case LIST:
		HandleList()
	case COPY:
		HandleCopy(args[1:])
	case WHO:
		HandleWho()
	case HELP:
		PrintHelp()
	default:
		texts.Error("Invalid command\n")
		PrintHelp()
	}
}

func HandleLogin() bool {
	var s string
	//var err error
	for {
		texts.Prompt("Login username: ")
		s = texts.GetLineString(os.Stdin)

		if checkIfUserExist(s) {
			break
		} else {
			texts.Error("Username not found\n")
		}
	}
	DataCache.Name = s

	texts.Prompt("Login password: ")
	p := texts.GetLineBytes(os.Stdin)
	DataCache.Password = p

	result := files.ReadData(&DataCache)
	if result {
		texts.Ok("~~~ Logged in as %s ~~~\n", DataCache.Name)
	} else {
		texts.Error("Wrong password\n")
	}

	return result
}

func HandleNewUser() {
	var TempCache entry.UserEntry

	isUsernameValid := false

	for !isUsernameValid {
		texts.Prompt("New username [1-32 chars]: ")
		s := texts.GetLineString(os.Stdin)

		if len(s) < 1 || len(s) > 32 {
			texts.Error("Username length not valid\n")
			continue
		}

		isUsernameValid = true

		for _, name := range Users {
			if name == s {
				isUsernameValid = false
				break
			}
		}

		if isUsernameValid {
			TempCache.Name = s
		} else {
			texts.Warning("Sorry, this name has been used\n")
		}
	}

	for {
		texts.Prompt("New password [AES, must be 16 chars]: ")
		p := texts.GetLineBytes(os.Stdin)
		if len(p) != 16 {
			texts.Error("Password length not valid, must be 16 chars! (AES requirement)\n")
			continue
		}
		TempCache.Password = p
		break
	}

	files.WriteData(&TempCache, true)

	Users = files.ScanAllUsers()
	texts.Ok("New user `%s` created\n", TempCache.Name)
}

func HandleAdd(args []string) {
	var n []byte
	if len(args) > 0 {
		texts.Prompt("Name: ")
		fmt.Printf("%s\n", args[0])
		n = []byte(args[0])
	} else {
		texts.Prompt("Name: ")
		n = texts.GetLineBytes(os.Stdin)
	}

	texts.Prompt("Description: ")
	s := texts.GetLineBytes(os.Stdin)

	texts.Prompt("Password: ")
	p := texts.GetLineBytes(os.Stdin)

	DataCache.Entry.AddEntry(entry.PasswordEntry{n, s, p})

	files.WriteData(&DataCache, true)
	files.ReadData(&DataCache)

	texts.Ok("Entry added\n")
}

func HandleRemove(args []string) {
	if len(args) <= 0 {
		texts.Error("Require name for removal\n")
		PrintHelp()
		return
	}
	DataCache.Entry.RemoveEntry([]byte(args[0]))
	DataCache.WriteData()

	texts.Ok("Removed\n")
}

func HandleClear() {
	if confirm.Ask("Clear all data. Are you sure?") {
		DataCache.Entry.Entries = []entry.PasswordEntry{}
		DataCache.WriteData()
	}

	texts.Ok("Cleared\n")
}

func HandleDestroy() {
	if confirm.Ask("Destroy current user `" + DataCache.Name + "`. Are you sure?") {
		files.DestroyUser(DataCache.Name)
		DataCache = entry.UserEntry{}
		Users = files.ScanAllUsers()

		if len(Users) <= 0 {
			HandleNewUser()
			HandleLogin()
		} else {
			HandleLogin()
		}
	}

	texts.Ok("Destroyed\n")
}

func HandleLookup(args []string) {
	if len(args) <= 0 {
		texts.Error("Require name for lookup\n")
		PrintHelp()
		return
	}

	for _, ent := range DataCache.Entry.Entries {
		if bytes.Equal(ent.Name, []byte(args[0])) {
			texts.Ok("%s\n", ent.Password)
			return
		}
	}
	texts.Error("Not found\n")
}

func HandleShowDetail(args []string) {
	if len(args) <= 0 {
		texts.Error("Require name for lookup\n")
		PrintHelp()
		return
	}

	for _, ent := range DataCache.Entry.Entries {
		if bytes.Equal(ent.Name, []byte(args[0])) {
			formatDetail(ent)
			return
		}
	}
	texts.Error("Not found\n")
}

func HandleList() {
	for _, ent := range DataCache.Entry.Entries {
		fmt.Println("============")
		formatDetail(ent)
	}
	fmt.Println("============")

	fmt.Printf("%v\n", DataCache)
}

func HandleCopy(args []string) {
	if len(args) <= 0 {
		texts.Error("Require name of password to be copied\n")
		PrintHelp()
		return
	}

	for _, ent := range DataCache.Entry.Entries {
		if bytes.Equal(ent.Name, []byte(args[0])) {
			clipboard.WriteAll(string(ent.Password))
			texts.Ok("Password copied to clipboard\n")
			return
		}
	}
	texts.Error("Not found\n")
}

func HandleWho() {
	texts.Ok("%s\n", DataCache.Name)
}

func PrintHelp() {
	fmt.Println(`
	Usage:
		blah...
	`)
}

func formatDetail(e entry.PasswordEntry) {
	fmt.Printf("Name: %s\nDescription: %s\nPassword: %s\n", e.Name, e.Desc, e.Password)
}

func checkIfUserExist(name string) bool {
	for _, user := range Users {
		if user == name {
			return true
		}
	}
	return false
}
