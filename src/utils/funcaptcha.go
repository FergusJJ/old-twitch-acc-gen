package utils

import (
	"encoding/json"
	"fmt"

	"github.com/FergusJJ/old-twitch-acc-gen/src/datastructures"
	bots_errors "github.com/FergusJJ/old-twitch-acc-gen/src/errors"
	"github.com/valyala/fasthttp"
)

var TWITCH_PK = "E5554D43-23CC-1982-971D-6A2262A2CA24"
var TWITCH_SURL = "https://twitch-api.arkoselabs.com"
var LOGIN_URL = "https://www.twitch.tv/signup"
var TWOCAP_KEY = ""

func MakeAPIRequest(taskId int) (token string, err error) {
	twocapResponse := &datastructures.TwoCapResponse{}
	funcapResponse := &datastructures.TwoCapResponse{}
	requestClient := &fasthttp.Client{Name: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.67 Safari/537.36"}
	requestURL := fmt.Sprintf("http://2captcha.com/in.php?key=%s&method=funcaptcha&publickey=%s&surl=%s&pageurl=%s&json=1", TWOCAP_KEY, TWITCH_PK, TWITCH_SURL, LOGIN_URL)
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(requestURL)
	resp := fasthttp.AcquireResponse()
	err = requestClient.Do(req, resp)
	if err != nil {
		WriteErrorToLogs(err)
		return "", bots_errors.FuncErrNetwork(taskId)
	}
	json.Unmarshal(resp.Body(), &twocapResponse)
	if twocapResponse.Status == 0 {
		return "", bots_errors.FuncBadTwoCapResponse(taskId)
	}
	//now need to get token from twocaptcha
	captchaNotReady := true
	for captchaNotReady {
		requestURL = fmt.Sprintf("http://2captcha.com/res.php?key=%s&action=get&id=%s&json=1", TWOCAP_KEY, twocapResponse.CapId)
		req = fasthttp.AcquireRequest()
		req.SetRequestURI(requestURL)
		resp = fasthttp.AcquireResponse()

		TimeOut(4000)

		err = requestClient.Do(req, resp)
		if err != nil {
			WriteErrorToLogs(err)
			return "", bots_errors.FuncErrNetwork(taskId)
		}
		json.Unmarshal(resp.Body(), &funcapResponse)

		if funcapResponse.Status == 1 {
			captchaNotReady = false
			return funcapResponse.CapId, nil
		}
	}
	return "", nil

}
