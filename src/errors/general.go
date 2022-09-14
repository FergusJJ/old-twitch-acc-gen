package bots_errors

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	colour "github.com/gookit/color"
)

var ColourRed = colour.FgRed.Render
var ColourGrey = colour.FgDarkGray.Render

//use this before writing the error to logs
func FuncErrUnexpected(taskId int) error {
	errorBody := fmt.Sprintf("%s %s", ColourGrey(fmt.Sprintf("< %d >", taskId)), ColourRed("there was an unexpected error, check logs for details"))
	return errors.New(errorBody)
}

func FuncErrNetwork(taskId int) error {
	errorBody := fmt.Sprintf("%s %s", ColourGrey(fmt.Sprintf("< %d >", taskId)), ColourRed("suspected network error, check logs for details"))
	return errors.New(errorBody)
}

func FuncUnhandledResp(taskId int) error {
	errorBody := fmt.Sprintf("%s %s", ColourGrey(fmt.Sprintf("< %d >", taskId)), ColourRed("unhandled register response code, check logs for details"))
	return errors.New(errorBody)
}

//https://golangcode.com/handle-ctrl-c-exit-in-terminal/
func SetupCloseHandler() {
	//sets up a goroutine to monitor for a ctrl-c interrupt
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		os.Exit(0)
	}()
}
