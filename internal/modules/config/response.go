package config

type ConfigResponse struct {
	ConfigID    int     `json:"configId"`
	ConfigName  string  `json:"configName"`
	ConfigKey   string  `json:"configKey"`
	ConfigValue string  `json:"configValue"`
	ConfigType  string  `json:"configType"`
	Remark      *string `json:"remark"`
}

type ConfigListResponse struct {
	Rows  []ConfigResponse `json:"rows"`
	Total int64            `json:"total"`
}
