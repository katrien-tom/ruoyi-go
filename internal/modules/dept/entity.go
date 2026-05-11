package dept

import "time"

// SysDept maps to the sys_dept table.
type SysDept struct {
	DeptID     int64      `gorm:"type:bigint(20);primaryKey;autoIncrement;comment:部门id;column:dept_id;" json:"deptId"`
	ParentID   int64      `gorm:"type:bigint(20);default:0;comment:父部门id;column:parent_id;" json:"parentId"`
	Ancestors  string     `gorm:"type:varchar(50);default:'';comment:祖级列表;column:ancestors;" json:"ancestors"`
	DeptName   string     `gorm:"type:varchar(30);default:'';comment:部门名称;column:dept_name;" json:"deptName"`
	OrderNum   int        `gorm:"type:int(4);default:0;comment:显示顺序;column:order_num;" json:"orderNum"`
	Leader     *string    `gorm:"type:varchar(20);comment:负责人;column:leader;" json:"leader"`
	Phone      *string    `gorm:"type:varchar(11);comment:联系电话;column:phone;" json:"phone"`
	Email      *string    `gorm:"type:varchar(50);comment:邮箱;column:email;" json:"email"`
	Status     string     `gorm:"type:char(1);default:'0';comment:部门状态（0正常 1停用）;column:status;" json:"status"`
	DelFlag    string     `gorm:"type:char(1);default:'0';comment:删除标志（0代表存在 2代表删除）;column:del_flag;" json:"delFlag"`
	CreateBy   string     `gorm:"type:varchar(64);default:'';comment:创建者;column:create_by;" json:"createBy"`
	CreateTime *time.Time `gorm:"type:datetime;comment:创建时间;column:create_time;" json:"createTime"`
	UpdateBy   string     `gorm:"type:varchar(64);default:'';comment:更新者;column:update_by;" json:"updateBy"`
	UpdateTime *time.Time `gorm:"type:datetime;comment:更新时间;column:update_time;" json:"updateTime"`
}

func (SysDept) TableName() string {
	return "sys_dept"
}
