package operlog

import "time"

type OperLogResponse struct {
	OperID        int64      `json:"operId"`
	Title         string     `json:"title"`
	BusinessType  int        `json:"businessType"`
	Method        string     `json:"method"`
	RequestMethod string     `json:"requestMethod"`
	OperatorType  int        `json:"operatorType"`
	OperName      string     `json:"operName"`
	DeptName      string     `json:"deptName"`
	OperURL       string     `json:"operUrl"`
	OperIP        string     `json:"operIp"`
	OperLocation  string     `json:"operLocation"`
	OperParam     string     `json:"operParam"`
	JSONResult    string     `json:"jsonResult"`
	Status        int        `json:"status"`
	ErrorMsg      string     `json:"errorMsg"`
	OperTime      *time.Time `json:"operTime"`
	CostTime      int64      `json:"costTime"`
}

type OperLogListResponse struct {
	Rows  []OperLogResponse `json:"rows"`
	Total int64             `json:"total"`
}
