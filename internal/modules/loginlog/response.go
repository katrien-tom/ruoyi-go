package loginlog

import "time"

type LoginLogResponse struct {
	InfoID        int64      `json:"infoId"`
	UserName      string     `json:"userName"`
	IPAddr        string     `json:"ipaddr"`
	LoginLocation string     `json:"loginLocation"`
	Browser       string     `json:"browser"`
	OS            string     `json:"os"`
	Status        string     `json:"status"`
	Msg           string     `json:"msg"`
	LoginTime     *time.Time `json:"loginTime"`
}

type LoginLogListResponse struct {
	Rows  []LoginLogResponse `json:"rows"`
	Total int64              `json:"total"`
}
