package operlog

import (
	"context"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

type operLogFilter struct {
	Title        string
	BusinessType int
	OperName     string
	Status       int
	BeginTime    string
	EndTime      string
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindAll(ctx context.Context, offset, limit int, f operLogFilter) ([]SysOperLog, int64, error) {
	var total int64
	base := r.db.WithContext(ctx).Model(&SysOperLog{})

	if f.Title != "" {
		base = base.Where("title LIKE ?", "%"+f.Title+"%")
	}
	if f.BusinessType > 0 {
		base = base.Where("business_type = ?", f.BusinessType)
	}
	if f.OperName != "" {
		base = base.Where("oper_name LIKE ?", "%"+f.OperName+"%")
	}
	if f.Status > 0 {
		base = base.Where("status = ?", f.Status)
	}
	if f.BeginTime != "" {
		base = base.Where("oper_time >= ?", f.BeginTime)
	}
	if f.EndTime != "" {
		base = base.Where("oper_time <= ?", f.EndTime)
	}

	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var logs []SysOperLog
	if err := base.Order("oper_id DESC").Offset(offset).Limit(limit).Find(&logs).Error; err != nil {
		return nil, 0, err
	}
	return logs, total, nil
}

func (r *Repository) Delete(ctx context.Context, ids []int64) error {
	return r.db.WithContext(ctx).Where("oper_id IN ?", ids).Delete(&SysOperLog{}).Error
}

func (r *Repository) Clean(ctx context.Context) error {
	return r.db.WithContext(ctx).Where("1 = 1").Delete(&SysOperLog{}).Error
}

func (r *Repository) Create(ctx context.Context, log *SysOperLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}
