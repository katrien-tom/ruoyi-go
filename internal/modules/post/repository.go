package post

import (
	"context"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

type postListFilter struct {
	PostCode string
	PostName string
	Status   string
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindAll(ctx context.Context, offset, limit int, filter postListFilter) ([]SysPost, int64, error) {
	var total int64
	base := r.db.WithContext(ctx).Model(&SysPost{})

	if filter.PostCode != "" {
		base = base.Where("post_code LIKE ?", "%"+filter.PostCode+"%")
	}
	if filter.PostName != "" {
		base = base.Where("post_name LIKE ?", "%"+filter.PostName+"%")
	}
	if filter.Status != "" {
		base = base.Where("status = ?", filter.Status)
	}

	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var posts []SysPost
	if err := base.Order("post_sort").Offset(offset).Limit(limit).Find(&posts).Error; err != nil {
		return nil, 0, err
	}
	return posts, total, nil
}

func (r *Repository) FindByID(ctx context.Context, id int64) (*SysPost, error) {
	var post SysPost
	if err := r.db.WithContext(ctx).Where("post_id = ?", id).First(&post).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *Repository) Create(ctx context.Context, p *SysPost) error {
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *Repository) Update(ctx context.Context, p *SysPost) error {
	return r.db.WithContext(ctx).Model(p).Where("post_id = ?", p.PostID).Updates(p).Error
}

func (r *Repository) Delete(ctx context.Context, ids []int64) error {
	return r.db.WithContext(ctx).Where("post_id IN ?", ids).Delete(&SysPost{}).Error
}

func (r *Repository) ExistsByCodeExcluding(ctx context.Context, code string, excludeID int64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&SysPost{}).
		Where("post_code = ? AND post_id <> ?", code, excludeID).
		Count(&count).Error
	return count > 0, err
}

func (r *Repository) ExistsByNameExcluding(ctx context.Context, name string, excludeID int64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&SysPost{}).
		Where("post_name = ? AND post_id <> ?", name, excludeID).
		Count(&count).Error
	return count > 0, err
}

func (r *Repository) FindAllForSelect(ctx context.Context) ([]SysPost, error) {
	var posts []SysPost
	err := r.db.WithContext(ctx).Where("status = ?", "0").Order("post_sort").Find(&posts).Error
	return posts, err
}
