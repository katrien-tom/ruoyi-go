package role

import "time"

// SysRole maps to the sys_role table.
type SysRole struct {
	RoleID            int64      `gorm:"type:bigint(20);primaryKey;autoIncrement;comment:角色ID;column:role_id;" json:"roleId"`
	RoleName          string     `gorm:"type:varchar(30);not null;comment:角色名称;column:role_name;" json:"roleName"`
	RoleKey           string     `gorm:"type:varchar(100);not null;comment:角色权限字符串;column:role_key;" json:"roleKey"`
	RoleSort          int        `gorm:"type:int(4);not null;comment:显示顺序;column:role_sort;" json:"roleSort"`
	DataScope         string     `gorm:"type:char(1);default:'1';comment:数据范围;column:data_scope;" json:"dataScope"`
	MenuCheckStrictly int8       `gorm:"type:tinyint(1);default:1;comment:菜单树选择项是否关联显示;column:menu_check_strictly;" json:"menuCheckStrictly"`
	DeptCheckStrictly int8       `gorm:"type:tinyint(1);default:1;comment:部门树选择项是否关联显示;column:dept_check_strictly;" json:"deptCheckStrictly"`
	Status            string     `gorm:"type:char(1);not null;comment:角色状态（0正常 1停用）;column:status;" json:"status"`
	DelFlag           string     `gorm:"type:char(1);default:'0';comment:删除标志（0代表存在 2代表删除）;column:del_flag;" json:"delFlag"`
	CreateBy          string     `gorm:"type:varchar(64);default:'';comment:创建者;column:create_by;" json:"createBy"`
	CreateTime        *time.Time `gorm:"type:datetime;comment:创建时间;column:create_time;" json:"createTime"`
	UpdateBy          string     `gorm:"type:varchar(64);default:'';comment:更新者;column:update_by;" json:"updateBy"`
	UpdateTime        *time.Time `gorm:"type:datetime;comment:更新时间;column:update_time;" json:"updateTime"`
	Remark            *string    `gorm:"type:varchar(500);comment:备注;column:remark;" json:"remark"`
}

func (SysRole) TableName() string {
	return "sys_role"
}

// SysRoleMenu maps to the sys_role_menu table.
type SysRoleMenu struct {
	RoleID  int64 `gorm:"type:bigint(20);primaryKey;comment:角色ID;column:role_id;" json:"roleId"`
	MenuID  int64 `gorm:"type:bigint(20);primaryKey;comment:菜单ID;column:menu_id;" json:"menuId"`
}
