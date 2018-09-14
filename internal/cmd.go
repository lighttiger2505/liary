package internal

import (
	"fmt"
	"io/ioutil"
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

func GetAppendValue(args []string) (string, error) {
	var val string
	if terminal.IsTerminal(0) {
		if len(args) > 0 {
			val = args[0]
		} else {
			val = ""
		}
	} else {
		b, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return "", fmt.Errorf("Failed make diary file. %s", err.Error())
		}
		val = string(b)
	}
	return val, nil
}
