package dict

type DictTypeListRequest struct {
	DictName string `json:"dictName" form:"dictName"`
	DictType string `json:"dictType" form:"dictType"`
	Status   string `json:"status" form:"status"`
}

type AddDictTypeRequest struct {
	DictName string  `json:"dictName" binding:"required"`
	DictType string  `json:"dictType" binding:"required"`
	Status   string  `json:"status"`
	Remark   *string `json:"remark"`
}

type EditDictTypeRequest struct {
	DictID   int64   `json:"dictId" binding:"required"`
	DictName string  `json:"dictName" binding:"required"`
	DictType string  `json:"dictType" binding:"required"`
	Status   string  `json:"status"`
	Remark   *string `json:"remark"`
}

type DictDataListRequest struct {
	DictType  string `json:"dictType" form:"dictType"`
	DictLabel string `json:"dictLabel" form:"dictLabel"`
	Status    string `json:"status" form:"status"`
}

type AddDictDataRequest struct {
	DictType  string  `json:"dictType" binding:"required"`
	DictLabel string  `json:"dictLabel" binding:"required"`
	DictValue string  `json:"dictValue" binding:"required"`
	DictSort  int     `json:"dictSort"`
	CSSClass  *string `json:"cssClass"`
	ListClass *string `json:"listClass"`
	IsDefault string  `json:"isDefault"`
	Status    string  `json:"status"`
	Remark    *string `json:"remark"`
}

type EditDictDataRequest struct {
	DictCode  int64   `json:"dictCode" binding:"required"`
	DictType  string  `json:"dictType" binding:"required"`
	DictLabel string  `json:"dictLabel" binding:"required"`
	DictValue string  `json:"dictValue" binding:"required"`
	DictSort  int     `json:"dictSort"`
	CSSClass  *string `json:"cssClass"`
	ListClass *string `json:"listClass"`
	IsDefault string  `json:"isDefault"`
	Status    string  `json:"status"`
	Remark    *string `json:"remark"`
}
