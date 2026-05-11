package job

import "time"

// SysJob maps to the sys_job table.
type SysJob struct {
	JobID          int64      `gorm:"type:bigint(20);primaryKey;autoIncrement;comment:任务ID;column:job_id;" json:"jobId"`
	JobName        string     `gorm:"type:varchar(64);default:'';comment:任务名称;column:job_name;" json:"jobName"`
	JobGroup       string     `gorm:"type:varchar(64);default:'DEFAULT';comment:任务组名;column:job_group;" json:"jobGroup"`
	InvokeTarget   string     `gorm:"type:varchar(500);not null;comment:调用目标字符串;column:invoke_target;" json:"invokeTarget"`
	CronExpression string     `gorm:"type:varchar(255);default:'';comment:cron执行表达式;column:cron_expression;" json:"cronExpression"`
	MisfirePolicy  string     `gorm:"type:varchar(20);default:'3';comment:计划执行错误策略;column:misfire_policy;" json:"misfirePolicy"`
	Concurrent     string     `gorm:"type:char(1);default:'1';comment:是否并发执行（0允许 1禁止）;column:concurrent;" json:"concurrent"`
	Status         string     `gorm:"type:char(1);default:'0';comment:状态（0正常 1暂停）;column:status;" json:"status"`
	CreateBy       string     `gorm:"type:varchar(64);default:'';comment:创建者;column:create_by;" json:"createBy"`
	CreateTime     *time.Time `gorm:"type:datetime;comment:创建时间;column:create_time;" json:"createTime"`
	UpdateBy       string     `gorm:"type:varchar(64);default:'';comment:更新者;column:update_by;" json:"updateBy"`
	UpdateTime     *time.Time `gorm:"type:datetime;comment:更新时间;column:update_time;" json:"updateTime"`
	Remark         string     `gorm:"type:varchar(500);default:'';comment:备注信息;column:remark;" json:"remark"`
}

func (SysJob) TableName() string {
	return "sys_job"
}

// SysJobLog maps to the sys_job_log table.
type SysJobLog struct {
	JobLogID      int64      `gorm:"type:bigint(20);primaryKey;autoIncrement;comment:任务日志ID;column:job_log_id;" json:"jobLogId"`
	JobName       string     `gorm:"type:varchar(64);not null;comment:任务名称;column:job_name;" json:"jobName"`
	JobGroup      string     `gorm:"type:varchar(64);not null;comment:任务组名;column:job_group;" json:"jobGroup"`
	InvokeTarget  string     `gorm:"type:varchar(500);not null;comment:调用目标字符串;column:invoke_target;" json:"invokeTarget"`
	JobMessage    *string    `gorm:"type:varchar(500);comment:日志信息;column:job_message;" json:"jobMessage"`
	Status        string     `gorm:"type:char(1);default:'0';comment:执行状态（0正常 1失败）;column:status;" json:"status"`
	ExceptionInfo string     `gorm:"type:varchar(2000);default:'';comment:异常信息;column:exception_info;" json:"exceptionInfo"`
	CreateTime    *time.Time `gorm:"type:datetime;comment:创建时间;column:create_time;" json:"createTime"`
}

func (SysJobLog) TableName() string {
	return "sys_job_log"
}
