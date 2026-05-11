package post

import "time"

// SysPost maps to the sys_post table.
type SysPost struct {
	PostID     int64      `gorm:"type:bigint(20);primaryKey;autoIncrement;comment:岗位ID;column:post_id;" json:"postId"`
	PostCode   string     `gorm:"type:varchar(64);not null;comment:岗位编码;column:post_code;" json:"postCode"`
	PostName   string     `gorm:"type:varchar(50);not null;comment:岗位名称;column:post_name;" json:"postName"`
	PostSort   int        `gorm:"type:int(4);not null;comment:显示顺序;column:post_sort;" json:"postSort"`
	Status     string     `gorm:"type:char(1);not null;comment:状态（0正常 1停用）;column:status;" json:"status"`
	CreateBy   string     `gorm:"type:varchar(64);default:'';comment:创建者;column:create_by;" json:"createBy"`
	CreateTime *time.Time `gorm:"type:datetime;comment:创建时间;column:create_time;" json:"createTime"`
	UpdateBy   string     `gorm:"type:varchar(64);default:'';comment:更新者;column:update_by;" json:"updateBy"`
	UpdateTime *time.Time `gorm:"type:datetime;comment:更新时间;column:update_time;" json:"updateTime"`
	Remark     *string    `gorm:"type:varchar(500);comment:备注;column:remark;" json:"remark"`
}

func (SysPost) TableName() string {
	return "sys_post"
}
