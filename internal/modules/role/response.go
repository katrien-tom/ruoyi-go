package role

type RoleResponse struct {
	RoleID            int64   `json:"roleId"`
	RoleName          string  `json:"roleName"`
	RoleKey           string  `json:"roleKey"`
	RoleSort          int     `json:"roleSort"`
	DataScope         string  `json:"dataScope"`
	MenuCheckStrictly int8    `json:"menuCheckStrictly"`
	DeptCheckStrictly int8    `json:"deptCheckStrictly"`
	Status            string  `json:"status"`
	DelFlag           string  `json:"delFlag"`
	Remark            *string `json:"remark"`
}

type RoleListResponse struct {
	Rows  []RoleResponse `json:"rows"`
	Total int64          `json:"total"`
}

type AuthUserResponse struct {
	UserID      int64  `json:"userId"`
	UserName    string `json:"userName"`
	NickName    string `json:"nickName"`
	Phonenumber string `json:"phonenumber"`
	Status      string `json:"status"`
	DeptID      *int64  `json:"deptId"`
	DeptName    string `json:"deptName"`
}

type AuthUserListResponse struct {
	Rows  []AuthUserResponse `json:"rows"`
	Total int64              `json:"total"`
}
