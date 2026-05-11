package notice

type NoticeResponse struct {
	NoticeID      int    `json:"noticeId"`
	NoticeTitle   string `json:"noticeTitle"`
	NoticeType    string `json:"noticeType"`
	NoticeContent string `json:"noticeContent"`
	Status        string `json:"status"`
	CreateBy      string `json:"createBy"`
	CreateTime    string `json:"createTime"`
	Remark        *string `json:"remark"`
}

type NoticeListResponse struct {
	Rows  []NoticeResponse `json:"rows"`
	Total int64            `json:"total"`
}
