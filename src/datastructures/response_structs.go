package datastructures

type UsernameVerificationResponse [1]struct {
	Data struct {
		IsUsernameAvailable bool `json:"isUsernameAvailable"`
	} `json:"data"`
	Extensions struct {
		DurationMilliseconds int    `json:"durationMilliseconds"`
		OperationName        string `json:"operationName"`
		RequestID            string `json:"requestID"`
	} `json:"extensions"`
}

type TwoCapResponse struct {
	Status int    `json:"status"`
	CapId  string `json:"request"`
}

type BadSignUpResponse struct {
	CaptchaProof     string   `json:"captcha_proof"`
	Error            string   `json:"error"`
	Errors           []string `json:"errors"`
	ErrorCode        int      `json:"error_code"`
	ErrorDescription string   `json:"error_description"`
}

type GoodSignUpResponse struct {
	CaptchaProof string `json:"captcha_proof"`
	AccessToken  string `json:"access_token"`
	RedirectPath string `json:"redirect_path"`
	UserID       string `json:"userID"`
}
