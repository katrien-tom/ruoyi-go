package menu

type AddMenuRequest struct {
	MenuName  string  `json:"menuName" binding:"required"`
	ParentID  int64   `json:"parentId"`
	OrderNum  int     `json:"orderNum"`
	Path      string  `json:"path"`
	Component *string `json:"component"`
	Query     *string `json:"query"`
	RouteName string  `json:"routeName"`
	IsFrame   int     `json:"isFrame"`
	IsCache   int     `json:"isCache"`
	MenuType  string  `json:"menuType" binding:"required"`
	Visible   string  `json:"visible"`
	Status    string  `json:"status"`
	Perms     *string `json:"perms"`
	Icon      string  `json:"icon"`
}

type EditMenuRequest struct {
	MenuID    int64   `json:"menuId" binding:"required"`
	MenuName  string  `json:"menuName" binding:"required"`
	ParentID  int64   `json:"parentId"`
	OrderNum  int     `json:"orderNum"`
	Path      string  `json:"path"`
	Component *string `json:"component"`
	Query     *string `json:"query"`
	RouteName string  `json:"routeName"`
	IsFrame   int     `json:"isFrame"`
	IsCache   int     `json:"isCache"`
	MenuType  string  `json:"menuType" binding:"required"`
	Visible   string  `json:"visible"`
	Status    string  `json:"status"`
	Perms     *string `json:"perms"`
	Icon      string  `json:"icon"`
}
