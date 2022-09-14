package main

import (
	"bufio"
	"fmt"
	"os"

	bots_errors "github.com/FergusJJ/old-twitch-acc-gen/src/errors"
	utils "github.com/FergusJJ/old-twitch-acc-gen/src/utils"
)

var TWOCAP_KEY = ""

func main() {
	bots_errors.SetupCloseHandler()
	err := utils.ReadEmailCSV()
	if err != nil {
		fmt.Print(err)
		input := bufio.NewScanner(os.Stdin)
		input.Scan()
		return
	}

	err = utils.LoadTwoCapKey()
	if err != nil {
		fmt.Print(err)
		input := bufio.NewScanner(os.Stdin)
		input.Scan()
		return
	}

	err = utils.BeginSignup()
	if err != nil {
		fmt.Print(err)
		input := bufio.NewScanner(os.Stdin)
		input.Scan()
		return
	}

}
