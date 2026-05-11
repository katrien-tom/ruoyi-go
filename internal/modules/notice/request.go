package notice

type NoticeListRequest struct {
	NoticeTitle string `json:"noticeTitle" form:"noticeTitle"`
	NoticeType  string `json:"noticeType" form:"noticeType"`
	CreateBy    string `json:"createBy" form:"createBy"`
}

type AddNoticeRequest struct {
	NoticeTitle   string `json:"noticeTitle" binding:"required"`
	NoticeType    string `json:"noticeType" binding:"required"`
	NoticeContent string `json:"noticeContent"`
	Status        string `json:"status"`
	Remark        *string `json:"remark"`
}

type EditNoticeRequest struct {
	NoticeID      int     `json:"noticeId" binding:"required"`
	NoticeTitle   string  `json:"noticeTitle" binding:"required"`
	NoticeType    string  `json:"noticeType" binding:"required"`
	NoticeContent string  `json:"noticeContent"`
	Status        string  `json:"status"`
	Remark        *string `json:"remark"`
}
