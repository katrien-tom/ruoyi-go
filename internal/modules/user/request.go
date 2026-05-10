package user

type UserListRequest struct {
	PageNum     int    `json:"pageNum" form:"pageNum"`
	PageSize    int    `json:"pageSize" form:"pageSize"`
	UserName    string `json:"userName" form:"userName"`
	Phonenumber string `json:"phonenumber" form:"phonenumber"`
	Status      string `json:"status" form:"status"`
	DeptID      *int64 `json:"deptId" form:"deptId"`
	BeginTime   string `json:"beginTime" form:"beginTime"`
	EndTime     string `json:"endTime" form:"endTime"`
}

type AddUserRequest struct {
	UserName    string  `json:"userName" binding:"required"`
	NickName    string  `json:"nickName" binding:"required"`
	Password    string  `json:"password" binding:"required"`
	DeptID      *int64  `json:"deptId"`
	Email       string  `json:"email"`
	Phonenumber string  `json:"phonenumber"`
	Sex         string  `json:"sex"`
	Status      string  `json:"status"`
	PostIDs     []int64 `json:"postIds"`
	RoleIDs     []int64 `json:"roleIds"`
	Remark      *string `json:"remark"`
}

type EditUserRequest struct {
	UserID      int64   `json:"userId" binding:"required"`
	NickName    string  `json:"nickName" binding:"required"`
	DeptID      *int64  `json:"deptId"`
	Email       string  `json:"email"`
	Phonenumber string  `json:"phonenumber"`
	Sex         string  `json:"sex"`
	Status      string  `json:"status"`
	PostIDs     []int64 `json:"postIds"`
	RoleIDs     []int64 `json:"roleIds"`
	Remark      *string `json:"remark"`
}

type ResetPwdRequest struct {
	UserID   int64  `json:"userId" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ChangeStatusRequest struct {
	UserID int64  `json:"userId" binding:"required"`
	Status string `json:"status" binding:"required"`
}

type AuthRoleRequest struct {
	UserID  int64   `json:"userId" binding:"required"`
	RoleIDs []int64 `json:"roleIds"`
}
