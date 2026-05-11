package menu

import (
	"context"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindByUserID(ctx context.Context, userID int64) ([]SysMenu, error) {
	var menus []SysMenu
	err := r.db.WithContext(ctx).
		Table("sys_menu m").
		Select("DISTINCT m.*").
		Joins("JOIN sys_role_menu rm ON rm.menu_id = m.menu_id").
		Joins("JOIN sys_user_role ur ON ur.role_id = rm.role_id").
		Where("ur.user_id = ? AND m.menu_type IN ? AND m.status = ?", userID, []string{"M", "C"}, "0").
		Order("m.parent_id, m.order_num").
		Scan(&menus).Error
	if err != nil {
		return nil, err
	}
	return menus, nil
}

func (r *Repository) FindAll(ctx context.Context) ([]SysMenu, error) {
	var menus []SysMenu
	err := r.db.WithContext(ctx).Order("parent_id, order_num").Find(&menus).Error
	return menus, err
}

func (r *Repository) FindAllByType(ctx context.Context, menuTypes []string) ([]SysMenu, error) {
	var menus []SysMenu
	err := r.db.WithContext(ctx).
		Where("menu_type IN ? AND status = ?", menuTypes, "0").
		Order("parent_id, order_num").
		Find(&menus).Error
	return menus, err
}

func (r *Repository) FindByUserIDAndTypes(ctx context.Context, userID int64, menuTypes []string) ([]SysMenu, error) {
	var menus []SysMenu
	err := r.db.WithContext(ctx).
		Table("sys_menu m").
		Select("DISTINCT m.*").
		Joins("JOIN sys_role_menu rm ON rm.menu_id = m.menu_id").
		Joins("JOIN sys_user_role ur ON ur.role_id = rm.role_id").
		Where("ur.user_id = ? AND m.menu_type IN ? AND m.status = ?", userID, menuTypes, "0").
		Order("m.parent_id, m.order_num").
		Scan(&menus).Error
	if err != nil {
		return nil, err
	}
	return menus, nil
}

func (r *Repository) FindByID(ctx context.Context, menuID int64) (*SysMenu, error) {
	var m SysMenu
	if err := r.db.WithContext(ctx).Where("menu_id = ?", menuID).First(&m).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *Repository) Create(ctx context.Context, m *SysMenu) error {
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *Repository) Update(ctx context.Context, m *SysMenu) error {
	return r.db.WithContext(ctx).Model(m).Where("menu_id = ?", m.MenuID).Updates(m).Error
}

func (r *Repository) Delete(ctx context.Context, menuID int64) error {
	return r.db.WithContext(ctx).Where("menu_id = ?", menuID).Delete(&SysMenu{}).Error
}

func (r *Repository) HasChildren(ctx context.Context, menuID int64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&SysMenu{}).Where("parent_id = ?", menuID).Count(&count).Error
	return count > 0, err
}

func (r *Repository) FindRoleMenuIDs(ctx context.Context, roleID int64) ([]int64, error) {
	var ids []int64
	err := r.db.WithContext(ctx).
		Table("sys_role_menu").
		Select("menu_id").
		Where("role_id = ?", roleID).
		Scan(&ids).Error
	return ids, err
}

func (r *Repository) ReplaceRoleMenus(ctx context.Context, roleID int64, menuIDs []int64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("role_id = ?", roleID).Delete(&struct {
			RoleID int64 `gorm:"column:role_id"`
			MenuID int64 `gorm:"column:menu_id"`
		}{}).Error; err != nil {
			return err
		}
		if len(menuIDs) == 0 {
			return nil
		}
		type sysRoleMenu struct {
			RoleID int64 `gorm:"column:role_id"`
			MenuID int64 `gorm:"column:menu_id"`
		}
		rows := make([]sysRoleMenu, len(menuIDs))
		for i, mid := range menuIDs {
			rows[i] = sysRoleMenu{RoleID: roleID, MenuID: mid}
		}
		return tx.Table("sys_role_menu").Create(&rows).Error
	})
}

func (r *Repository) FindMenuPermsByUserID(ctx context.Context, userID int64) ([]string, error) {
	var perms []string
	err := r.db.WithContext(ctx).
		Table("sys_menu m").
		Select("DISTINCT m.perms").
		Joins("JOIN sys_role_menu rm ON rm.menu_id = m.menu_id").
		Joins("JOIN sys_user_role ur ON ur.role_id = rm.role_id").
		Where("ur.user_id = ? AND m.status = ? AND m.perms <> '' AND m.perms IS NOT NULL", userID, "0").
		Scan(&perms).Error
	return perms, err
}

// TableName for sys_role_menu raw operations
func (sysRoleMenu) TableName() string { return "sys_role_menu" }

type sysRoleMenu struct {
	RoleID int64 `gorm:"column:role_id;primaryKey"`
	MenuID int64 `gorm:"column:menu_id;primaryKey"`
}
