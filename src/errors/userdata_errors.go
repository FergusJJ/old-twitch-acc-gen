package bots_errors

import (
	"errors"
	"fmt"
)

//used when the login.csv file is empty
var ErrNoEmails = errors.New(ColourRed("no stored emails to login"))

//Should be used if either email or password is missing from a csv row
type ErrIncompleteRow struct {
	EmailName      string
	PasswordEnding string
	UserName       string
	UserPassword   string
	IsEmail        bool
	IsPassword     bool
	IsUserName     bool
	IsUserPassword bool
}

func FuncErrIncompleteUserData(errorStruct ErrIncompleteRow) {
	if !errorStruct.IsEmail {
		errorString := ColourRed(fmt.Sprintf("missing email in user data with password ending in: %s", errorStruct.PasswordEnding))
		err := errors.New(errorString)
		fmt.Println(err)
		return
	}
	if !errorStruct.IsPassword {
		errorString := ColourRed(fmt.Sprintf("missing password in user data with email: %s", errorStruct.EmailName))
		err := errors.New(errorString)
		fmt.Println(err)
		return
	}
	if !errorStruct.IsUserName {
		errorString := ColourRed(fmt.Sprintf("missing username in user data with email: %s", errorStruct.EmailName))
		err := errors.New(errorString)
		fmt.Println(err)
		return
	}
	if !errorStruct.IsUserPassword {
		errorString := ColourRed(fmt.Sprintf("missing userpassword in user data with email: %s", errorStruct.EmailName))
		err := errors.New(errorString)
		fmt.Println(err)
		return
	}
}
