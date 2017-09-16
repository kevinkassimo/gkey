package texts

import(
	"strings"
	"log"
	"bufio"
	"io"
	"os"
	"fmt"
)

func isWhitespace(r rune) bool {
	return r == ' ' || r == '\n' || r == '\t'
}

func SplitBySpace(s string) []string {
	return strings.Fields(s)
}

func SplitByFirstSpace(s string) []string {
	str := strings.Trim(s, "\t\n ")

	isSpaceFound := false
	startSpaceIndex := -1

	for i, r := range str {
		if isWhitespace(r) {
			isSpaceFound = true
			startSpaceIndex = i
			break
		}
	}

	if isSpaceFound {
		for i, r := range str[startSpaceIndex:] {
			if !isWhitespace(r) {
				return []string{str[:startSpaceIndex], str[(startSpaceIndex+i):]}
			}
		}
		log.Fatalf("Split by first space failed: %s\n", str)
		return nil
	} else {
		return []string{str}
	}
}

func GetLineString(rd io.Reader) string {
	reader := bufio.NewReader(rd)
	readStr, err := reader.ReadString('\n')
	if err != nil {
		checkEOF(err)
		panic(err)
	}
	return readStr[:len(readStr)-1]
}

func GetLineStringTrimmed(rd io.Reader) string {
	reader := bufio.NewReader(rd)
	readStr, err := reader.ReadString('\n')
	if err != nil {
		checkEOF(err)
		panic(err)
	}
	readStr = strings.Trim(readStr, "\t ")
	return readStr[:len(readStr)-1]
}

func GetLineBytes(rd io.Reader) []byte {
	reader := bufio.NewReader(rd)
	readStr, err := reader.ReadString('\n')
	if err != nil {
		checkEOF(err)
		panic(err)
	}
	return []byte(readStr[:len(readStr)-1])
}
func GetLineBytesTrimmed(rd io.Reader) []byte {
	reader := bufio.NewReader(rd)
	readStr, err := reader.ReadString('\n')
	if err != nil {
		checkEOF(err)
		panic(err)
	}
	readStr = strings.Trim(readStr, "\t ")
	return []byte(readStr[:len(readStr)-1])
}

func checkEOF(err error) {
	if err.Error() == "EOF" {
		fmt.Printf("\n")
		os.Exit(0)
	}
}