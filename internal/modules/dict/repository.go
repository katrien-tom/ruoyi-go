package dict

import (
	"context"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

type dictTypeFilter struct {
	DictName string
	DictType string
	Status   string
}

type dictDataFilter struct {
	DictType string
	DictLabel string
	Status   string
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// DictType operations

func (r *Repository) FindDictTypes(ctx context.Context, offset, limit int, filter dictTypeFilter) ([]SysDictType, int64, error) {
	var total int64
	base := r.db.WithContext(ctx).Model(&SysDictType{})

	if filter.DictName != "" {
		base = base.Where("dict_name LIKE ?", "%"+filter.DictName+"%")
	}
	if filter.DictType != "" {
		base = base.Where("dict_type LIKE ?", "%"+filter.DictType+"%")
	}
	if filter.Status != "" {
		base = base.Where("status = ?", filter.Status)
	}

	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var types []SysDictType
	if err := base.Order("dict_id").Offset(offset).Limit(limit).Find(&types).Error; err != nil {
		return nil, 0, err
	}
	return types, total, nil
}

func (r *Repository) FindDictTypeByID(ctx context.Context, id int64) (*SysDictType, error) {
	var dt SysDictType
	if err := r.db.WithContext(ctx).Where("dict_id = ?", id).First(&dt).Error; err != nil {
		return nil, err
	}
	return &dt, nil
}

func (r *Repository) CreateDictType(ctx context.Context, dt *SysDictType) error {
	return r.db.WithContext(ctx).Create(dt).Error
}

func (r *Repository) UpdateDictType(ctx context.Context, dt *SysDictType) error {
	return r.db.WithContext(ctx).Model(dt).Where("dict_id = ?", dt.DictID).Updates(dt).Error
}

func (r *Repository) DeleteDictTypes(ctx context.Context, ids []int64) error {
	return r.db.WithContext(ctx).Where("dict_id IN ?", ids).Delete(&SysDictType{}).Error
}

func (r *Repository) ExistsDictTypeExcluding(ctx context.Context, dictType string, excludeID int64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&SysDictType{}).
		Where("dict_type = ? AND dict_id <> ?", dictType, excludeID).
		Count(&count).Error
	return count > 0, err
}

// DictData operations

func (r *Repository) FindDictDatas(ctx context.Context, offset, limit int, filter dictDataFilter) ([]SysDictData, int64, error) {
	var total int64
	base := r.db.WithContext(ctx).Model(&SysDictData{})

	if filter.DictType != "" {
		base = base.Where("dict_type = ?", filter.DictType)
	}
	if filter.DictLabel != "" {
		base = base.Where("dict_label LIKE ?", "%"+filter.DictLabel+"%")
	}
	if filter.Status != "" {
		base = base.Where("status = ?", filter.Status)
	}

	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var data []SysDictData
	if err := base.Order("dict_sort").Offset(offset).Limit(limit).Find(&data).Error; err != nil {
		return nil, 0, err
	}
	return data, total, nil
}

func (r *Repository) FindDictDataByID(ctx context.Context, id int64) (*SysDictData, error) {
	var dd SysDictData
	if err := r.db.WithContext(ctx).Where("dict_code = ?", id).First(&dd).Error; err != nil {
		return nil, err
	}
	return &dd, nil
}

func (r *Repository) FindDictDataByType(ctx context.Context, dictType string) ([]SysDictData, error) {
	var data []SysDictData
	if err := r.db.WithContext(ctx).Where("dict_type = ? AND status = ?", dictType, "0").Order("dict_sort").Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (r *Repository) CreateDictData(ctx context.Context, dd *SysDictData) error {
	return r.db.WithContext(ctx).Create(dd).Error
}

func (r *Repository) UpdateDictData(ctx context.Context, dd *SysDictData) error {
	return r.db.WithContext(ctx).Model(dd).Where("dict_code = ?", dd.DictCode).Updates(dd).Error
}

func (r *Repository) DeleteDictDatas(ctx context.Context, ids []int64) error {
	return r.db.WithContext(ctx).Where("dict_code IN ?", ids).Delete(&SysDictData{}).Error
}
