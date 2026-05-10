package menu

import (
	"context"
	"fmt"
)

var errMenuHasChildren = fmt.Errorf("menu has children")

type Service struct {
	repo *Repository
}

func NewService(r *Repository) *Service {
	return &Service{repo: r}
}

func (s *Service) FindByUserID(ctx context.Context, userID int64) ([]SysMenu, error) {
	return s.repo.FindByUserID(ctx, userID)
}

func (s *Service) GetMenuList(ctx context.Context) ([]MenuResponse, error) {
	menus, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return buildMenuTree(menus, 0), nil
}

func (s *Service) GetMenuTreeSelect(ctx context.Context) ([]MenuTreeSelectResponse, error) {
	menus, err := s.repo.FindAllByType(ctx, []string{"M", "C"})
	if err != nil {
		return nil, err
	}
	return buildTreeSelect(menus, 0), nil
}

func (s *Service) GetByID(ctx context.Context, menuID int64) (*MenuResponse, error) {
	m, err := s.repo.FindByID(ctx, menuID)
	if err != nil {
		return nil, err
	}
	return &MenuResponse{
		MenuID:    m.MenuID,
		MenuName:  m.MenuName,
		ParentID:  m.ParentID,
		OrderNum:  m.OrderNum,
		Path:      m.Path,
		Component: m.Component,
		Query:     m.Query,
		RouteName: m.RouteName,
		IsFrame:   m.IsFrame,
		IsCache:   m.IsCache,
		MenuType:  m.MenuType,
		Visible:   m.Visible,
		Status:    m.Status,
		Perms:     m.Perms,
		Icon:      m.Icon,
	}, nil
}

func (s *Service) AddMenu(ctx context.Context, req AddMenuRequest) error {
	m := &SysMenu{
		MenuName:  req.MenuName,
		ParentID:  req.ParentID,
		OrderNum:  req.OrderNum,
		Path:      req.Path,
		Component: req.Component,
		Query:     req.Query,
		RouteName: req.RouteName,
		IsFrame:   req.IsFrame,
		IsCache:   req.IsCache,
		MenuType:  req.MenuType,
		Visible:   req.Visible,
		Status:    req.Status,
		Perms:     req.Perms,
		Icon:      req.Icon,
	}
	if m.Status == "" {
		m.Status = "0"
	}
	if m.Visible == "" {
		m.Visible = "0"
	}
	return s.repo.Create(ctx, m)
}

func (s *Service) EditMenu(ctx context.Context, req EditMenuRequest) error {
	m := &SysMenu{
		MenuID:    req.MenuID,
		MenuName:  req.MenuName,
		ParentID:  req.ParentID,
		OrderNum:  req.OrderNum,
		Path:      req.Path,
		Component: req.Component,
		Query:     req.Query,
		RouteName: req.RouteName,
		IsFrame:   req.IsFrame,
		IsCache:   req.IsCache,
		MenuType:  req.MenuType,
		Visible:   req.Visible,
		Status:    req.Status,
		Perms:     req.Perms,
		Icon:      req.Icon,
	}
	return s.repo.Update(ctx, m)
}

func (s *Service) DeleteMenu(ctx context.Context, menuID int64) error {
	has, err := s.repo.HasChildren(ctx, menuID)
	if err != nil {
		return err
	}
	if has {
		return errMenuHasChildren
	}
	return s.repo.Delete(ctx, menuID)
}

func (s *Service) GetRoleMenuTreeSelect(ctx context.Context, roleID int64) (map[string]any, error) {
	menus, err := s.repo.FindAllByType(ctx, []string{"M", "C"})
	if err != nil {
		return nil, err
	}
	tree := buildTreeSelect(menus, 0)

	checkedKeys, err := s.repo.FindRoleMenuIDs(ctx, roleID)
	if err != nil {
		return nil, err
	}

	return map[string]any{
		"menus":       tree,
		"checkedKeys": checkedKeys,
	}, nil
}

func buildMenuTree(menus []SysMenu, parentID int64) []MenuResponse {
	var result []MenuResponse
	for _, m := range menus {
		if m.ParentID != parentID {
			continue
		}
		node := MenuResponse{
			MenuID:    m.MenuID,
			MenuName:  m.MenuName,
			ParentID:  m.ParentID,
			OrderNum:  m.OrderNum,
			Path:      m.Path,
			Component: m.Component,
			Query:     m.Query,
			RouteName: m.RouteName,
			IsFrame:   m.IsFrame,
			IsCache:   m.IsCache,
			MenuType:  m.MenuType,
			Visible:   m.Visible,
			Status:    m.Status,
			Perms:     m.Perms,
			Icon:      m.Icon,
			Children:  buildMenuTree(menus, m.MenuID),
		}
		if node.Children == nil {
			node.Children = []MenuResponse{}
		}
		result = append(result, node)
	}
	if result == nil {
		result = []MenuResponse{}
	}
	return result
}

func buildTreeSelect(menus []SysMenu, parentID int64) []MenuTreeSelectResponse {
	var result []MenuTreeSelectResponse
	for _, m := range menus {
		if m.ParentID != parentID {
			continue
		}
		node := MenuTreeSelectResponse{
			ID:       m.MenuID,
			Label:    m.MenuName,
			Children: buildTreeSelect(menus, m.MenuID),
		}
		if node.Children == nil {
			node.Children = []MenuTreeSelectResponse{}
		}
		result = append(result, node)
	}
	if result == nil {
		result = []MenuTreeSelectResponse{}
	}
	return result
}
