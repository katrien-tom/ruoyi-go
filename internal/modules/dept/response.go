package dept

type DeptResponse struct {
	DeptID    int64          `json:"deptId"`
	ParentID  int64          `json:"parentId"`
	Ancestors string         `json:"ancestors"`
	DeptName  string         `json:"deptName"`
	OrderNum  int            `json:"orderNum"`
	Leader    *string        `json:"leader"`
	Phone     *string        `json:"phone"`
	Email     *string        `json:"email"`
	Status    string         `json:"status"`
	DelFlag   string         `json:"delFlag"`
	Children  []DeptResponse `json:"children"`
}

type DeptTreeSelectResponse struct {
	ID       int64                    `json:"id"`
	Label    string                   `json:"label"`
	Children []DeptTreeSelectResponse `json:"children"`
}
