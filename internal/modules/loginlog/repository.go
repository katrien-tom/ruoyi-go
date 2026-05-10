package loginlog

import (
	"context"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

type loginLogFilter struct {
	UserName  string
	IPAddr    string
	Status    string
	BeginTime string
	EndTime   string
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindAll(ctx context.Context, offset, limit int, f loginLogFilter) ([]SysLogininfor, int64, error) {
	var total int64
	base := r.db.WithContext(ctx).Model(&SysLogininfor{})

	if f.UserName != "" {
		base = base.Where("user_name LIKE ?", "%"+f.UserName+"%")
	}
	if f.IPAddr != "" {
		base = base.Where("ipaddr LIKE ?", "%"+f.IPAddr+"%")
	}
	if f.Status != "" {
		base = base.Where("status = ?", f.Status)
	}
	if f.BeginTime != "" {
		base = base.Where("login_time >= ?", f.BeginTime)
	}
	if f.EndTime != "" {
		base = base.Where("login_time <= ?", f.EndTime)
	}

	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var logs []SysLogininfor
	if err := base.Order("info_id DESC").Offset(offset).Limit(limit).Find(&logs).Error; err != nil {
		return nil, 0, err
	}
	return logs, total, nil
}

func (r *Repository) Delete(ctx context.Context, ids []int64) error {
	return r.db.WithContext(ctx).Where("info_id IN ?", ids).Delete(&SysLogininfor{}).Error
}

func (r *Repository) Clean(ctx context.Context) error {
	return r.db.WithContext(ctx).Where("1 = 1").Delete(&SysLogininfor{}).Error
}

func (r *Repository) Create(ctx context.Context, log *SysLogininfor) error {
	return r.db.WithContext(ctx).Create(log).Error
}
