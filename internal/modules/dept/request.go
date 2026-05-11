package dept

type DeptListRequest struct {
	DeptName string `json:"deptName" form:"deptName"`
	Status   string `json:"status" form:"status"`
}

type AddDeptRequest struct {
	ParentID  int64   `json:"parentId"`
	DeptName  string  `json:"deptName" binding:"required"`
	OrderNum  int     `json:"orderNum"`
	Leader    *string `json:"leader"`
	Phone     *string `json:"phone"`
	Email     *string `json:"email"`
	Status    string  `json:"status"`
}

type EditDeptRequest struct {
	DeptID    int64   `json:"deptId" binding:"required"`
	ParentID  int64   `json:"parentId"`
	DeptName  string  `json:"deptName" binding:"required"`
	OrderNum  int     `json:"orderNum"`
	Leader    *string `json:"leader"`
	Phone     *string `json:"phone"`
	Email     *string `json:"email"`
	Status    string  `json:"status"`
}
