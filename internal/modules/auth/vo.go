package auth

type CaptchaVO struct {
	CaptchaEnabled bool   `json:"captchaEnabled"`
	Img            string `json:"img"`
	UUID           string `json:"uuid"`
}

type LoginVO struct {
	Token     string `json:"token"`
	TokenType string `json:"tokenType"`
	UserID    int64  `json:"userId"`
	UserName  string `json:"userName"`
	NickName  string `json:"nickName"`
}
