package utils

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"

	datastructures "github.com/FergusJJ/old-twitch-acc-gen/src/datastructures"
	bots_errors "github.com/FergusJJ/old-twitch-acc-gen/src/errors"
)

func ReadEmailCSV() error {
	csvFile, err := os.Open("userdata/accounts.csv")
	if err != nil {
		WriteErrorToLogs(err)
		return bots_errors.FuncErrUnexpected(-1)
	}
	defer csvFile.Close()

	csvFileLines, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		WriteErrorToLogs(err)
		return bots_errors.FuncErrUnexpected(-1)
	}
	if len(csvFileLines) == 1 {
		return bots_errors.ErrNoEmails
	}
	for i, line := range csvFileLines {
		if i == 0 {
			continue
		}
		if len(line[0]) == 0 || len(line[1]) == 0 || len(line[2]) == 0 || len(line[3]) == 0 {
			err := &bots_errors.ErrIncompleteRow{}

			if len(line[0]) == 0 {
				err.EmailName = ""
				err.IsEmail = false
			} else {
				err.EmailName = line[0]
				err.IsEmail = true
			}

			if len(line[1]) == 0 {
				err.PasswordEnding = ""
				err.IsPassword = false
			} else {
				err.PasswordEnding = line[1][len(line[1])-4:]
				err.IsPassword = true
			}

			if len(line[2]) == 0 {
				err.UserName = ""
				err.IsUserName = false
			} else {
				err.UserName = line[2]
				err.IsUserName = true
			}
			if len(line[3]) == 0 {
				err.UserPassword = ""
				err.IsUserPassword = false
			} else {
				err.UserPassword = line[3][len(line[3])-4:]
				err.IsUserPassword = true
			}

			bots_errors.FuncErrIncompleteUserData(*err)
			return errors.New("")
		}

		emailStruct := &datastructures.Account{
			Email:        line[0],
			Password:     line[1],
			UserName:     line[2],
			UserPassword: line[3],
		}
		datastructures.AccountsMap[i] = emailStruct
	}

	return nil
}

func LoadTwoCapKey() error {
	type jsonSettings struct {
		TwoCapKey    string `json:"2captcha"`
		ShowPassword bool   `json:"showPassword"`
	}
	var userSettings jsonSettings
	jsonFile, err := os.Open("userdata/settings.json")
	if err != nil {
		WriteErrorToLogs(err)
		return bots_errors.FuncErrUnexpected(-1)
	}
	defer jsonFile.Close()
	bytes, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(bytes, &userSettings)
	TWOCAP_KEY = userSettings.TwoCapKey
	SHOW_PASS = userSettings.ShowPassword
	return nil
}
