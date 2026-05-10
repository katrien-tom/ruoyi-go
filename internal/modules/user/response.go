package user

import "time"

type UserResponse struct {
	UserID      int64      `json:"userId"`
	DeptID      *int64     `json:"deptId"`
	DeptName    string     `json:"deptName"`
	UserName    string     `json:"userName"`
	NickName    string     `json:"nickName"`
	UserType    string     `json:"userType"`
	Email       string     `json:"email"`
	Phonenumber string     `json:"phonenumber"`
	Sex         string     `json:"sex"`
	Avatar      string     `json:"avatar"`
	Status      string     `json:"status"`
	DelFlag     string     `json:"delFlag"`
	LoginIP     string     `json:"loginIp"`
	LoginDate   *time.Time `json:"loginDate"`
	CreateBy    string     `json:"createBy"`
	CreateTime  *time.Time `json:"createTime"`
	UpdateBy    string     `json:"updateBy"`
	UpdateTime  *time.Time `json:"updateTime"`
	Remark      *string    `json:"remark"`
	Dept        *DeptInfo  `json:"dept"`
	Roles       []RoleInfo `json:"roles"`
	PostIDs     []int64    `json:"postIds"`
}

type DeptInfo struct {
	DeptID   int64  `json:"deptId"`
	DeptName string `json:"deptName"`
}

type RoleInfo struct {
	RoleID   int64  `json:"roleId"`
	RoleName string `json:"roleName"`
	RoleKey  string `json:"roleKey"`
}

type UserListResponse struct {
	Rows  []UserResponse `json:"rows"`
	Total int64          `json:"total"`
}

type AuthRoleResponse struct {
	Roles       []RoleInfo `json:"roles"`
	UserRoleIDs []int64    `json:"userRoleIds"`
}
