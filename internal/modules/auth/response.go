package auth

import "github.com/banyejiu/ruoyi-go/internal/security"

type CaptchaResponse struct {
	CaptchaEnabled bool   `json:"captchaEnabled"`
	Img            string `json:"img"`
	UUID           string `json:"uuid"`
}

type LoginResponse struct {
	Token     string `json:"token"`
	TokenType string `json:"tokenType"`
	UserID    int64  `json:"userId"`
	UserName  string `json:"userName"`
	NickName  string `json:"nickName"`
	DeptID    *int64 `json:"deptId,omitempty"`
	ExpiresAt int64  `json:"expiresAt"`
}

type GetInfoResponse struct {
	User        *security.UserInfo `json:"user"`
	Roles       []string           `json:"roles"`
	Permissions []string           `json:"permissions"`
}

type OnlineUserResponse struct {
	TokenID    string `json:"tokenId"`
	UserID     int64  `json:"userId"`
	UserName   string `json:"userName"`
	NickName   string `json:"nickName"`
	DeptID     *int64 `json:"deptId,omitempty"`
	IPAddr     string `json:"ipaddr,omitempty"`
	LoginIP    string `json:"loginIp,omitempty"`
	Browser    string `json:"browser,omitempty"`
	OS         string `json:"os,omitempty"`
	LoginTime  int64  `json:"loginTime"`
	ExpireTime int64  `json:"expireTime"`
}
