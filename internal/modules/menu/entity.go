package menu

// SysMenu maps to the sys_menu table.
type SysMenu struct {
	MenuID    int64   `gorm:"column:menu_id" json:"menuId"`
	MenuName  string  `gorm:"column:menu_name" json:"menuName"`
	ParentID  int64   `gorm:"column:parent_id" json:"parentId"`
	OrderNum  int     `gorm:"column:order_num" json:"orderNum"`
	Path      string  `gorm:"column:path" json:"path"`
	Component *string `gorm:"column:component" json:"component"`
	Query     *string `gorm:"column:query" json:"query"`
	RouteName string  `gorm:"column:route_name" json:"routeName"`
	IsFrame   int     `gorm:"column:is_frame" json:"isFrame"`
	IsCache   int     `gorm:"column:is_cache" json:"isCache"`
	MenuType  string  `gorm:"column:menu_type" json:"menuType"`
	Visible   string  `gorm:"column:visible" json:"visible"`
	Status    string  `gorm:"column:status" json:"status"`
	Perms     *string `gorm:"column:perms" json:"perms"`
	Icon      string  `gorm:"column:icon" json:"icon"`
}

func (SysMenu) TableName() string {
	return "sys_menu"
}
