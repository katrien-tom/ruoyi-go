package loginlog

import "time"

// SysLogininfor maps to the sys_logininfor table.
type SysLogininfor struct {
	InfoID        int64      `gorm:"type:bigint(20);primaryKey;autoIncrement;comment:访问ID;column:info_id;" json:"infoId"`
	UserName      string     `gorm:"type:varchar(50);default:'';comment:用户账号;column:user_name;" json:"userName"`
	IPAddr        string     `gorm:"type:varchar(128);default:'';comment:登录IP地址;column:ipaddr;" json:"ipaddr"`
	LoginLocation string     `gorm:"type:varchar(255);default:'';comment:登录地点;column:login_location;" json:"loginLocation"`
	Browser       string     `gorm:"type:varchar(50);default:'';comment:浏览器类型;column:browser;" json:"browser"`
	OS            string     `gorm:"type:varchar(50);default:'';comment:操作系统;column:os;" json:"os"`
	Status        string     `gorm:"type:char(1);default:'0';comment:登录状态（0成功 1失败）;column:status;" json:"status"`
	Msg           string     `gorm:"type:varchar(255);default:'';comment:提示消息;column:msg;" json:"msg"`
	LoginTime     *time.Time `gorm:"type:datetime;comment:访问时间;column:login_time;" json:"loginTime"`
}

func (SysLogininfor) TableName() string {
	return "sys_logininfor"
}
