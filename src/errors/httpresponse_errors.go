package bots_errors

import (
	"errors"
	"fmt"
)

func FuncBadTwoCapResponse(taskId int) error {
	errorBody := fmt.Sprintf("%s %s", ColourGrey(fmt.Sprintf("< %d >", taskId)), ColourRed("bad response recieved from 2captcha api"))
	return errors.New(errorBody)
}
