package operlog

import "time"

// SysOperLog maps to the sys_oper_log table.
type SysOperLog struct {
	OperID        int64      `gorm:"type:bigint(20);primaryKey;autoIncrement;comment:日志主键;column:oper_id;" json:"operId"`
	Title         string     `gorm:"type:varchar(50);default:'';comment:模块标题;column:title;" json:"title"`
	BusinessType  int        `gorm:"type:int(2);default:0;comment:业务类型（0其它 1新增 2修改 3删除）;column:business_type;" json:"businessType"`
	Method        string     `gorm:"type:varchar(200);default:'';comment:方法名称;column:method;" json:"method"`
	RequestMethod string     `gorm:"type:varchar(10);default:'';comment:请求方式;column:request_method;" json:"requestMethod"`
	OperatorType  int        `gorm:"type:int(1);default:0;comment:操作类别（0其它 1后台用户 2手机端用户）;column:operator_type;" json:"operatorType"`
	OperName      string     `gorm:"type:varchar(50);default:'';comment:操作人员;column:oper_name;" json:"operName"`
	DeptName      string     `gorm:"type:varchar(50);default:'';comment:部门名称;column:dept_name;" json:"deptName"`
	OperURL       string     `gorm:"type:varchar(255);default:'';comment:请求URL;column:oper_url;" json:"operUrl"`
	OperIP        string     `gorm:"type:varchar(128);default:'';comment:主机地址;column:oper_ip;" json:"operIp"`
	OperLocation  string     `gorm:"type:varchar(255);default:'';comment:操作地点;column:oper_location;" json:"operLocation"`
	OperParam     string     `gorm:"type:varchar(2000);default:'';comment:请求参数;column:oper_param;" json:"operParam"`
	JSONResult    string     `gorm:"type:varchar(2000);default:'';comment:返回参数;column:json_result;" json:"jsonResult"`
	Status        int        `gorm:"type:int(1);default:0;comment:操作状态（0正常 1异常）;column:status;" json:"status"`
	ErrorMsg      string     `gorm:"type:varchar(2000);default:'';comment:错误消息;column:error_msg;" json:"errorMsg"`
	OperTime      *time.Time `gorm:"type:datetime;comment:操作时间;column:oper_time;" json:"operTime"`
	CostTime      int64      `gorm:"type:bigint(20);default:0;comment:消耗时间;column:cost_time;" json:"costTime"`
}

func (SysOperLog) TableName() string {
	return "sys_oper_log"
}
