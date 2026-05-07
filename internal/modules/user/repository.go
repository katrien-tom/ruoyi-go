package user

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

func (r *Repository) FindByUserName(ctx context.Context, userName string) (*SysUser, error) {
	var user SysUser
	if err := r.db.WithContext(ctx).
		Where("user_name = ? AND del_flag = ?", userName, "0").
		First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) FindPermissionsByUserID(ctx context.Context, userID int64) ([]string, error) {
	var permissions []string

	err := r.db.WithContext(ctx).
		Table("sys_menu m").
		Select("DISTINCT m.perms").
		Joins("JOIN sys_role_menu rm ON rm.menu_id = m.menu_id").
		Joins("JOIN sys_user_role ur ON ur.role_id = rm.role_id").
		Where("ur.user_id = ? AND m.status = ? AND m.perms <> '' AND m.perms IS NOT NULL", userID, "0").
		Scan(&permissions).Error
	if err != nil {
		return nil, err
	}

	return permissions, nil
}

func (r *Repository) FindRoleKeysByUserID(ctx context.Context, userID int64) ([]string, error) {
	var roleKeys []string

	err := r.db.WithContext(ctx).
		Table("sys_role r").
		Select("DISTINCT r.role_key").
		Joins("JOIN sys_user_role ur ON ur.role_id = r.role_id").
		Where("ur.user_id = ? AND r.status = ? AND r.del_flag = ? AND r.role_key <> ''", userID, "0", "0").
		Scan(&roleKeys).Error
	if err != nil {
		return nil, err
	}

	return roleKeys, nil
}
