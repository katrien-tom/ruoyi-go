package config

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

func (r *Repository) FindAll(ctx context.Context, offset, limit int, configName, configKey, configType string) ([]SysConfig, int64, error) {
	var total int64
	base := r.db.WithContext(ctx).Model(&SysConfig{})

	if configName != "" {
		base = base.Where("config_name LIKE ?", "%"+configName+"%")
	}
	if configKey != "" {
		base = base.Where("config_key LIKE ?", "%"+configKey+"%")
	}
	if configType != "" {
		base = base.Where("config_type = ?", configType)
	}

	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var configs []SysConfig
	if err := base.Order("config_id").Offset(offset).Limit(limit).Find(&configs).Error; err != nil {
		return nil, 0, err
	}
	return configs, total, nil
}

func (r *Repository) FindByID(ctx context.Context, id int) (*SysConfig, error) {
	var cfg SysConfig
	if err := r.db.WithContext(ctx).Where("config_id = ?", id).First(&cfg).Error; err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (r *Repository) FindByKey(ctx context.Context, key string) (*SysConfig, error) {
	var cfg SysConfig
	if err := r.db.WithContext(ctx).Where("config_key = ?", key).First(&cfg).Error; err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (r *Repository) Update(ctx context.Context, cfg *SysConfig) error {
	return r.db.WithContext(ctx).Model(cfg).Where("config_id = ?", cfg.ConfigID).Updates(cfg).Error
}

func (r *Repository) Delete(ctx context.Context, ids []int) error {
	return r.db.WithContext(ctx).Where("config_id IN ?", ids).Delete(&SysConfig{}).Error
}

func (r *Repository) ExistsByKeyExcluding(ctx context.Context, key string, excludeID int) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&SysConfig{}).
		Where("config_key = ? AND config_id <> ?", key, excludeID).
		Count(&count).Error
	return count > 0, err
}
