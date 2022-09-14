package utils

import (
	"log"
	"os"
	"time"

	colour "github.com/gookit/color"
)

var SHOW_PASS = false
var ColourGreen = colour.FgGreen.Render

func WriteErrorToLogs(unknownError error) {
	file, err := os.OpenFile("userdata/errorlogs.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer file.Close()
	log.SetOutput(file)
	log.Println(unknownError)
}

func TimeOut(delay int) {
	time.Sleep(time.Millisecond * time.Duration(delay))
}
