package config

import "time"

// SysConfig maps to the sys_config table.
type SysConfig struct {
	ConfigID    int        `gorm:"type:int(5);primaryKey;autoIncrement;comment:参数主键;column:config_id;" json:"configId"`
	ConfigName  string     `gorm:"type:varchar(100);default:'';comment:参数名称;column:config_name;" json:"configName"`
	ConfigKey   string     `gorm:"type:varchar(100);default:'';comment:参数键名;column:config_key;" json:"configKey"`
	ConfigValue string     `gorm:"type:varchar(500);default:'';comment:参数键值;column:config_value;" json:"configValue"`
	ConfigType  string     `gorm:"type:char(1);default:'N';comment:系统内置（Y是 N否）;column:config_type;" json:"configType"`
	CreateBy    string     `gorm:"type:varchar(64);default:'';comment:创建者;column:create_by;" json:"createBy"`
	CreateTime  *time.Time `gorm:"type:datetime;comment:创建时间;column:create_time;" json:"createTime"`
	UpdateBy    string     `gorm:"type:varchar(64);default:'';comment:更新者;column:update_by;" json:"updateBy"`
	UpdateTime  *time.Time `gorm:"type:datetime;comment:更新时间;column:update_time;" json:"updateTime"`
	Remark      *string    `gorm:"type:varchar(500);comment:备注;column:remark;" json:"remark"`
}

func (SysConfig) TableName() string {
	return "sys_config"
}
