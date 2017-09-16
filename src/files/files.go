package files

import (
	"fmt"
	"github.com/kevinkassimo/gkey/src/entry"
	"io/ioutil"
	"log"
	"os"
)

const (
	BASE_DIR  = "/usr/local/share/pass"
	USER_DIR  = "/usr/local/share/pass/users"
	DATA_JSON = "data.json"
)

func ScanAllUsers() []string {
	folders, err := ioutil.ReadDir(USER_DIR)
	if err != nil {
		log.Fatal(err)
	}

	userList := []string{}

	for _, folder := range folders {
		fmt.Printf("%s\n", folder.Name())
		userList = append(userList, folder.Name())
	}
	return userList
}

func ReadData(e *entry.UserEntry) bool {
	e.FileName = USER_DIR + "/" + e.Name + "/" + DATA_JSON
	return e.ReadData()
}

func WriteData(e *entry.UserEntry, shouldCreate bool) {
	e.FileName = USER_DIR + "/" + e.Name + "/" + DATA_JSON

	if shouldCreate {
		os.Mkdir(USER_DIR+"/"+e.Name, 0770)
	}

	e.WriteData()
}

func DestroyUser(username string) {
	folderName := USER_DIR + "/" + username
	if err := os.RemoveAll(folderName); err != nil {
		log.Fatalf("Cannot remove %s: %s", folderName, err)
	}
}
