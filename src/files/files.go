package files

import (
	"github.com/kevinkassimo/gokey/src/entry"
	"io/ioutil"
	"log"
	"os"
	"os/user"
)

var (
	BASE_DIR string
	USER_DIR string
)

const (
	//BASE_DIR  = "/usr/local/share/pass"
	//USER_DIR  = "/usr/local/share/pass/users"
	DATA_JSON = "data.json"
)

func init() {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	BASE_DIR = usr.HomeDir + "/.gokey_store"
	USER_DIR = usr.HomeDir + "/.gokey_store/users"

	checkDirs()
}

func checkDirs() {
	err := os.MkdirAll(USER_DIR, 0770)
	if err != nil {
		panic(err)
	}
}

func ScanAllUsers() []string {
	folders, err := ioutil.ReadDir(USER_DIR)
	if err != nil {
		log.Fatal(err)
	}

	userList := []string{}

	for _, folder := range folders {
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
