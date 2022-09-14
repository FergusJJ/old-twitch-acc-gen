package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"sync"
	"time"

	datastructures "github.com/FergusJJ/old-twitch-acc-gen/src/datastructures"
	bots_errors "github.com/FergusJJ/old-twitch-acc-gen/src/errors"
)

func CheckExistingUsername(credentials *datastructures.Account, taskId int, client *http.Client) (bool, error) {

	url := "https://gql.twitch.tv/gql"
	postBodyString := fmt.Sprintf(`[{"operationName": "UsernameValidator_User","variables": {"username": "%s"},"extensions": {"persistedQuery": {"version": 1,"sha256Hash": "fd1085cf8350e309b725cf8ca91cd90cac03909a3edeeedbd0872ac912f3d660"}}}]`, credentials.UserName)
	jsonString := []byte(postBodyString)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonString))
	if err != nil {
		WriteErrorToLogs(err)
		return false, bots_errors.FuncErrUnexpected(taskId)
	}
	req.Header = http.Header{"Client-Id": {"kimne78kx3ncx6brgo4mv6wki5h1ko"}}
	resp, err := client.Do(req)
	if err != nil {
		//network err
		WriteErrorToLogs(err)
		return false, bots_errors.FuncErrNetwork(taskId)
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		WriteErrorToLogs(err)
		return false, bots_errors.FuncErrUnexpected(taskId)
	}
	verificationResponse := &datastructures.UsernameVerificationResponse{}
	json.Unmarshal(b, verificationResponse)

	usernameAvailable := verificationResponse[0].Data.IsUsernameAvailable
	return usernameAvailable, nil
}

func BeginSignup() error {
	rand.Seed(time.Now().UnixNano())
	var Wg sync.WaitGroup
	Wg.Add(len(datastructures.AccountsMap))
	if len(datastructures.AccountsMap) == 0 {
		return bots_errors.ErrNoEmails
	}

	for i, account := range datastructures.AccountsMap {

		go func(account *datastructures.Account, i int) {
			fmt.Printf("%s %s", bots_errors.ColourGrey(fmt.Sprintf("< %d >", i)), ColourGreen("starting\n"))

			jar, err := cookiejar.New(nil)
			if err != nil {
				WriteErrorToLogs(err)
				fmt.Println(bots_errors.FuncErrUnexpected(i))
				return
			}
			charlesProxy, err := url.Parse("http://localhost:8888")
			if err != nil {
				WriteErrorToLogs(err)
				fmt.Println(bots_errors.FuncErrUnexpected(i))
				return
			}

			client := &http.Client{
				Jar:       jar,
				Transport: &http.Transport{Proxy: http.ProxyURL(charlesProxy)},
			}

			isUsernameAvailable, err := CheckExistingUsername(account, i, client)
			if err != nil {
				fmt.Println(err)
				return
			}
			if !isUsernameAvailable {
				fmt.Printf("%s %s", bots_errors.ColourGrey(fmt.Sprintf("< %d >", i)), bots_errors.ColourRed(fmt.Sprintf("username: %v is not available\n", account.UserName)))
				return
			}

			err = signup(account, i, client)
			if err != nil {
				//handle error appropriately
				WriteErrorToLogs(err)
				Wg.Done()
			}
			Wg.Done()
		}(account, i)
		//should allow for more reliability when requesting captcha response from twocaptcha api
		TimeOut(2000)
	}
	Wg.Wait()
	return nil
}

func signup(credentials *datastructures.Account, taskId int, client *http.Client) error {
	//create client ->  complete funcaptcha -> send data to login -> verify

	err := getSignupPage(client, taskId)
	if err != nil {
		fmt.Println(err)
		return errors.New("")
	}

	twoCapResponse, err := MakeAPIRequest(taskId)
	if err != nil {
		fmt.Println(err)
		return errors.New("")
	}
	fmt.Printf("%s %s\n", bots_errors.ColourGrey(fmt.Sprintf("< %d >", taskId)), ColourGreen("received funcaptcha token"))
	err = sendSignupForm(twoCapResponse, client, credentials, taskId)
	if err != nil {
		fmt.Println(err)
		return errors.New("")
	}
	//now need to verify twitch account

	return nil
}

func getSignupPage(client *http.Client, taskId int) error {
	req, err := http.NewRequest(http.MethodGet, LOGIN_URL, nil)
	if err != nil {
		WriteErrorToLogs(err)
		return bots_errors.FuncErrUnexpected(taskId)
	}
	req.Header = http.Header{
		"User-Agent": {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.67 Safari/537.36"},
		"Accept":     {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"},
	}

	resp, err := client.Do(req)
	if err != nil {
		//network err, probably timed out request
		WriteErrorToLogs(err)
		return bots_errors.FuncErrNetwork(taskId)
	}

	defer resp.Body.Close()

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		WriteErrorToLogs(err)
		return bots_errors.FuncErrUnexpected(taskId)
	}
	return nil

}

func sendSignupForm(twoCapToken string, client *http.Client, credentials *datastructures.Account, taskId int) error {

	day := rand.Intn(28-1+1) + 1
	month := rand.Intn(12-1+1) + 1
	year := rand.Intn(2004-1990+1) + 1990
	var requestUrl = "https://passport.twitch.tv/register"
	var requestBody = []byte(fmt.Sprintf(`{"username": "%s","password": "%s","email": "%s","birthday": {"day": %d,"month": %d,"year": %d},"client_id": "kimne78kx3ncx6brgo4mv6wki5h1ko","arkose": {"token": "%s"}}`, credentials.UserName, credentials.UserPassword, credentials.Email, day, month, year, twoCapToken))
	req, err := http.NewRequest(http.MethodPost, requestUrl, bytes.NewBuffer(requestBody))
	if err != nil {
		WriteErrorToLogs(err)
		return bots_errors.FuncErrUnexpected(taskId)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		//network err, probably timed out request
		WriteErrorToLogs(err)
		return bots_errors.FuncErrNetwork(taskId)
	}

	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		WriteErrorToLogs(err)
		return bots_errors.FuncErrUnexpected(taskId)
	}

	if resp.StatusCode > 399 {
		signupResponseStruct := &datastructures.BadSignUpResponse{}
		json.Unmarshal(b, signupResponseStruct)
		switch signupResponseStruct.ErrorCode {
		case 2015:
			fmt.Printf("%s %s", bots_errors.ColourGrey(fmt.Sprintf("< %d >", taskId)), bots_errors.ColourRed(fmt.Sprintf("email: %v is already in use\n", credentials.Email)))
			return nil
		case 1000:
			fmt.Printf("%s %s", bots_errors.ColourGrey(fmt.Sprintf("< %d >", taskId)), bots_errors.ColourRed("took too long to get captcha\n"))
			return nil

		default:
			err = fmt.Errorf("%v", signupResponseStruct)
			WriteErrorToLogs(err)
			return bots_errors.FuncUnhandledResp(taskId)
		}
	}
	signupResponseStruct := &datastructures.GoodSignUpResponse{}
	json.Unmarshal(b, signupResponseStruct)
	if SHOW_PASS {
		fmt.Printf("%s %s", bots_errors.ColourGrey(fmt.Sprintf("< %d >", taskId)), ColourGreen(fmt.Sprintf("created new account: %s - %s %s \n", credentials.Email, credentials.UserName, credentials.UserPassword)))
	} else {
		fmt.Printf("%s %s", bots_errors.ColourGrey(fmt.Sprintf("< %d >", taskId)), ColourGreen(fmt.Sprintf("created new account: %s - %s \n", credentials.Email, credentials.UserName)))
	}
	return nil

}
