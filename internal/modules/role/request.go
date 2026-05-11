package role

type RoleListRequest struct {
	RoleName string `json:"roleName" form:"roleName"`
	RoleKey  string `json:"roleKey" form:"roleKey"`
	Status   string `json:"status" form:"status"`
}

type AddRoleRequest struct {
	RoleName          string  `json:"roleName" binding:"required"`
	RoleKey           string  `json:"roleKey" binding:"required"`
	RoleSort          int     `json:"roleSort"`
	DataScope         string  `json:"dataScope"`
	MenuCheckStrictly int8    `json:"menuCheckStrictly"`
	DeptCheckStrictly int8    `json:"deptCheckStrictly"`
	Status            string  `json:"status"`
	MenuIDs           []int64 `json:"menuIds"`
	Remark            *string `json:"remark"`
}

type EditRoleRequest struct {
	RoleID            int64   `json:"roleId" binding:"required"`
	RoleName          string  `json:"roleName" binding:"required"`
	RoleKey           string  `json:"roleKey" binding:"required"`
	RoleSort          int     `json:"roleSort"`
	DataScope         string  `json:"dataScope"`
	MenuCheckStrictly int8    `json:"menuCheckStrictly"`
	DeptCheckStrictly int8    `json:"deptCheckStrictly"`
	Status            string  `json:"status"`
	MenuIDs           []int64 `json:"menuIds"`
	Remark            *string `json:"remark"`
}

type ChangeRoleStatusRequest struct {
	RoleID int64  `json:"roleId" binding:"required"`
	Status string `json:"status" binding:"required"`
}

type DataScopeRequest struct {
	RoleID  int64   `json:"roleId" binding:"required"`
	DataScope string `json:"dataScope" binding:"required"`
	DeptIDs  []int64 `json:"deptIds"`
}

type AuthUserListRequest struct {
	RoleID   int64  `json:"roleId" form:"roleId"`
	UserName string `json:"userName" form:"userName"`
	Phonenumber string `json:"phonenumber" form:"phonenumber"`
}

type AuthUserCancelRequest struct {
	RoleID int64  `json:"roleId" binding:"required"`
	UserID int64  `json:"userId" binding:"required"`
}

type AuthUserCancelAllRequest struct {
	RoleID int64   `json:"roleId" binding:"required"`
	UserIDs []int64 `json:"userIds" binding:"required"`
}

type AuthUserSelectAllRequest struct {
	RoleID int64   `json:"roleId" binding:"required"`
	UserIDs []int64 `json:"userIds" binding:"required"`
}
