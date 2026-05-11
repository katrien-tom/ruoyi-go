package loginlog

type LoginLogListRequest struct {
	UserName   string `json:"userName" form:"userName"`
	IPAddr     string `json:"ipaddr" form:"ipaddr"`
	Status     string `json:"status" form:"status"`
	BeginTime  string `json:"beginTime" form:"beginTime"`
	EndTime    string `json:"endTime" form:"endTime"`
}
