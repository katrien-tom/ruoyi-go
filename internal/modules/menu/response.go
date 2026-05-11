package menu

type MenuResponse struct {
	MenuID    int64            `json:"menuId"`
	MenuName  string           `json:"menuName"`
	ParentID  int64            `json:"parentId"`
	OrderNum  int              `json:"orderNum"`
	Path      string           `json:"path"`
	Component *string          `json:"component"`
	Query     *string          `json:"query"`
	RouteName string           `json:"routeName"`
	IsFrame   int              `json:"isFrame"`
	IsCache   int              `json:"isCache"`
	MenuType  string           `json:"menuType"`
	Visible   string           `json:"visible"`
	Status    string           `json:"status"`
	Perms     *string          `json:"perms"`
	Icon      string           `json:"icon"`
	Children  []MenuResponse   `json:"children"`
}

type MenuTreeSelectResponse struct {
	ID       int64                    `json:"id"`
	Label    string                   `json:"label"`
	Children []MenuTreeSelectResponse `json:"children"`
}
