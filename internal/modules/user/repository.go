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
