package util

import (
	"fmt"

	"golang.org/x/term"
)

func PromptForToken(prompt string, target *string) error {
	fmt.Printf("%s: ", prompt)

	token, err := term.ReadPassword(0)
	fmt.Println()
	if err != nil {
		return err
	}

	*target = string(token)

	return nil
}
