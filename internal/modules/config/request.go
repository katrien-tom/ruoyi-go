package config

type ConfigListRequest struct {
	ConfigName string `json:"configName" form:"configName"`
	ConfigKey  string `json:"configKey" form:"configKey"`
	ConfigType string `json:"configType" form:"configType"`
}

type EditConfigRequest struct {
	ConfigID    int     `json:"configId" binding:"required"`
	ConfigName  string  `json:"configName" binding:"required"`
	ConfigKey   string  `json:"configKey" binding:"required"`
	ConfigValue string  `json:"configValue" binding:"required"`
	ConfigType  string  `json:"configType"`
	Remark      *string `json:"remark"`
}
