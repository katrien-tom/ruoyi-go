package notice

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

func (r *Repository) FindAll(ctx context.Context, offset, limit int, noticeTitle, noticeType, createBy string) ([]SysNotice, int64, error) {
	var total int64
	base := r.db.WithContext(ctx).Model(&SysNotice{})

	if noticeTitle != "" {
		base = base.Where("notice_title LIKE ?", "%"+noticeTitle+"%")
	}
	if noticeType != "" {
		base = base.Where("notice_type = ?", noticeType)
	}
	if createBy != "" {
		base = base.Where("create_by LIKE ?", "%"+createBy+"%")
	}

	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var notices []SysNotice
	if err := base.Order("create_time DESC").Offset(offset).Limit(limit).Find(&notices).Error; err != nil {
		return nil, 0, err
	}
	return notices, total, nil
}

func (r *Repository) FindByID(ctx context.Context, id int) (*SysNotice, error) {
	var notice SysNotice
	if err := r.db.WithContext(ctx).Where("notice_id = ?", id).First(&notice).Error; err != nil {
		return nil, err
	}
	return &notice, nil
}

func (r *Repository) Create(ctx context.Context, n *SysNotice) error {
	return r.db.WithContext(ctx).Create(n).Error
}

func (r *Repository) Update(ctx context.Context, n *SysNotice) error {
	return r.db.WithContext(ctx).Model(n).Where("notice_id = ?", n.NoticeID).Updates(n).Error
}

func (r *Repository) Delete(ctx context.Context, ids []int) error {
	return r.db.WithContext(ctx).Where("notice_id IN ?", ids).Delete(&SysNotice{}).Error
}
