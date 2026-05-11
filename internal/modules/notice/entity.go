package notice

import "time"

// SysNotice maps to the sys_notice table.
type SysNotice struct {
	NoticeID      int        `gorm:"type:int(4);primaryKey;autoIncrement;comment:公告ID;column:notice_id;" json:"noticeId"`
	NoticeTitle   string     `gorm:"type:varchar(50);not null;comment:公告标题;column:notice_title;" json:"noticeTitle"`
	NoticeType    string     `gorm:"type:char(1);not null;comment:公告类型（1通知 2公告）;column:notice_type;" json:"noticeType"`
	NoticeContent []byte     `gorm:"type:longblob;comment:公告内容;column:notice_content;" json:"noticeContent"`
	Status        string     `gorm:"type:char(1);default:'0';comment:公告状态（0正常 1关闭）;column:status;" json:"status"`
	CreateBy      string     `gorm:"type:varchar(64);default:'';comment:创建者;column:create_by;" json:"createBy"`
	CreateTime    *time.Time `gorm:"type:datetime;comment:创建时间;column:create_time;" json:"createTime"`
	UpdateBy      string     `gorm:"type:varchar(64);default:'';comment:更新者;column:update_by;" json:"updateBy"`
	UpdateTime    *time.Time `gorm:"type:datetime;comment:更新时间;column:update_time;" json:"updateTime"`
	Remark        *string    `gorm:"type:varchar(255);comment:备注;column:remark;" json:"remark"`
}

func (SysNotice) TableName() string {
	return "sys_notice"
}
