package operlog

type OperLogListRequest struct {
	Title        string `json:"title" form:"title"`
	BusinessType int    `json:"businessType" form:"businessType"`
	OperName     string `json:"operName" form:"operName"`
	Status       int    `json:"status" form:"status"`
	BeginTime    string `json:"beginTime" form:"beginTime"`
	EndTime      string `json:"endTime" form:"endTime"`
}
