package confirm

import (
	"fmt"
	"strings"
	"github.com/kevinkassimo/gokey/src/texts"
)

// From https://siongui.github.io/2016/04/23/go-read-yes-no-from-console/
func Ask(q string) bool {
	var s string

	texts.Prompt("%s (y/N): ", q)
	_, err := fmt.Scan(&s)
	if err != nil {
		panic(err)
	}

	s = strings.TrimSpace(s)
	s = strings.ToLower(s)

	if s == "y" || s == "yes" {
		return true
	}
	return false
}
