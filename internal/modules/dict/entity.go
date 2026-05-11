package dict

import "time"

// SysDictType maps to the sys_dict_type table.
type SysDictType struct {
	DictID     int64      `gorm:"type:bigint(20);primaryKey;autoIncrement;comment:字典主键;column:dict_id;" json:"dictId"`
	DictName   string     `gorm:"type:varchar(100);default:'';comment:字典名称;column:dict_name;" json:"dictName"`
	DictType   string     `gorm:"type:varchar(100);default:'';comment:字典类型;column:dict_type;" json:"dictType"`
	Status     string     `gorm:"type:char(1);default:'0';comment:状态（0正常 1停用）;column:status;" json:"status"`
	CreateBy   string     `gorm:"type:varchar(64);default:'';comment:创建者;column:create_by;" json:"createBy"`
	CreateTime *time.Time `gorm:"type:datetime;comment:创建时间;column:create_time;" json:"createTime"`
	UpdateBy   string     `gorm:"type:varchar(64);default:'';comment:更新者;column:update_by;" json:"updateBy"`
	UpdateTime *time.Time `gorm:"type:datetime;comment:更新时间;column:update_time;" json:"updateTime"`
	Remark     *string    `gorm:"type:varchar(500);comment:备注;column:remark;" json:"remark"`
}

func (SysDictType) TableName() string {
	return "sys_dict_type"
}

// SysDictData maps to the sys_dict_data table.
type SysDictData struct {
	DictCode   int64      `gorm:"type:bigint(20);primaryKey;autoIncrement;comment:字典编码;column:dict_code;" json:"dictCode"`
	DictSort   int        `gorm:"type:int(4);default:0;comment:字典排序;column:dict_sort;" json:"dictSort"`
	DictLabel  string     `gorm:"type:varchar(100);default:'';comment:字典标签;column:dict_label;" json:"dictLabel"`
	DictValue  string     `gorm:"type:varchar(100);default:'';comment:字典键值;column:dict_value;" json:"dictValue"`
	DictType   string     `gorm:"type:varchar(100);default:'';comment:字典类型;column:dict_type;" json:"dictType"`
	CSSClass   *string    `gorm:"type:varchar(100);comment:样式属性;column:css_class;" json:"cssClass"`
	ListClass  *string    `gorm:"type:varchar(100);comment:表格回显样式;column:list_class;" json:"listClass"`
	IsDefault  string     `gorm:"type:char(1);default:'N';comment:是否默认（Y是 N否）;column:is_default;" json:"isDefault"`
	Status     string     `gorm:"type:char(1);default:'0';comment:状态（0正常 1停用）;column:status;" json:"status"`
	CreateBy   string     `gorm:"type:varchar(64);default:'';comment:创建者;column:create_by;" json:"createBy"`
	CreateTime *time.Time `gorm:"type:datetime;comment:创建时间;column:create_time;" json:"createTime"`
	UpdateBy   string     `gorm:"type:varchar(64);default:'';comment:更新者;column:update_by;" json:"updateBy"`
	UpdateTime *time.Time `gorm:"type:datetime;comment:更新时间;column:update_time;" json:"updateTime"`
	Remark     *string    `gorm:"type:varchar(500);comment:备注;column:remark;" json:"remark"`
}

func (SysDictData) TableName() string {
	return "sys_dict_data"
}
