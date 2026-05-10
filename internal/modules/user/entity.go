package user

import "time"

// SysUser maps to the sys_user table.
type SysUser struct {
	UserID        int64      `gorm:"type:bigint(20);primaryKey;autoIncrement;comment:用户ID;column:user_id;" json:"userId"`
	DeptID        *int64     `gorm:"type:bigint(20);comment:部门ID;column:dept_id;" json:"deptId"`
	UserName      string     `gorm:"type:varchar(30);not null;comment:用户账号;column:user_name;" json:"userName"`
	NickName      string     `gorm:"type:varchar(30);not null;comment:用户昵称;column:nick_name;" json:"nickName"`
	UserType      string     `gorm:"type:varchar(2);default:'00';comment:用户类型（00系统用户）;column:user_type;" json:"userType"`
	Email         string     `gorm:"type:varchar(50);default:'';comment:用户邮箱;column:email;" json:"email"`
	Phonenumber   string     `gorm:"type:varchar(11);default:'';comment:手机号码;column:phonenumber;" json:"phonenumber"`
	Sex           string     `gorm:"type:char(1);default:'0';comment:用户性别（0男 1女 2未知）;column:sex;" json:"sex"`
	Avatar        string     `gorm:"type:varchar(100);default:'';comment:头像地址;column:avatar;" json:"avatar"`
	Password      string     `gorm:"type:varchar(100);default:'';comment:密码;column:password;" json:"password"`
	Status        string     `gorm:"type:char(1);default:'0';comment:账号状态（0正常 1停用）;column:status;" json:"status"`
	DelFlag       string     `gorm:"type:char(1);default:'0';comment:删除标志（0代表存在 2代表删除）;column:del_flag;" json:"delFlag"`
	LoginIP       string     `gorm:"type:varchar(128);default:'';comment:最后登录IP;column:login_ip;" json:"loginIp"`
	LoginDate     *time.Time `gorm:"type:datetime;comment:最后登录时间;column:login_date;" json:"loginDate"`
	PwdUpdateDate *time.Time `gorm:"type:datetime;comment:密码最后更新时间;column:pwd_update_date;" json:"pwdUpdateDate"`
	CreateBy      string     `gorm:"type:varchar(64);default:'';comment:创建者;column:create_by;" json:"createBy"`
	CreateTime    *time.Time `gorm:"type:datetime;comment:创建时间;column:create_time;" json:"createTime"`
	UpdateBy      string     `gorm:"type:varchar(64);default:'';comment:更新者;column:update_by;" json:"updateBy"`
	UpdateTime    *time.Time `gorm:"type:datetime;comment:更新时间;column:update_time;" json:"updateTime"`
	Remark        *string    `gorm:"type:varchar(500);comment:备注;column:remark;" json:"remark"`
}

func (SysUser) TableName() string {
	return "sys_user"
}

// SysUserRole maps to the sys_user_role table.
type SysUserRole struct {
	UserID int64 `gorm:"type:bigint(20);primaryKey;comment:用户ID;column:user_id;" json:"userId"`
	RoleID int64 `gorm:"type:bigint(20);primaryKey;comment:角色ID;column:role_id;" json:"roleId"`
}

func (SysUserRole) TableName() string {
	return "sys_user_role"
}

// SysUserPost maps to the sys_user_post table.
type SysUserPost struct {
	UserID int64 `gorm:"type:bigint(20);primaryKey;comment:用户ID;column:user_id;" json:"userId"`
	PostID int64 `gorm:"type:bigint(20);primaryKey;comment:岗位ID;column:post_id;" json:"postId"`
}

func (SysUserPost) TableName() string {
	return "sys_user_post"
}
