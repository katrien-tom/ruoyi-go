package config

import (
	"context"
	"errors"
)

var errConfigKeyDuplicate = errors.New("config key already exists")

type Service struct {
	repo *Repository
}

func NewService(r *Repository) *Service {
	return &Service{repo: r}
}

func (s *Service) ListConfigs(ctx context.Context, req ConfigListRequest, pageNum, pageSize int) (*ConfigListResponse, error) {
	if pageNum <= 0 {
		pageNum = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	offset := (pageNum - 1) * pageSize

	configs, total, err := s.repo.FindAll(ctx, offset, pageSize, req.ConfigName, req.ConfigKey, req.ConfigType)
	if err != nil {
		return nil, err
	}

	rows := make([]ConfigResponse, len(configs))
	for i, c := range configs {
		rows[i] = ConfigResponse{
			ConfigID:    c.ConfigID,
			ConfigName:  c.ConfigName,
			ConfigKey:   c.ConfigKey,
			ConfigValue: c.ConfigValue,
			ConfigType:  c.ConfigType,
			Remark:      c.Remark,
		}
	}
	return &ConfigListResponse{Rows: rows, Total: total}, nil
}

func (s *Service) GetByID(ctx context.Context, id int) (*ConfigResponse, error) {
	c, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &ConfigResponse{
		ConfigID:    c.ConfigID,
		ConfigName:  c.ConfigName,
		ConfigKey:   c.ConfigKey,
		ConfigValue: c.ConfigValue,
		ConfigType:  c.ConfigType,
		Remark:      c.Remark,
	}, nil
}

func (s *Service) GetByKey(ctx context.Context, key string) (*ConfigResponse, error) {
	c, err := s.repo.FindByKey(ctx, key)
	if err != nil {
		return nil, err
	}
	return &ConfigResponse{
		ConfigID:    c.ConfigID,
		ConfigName:  c.ConfigName,
		ConfigKey:   c.ConfigKey,
		ConfigValue: c.ConfigValue,
		ConfigType:  c.ConfigType,
		Remark:      c.Remark,
	}, nil
}

func (s *Service) EditConfig(ctx context.Context, req EditConfigRequest) error {
	dup, err := s.repo.ExistsByKeyExcluding(ctx, req.ConfigKey, req.ConfigID)
	if err != nil {
		return err
	}
	if dup {
		return errConfigKeyDuplicate
	}

	cfg := &SysConfig{
		ConfigID:    req.ConfigID,
		ConfigName:  req.ConfigName,
		ConfigKey:   req.ConfigKey,
		ConfigValue: req.ConfigValue,
		ConfigType:  req.ConfigType,
		Remark:      req.Remark,
	}
	return s.repo.Update(ctx, cfg)
}

func (s *Service) DeleteConfigs(ctx context.Context, ids []int) error {
	return s.repo.Delete(ctx, ids)
}
