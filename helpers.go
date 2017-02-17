package migrates

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func wrapName(s string) string {
	return "`" + strings.Replace(s, "`", "\\`", -1) + "`"
}

func wrapNames(s []string) string {
	wrapped := make([]string, len(s))
	for key, value := range s {
		wrapped[key] = wrapName(value)
	}
	return strings.Join(wrapped, ",")
}

// askForConfirmation asks the user for confirmation. A user must type in "yes" or "no" and
// then press enter. It has fuzzy matching, so "y", "Y", "yes", "YES", and "Yes" all count as
// confirmations. If the input is not recognized, it will ask again. The function does not return
// until it gets a valid response from the user.
func askForConfirmation(s string) (bool, error) {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [y/n]: ", s)

		response, err := reader.ReadString('\n')
		if err != nil {
			return false, err
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return true, nil
		} else if response == "n" || response == "no" {
			return false, nil
		}
	}
}
