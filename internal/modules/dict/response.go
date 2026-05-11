package dict

type DictTypeResponse struct {
	DictID   int64   `json:"dictId"`
	DictName string  `json:"dictName"`
	DictType string  `json:"dictType"`
	Status   string  `json:"status"`
	Remark   *string `json:"remark"`
}

type DictTypeListResponse struct {
	Rows  []DictTypeResponse `json:"rows"`
	Total int64              `json:"total"`
}

type DictDataResponse struct {
	DictCode  int64   `json:"dictCode"`
	DictSort  int     `json:"dictSort"`
	DictLabel string  `json:"dictLabel"`
	DictValue string  `json:"dictValue"`
	DictType  string  `json:"dictType"`
	CSSClass  *string `json:"cssClass"`
	ListClass *string `json:"listClass"`
	IsDefault string  `json:"isDefault"`
	Status    string  `json:"status"`
	Remark    *string `json:"remark"`
}

type DictDataListResponse struct {
	Rows  []DictDataResponse `json:"rows"`
	Total int64              `json:"total"`
}
