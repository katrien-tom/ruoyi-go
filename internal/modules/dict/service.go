package dict

import (
	"context"
	"errors"
)

var errDictTypeDuplicate = errors.New("dict type already exists")

type Service struct {
	repo *Repository
}

func NewService(r *Repository) *Service {
	return &Service{repo: r}
}

// DictType

func (s *Service) ListDictTypes(ctx context.Context, req DictTypeListRequest, pageNum, pageSize int) (*DictTypeListResponse, error) {
	if pageNum <= 0 {
		pageNum = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	offset := (pageNum - 1) * pageSize

	types, total, err := s.repo.FindDictTypes(ctx, offset, pageSize, dictTypeFilter{
		DictName: req.DictName,
		DictType: req.DictType,
		Status:   req.Status,
	})
	if err != nil {
		return nil, err
	}

	rows := make([]DictTypeResponse, len(types))
	for i, t := range types {
		rows[i] = DictTypeResponse{
			DictID:   t.DictID,
			DictName: t.DictName,
			DictType: t.DictType,
			Status:   t.Status,
			Remark:   t.Remark,
		}
	}
	return &DictTypeListResponse{Rows: rows, Total: total}, nil
}

func (s *Service) GetDictTypeByID(ctx context.Context, id int64) (*DictTypeResponse, error) {
	t, err := s.repo.FindDictTypeByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &DictTypeResponse{
		DictID:   t.DictID,
		DictName: t.DictName,
		DictType: t.DictType,
		Status:   t.Status,
		Remark:   t.Remark,
	}, nil
}

func (s *Service) AddDictType(ctx context.Context, req AddDictTypeRequest) error {
	dup, err := s.repo.ExistsDictTypeExcluding(ctx, req.DictType, 0)
	if err != nil {
		return err
	}
	if dup {
		return errDictTypeDuplicate
	}

	dt := &SysDictType{
		DictName: req.DictName,
		DictType: req.DictType,
		Status:   req.Status,
		Remark:   req.Remark,
	}
	if dt.Status == "" {
		dt.Status = "0"
	}
	return s.repo.CreateDictType(ctx, dt)
}

func (s *Service) EditDictType(ctx context.Context, req EditDictTypeRequest) error {
	dup, err := s.repo.ExistsDictTypeExcluding(ctx, req.DictType, req.DictID)
	if err != nil {
		return err
	}
	if dup {
		return errDictTypeDuplicate
	}

	dt := &SysDictType{
		DictID:   req.DictID,
		DictName: req.DictName,
		DictType: req.DictType,
		Status:   req.Status,
		Remark:   req.Remark,
	}
	return s.repo.UpdateDictType(ctx, dt)
}

func (s *Service) DeleteDictTypes(ctx context.Context, ids []int64) error {
	return s.repo.DeleteDictTypes(ctx, ids)
}

// DictData

func (s *Service) ListDictDatas(ctx context.Context, req DictDataListRequest, pageNum, pageSize int) (*DictDataListResponse, error) {
	if pageNum <= 0 {
		pageNum = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	offset := (pageNum - 1) * pageSize

	data, total, err := s.repo.FindDictDatas(ctx, offset, pageSize, dictDataFilter{
		DictType:  req.DictType,
		DictLabel: req.DictLabel,
		Status:    req.Status,
	})
	if err != nil {
		return nil, err
	}

	rows := make([]DictDataResponse, len(data))
	for i, d := range data {
		rows[i] = DictDataResponse{
			DictCode:  d.DictCode,
			DictSort:  d.DictSort,
			DictLabel: d.DictLabel,
			DictValue: d.DictValue,
			DictType:  d.DictType,
			CSSClass:  d.CSSClass,
			ListClass: d.ListClass,
			IsDefault: d.IsDefault,
			Status:    d.Status,
			Remark:    d.Remark,
		}
	}
	return &DictDataListResponse{Rows: rows, Total: total}, nil
}

func (s *Service) GetDictDataByType(ctx context.Context, dictType string) ([]DictDataResponse, error) {
	data, err := s.repo.FindDictDataByType(ctx, dictType)
	if err != nil {
		return nil, err
	}
	rows := make([]DictDataResponse, len(data))
	for i, d := range data {
		rows[i] = DictDataResponse{
			DictCode:  d.DictCode,
			DictSort:  d.DictSort,
			DictLabel: d.DictLabel,
			DictValue: d.DictValue,
			DictType:  d.DictType,
			CSSClass:  d.CSSClass,
			ListClass: d.ListClass,
			IsDefault: d.IsDefault,
			Status:    d.Status,
			Remark:    d.Remark,
		}
	}
	return rows, nil
}

func (s *Service) GetDictDataByID(ctx context.Context, id int64) (*DictDataResponse, error) {
	d, err := s.repo.FindDictDataByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &DictDataResponse{
		DictCode:  d.DictCode,
		DictSort:  d.DictSort,
		DictLabel: d.DictLabel,
		DictValue: d.DictValue,
		DictType:  d.DictType,
		CSSClass:  d.CSSClass,
		ListClass: d.ListClass,
		IsDefault: d.IsDefault,
		Status:    d.Status,
		Remark:    d.Remark,
	}, nil
}

func (s *Service) AddDictData(ctx context.Context, req AddDictDataRequest) error {
	dd := &SysDictData{
		DictType:  req.DictType,
		DictLabel: req.DictLabel,
		DictValue: req.DictValue,
		DictSort:  req.DictSort,
		CSSClass:  req.CSSClass,
		ListClass: req.ListClass,
		IsDefault: req.IsDefault,
		Status:    req.Status,
		Remark:    req.Remark,
	}
	if dd.Status == "" {
		dd.Status = "0"
	}
	if dd.IsDefault == "" {
		dd.IsDefault = "N"
	}
	return s.repo.CreateDictData(ctx, dd)
}

func (s *Service) EditDictData(ctx context.Context, req EditDictDataRequest) error {
	dd := &SysDictData{
		DictCode:  req.DictCode,
		DictType:  req.DictType,
		DictLabel: req.DictLabel,
		DictValue: req.DictValue,
		DictSort:  req.DictSort,
		CSSClass:  req.CSSClass,
		ListClass: req.ListClass,
		IsDefault: req.IsDefault,
		Status:    req.Status,
		Remark:    req.Remark,
	}
	return s.repo.UpdateDictData(ctx, dd)
}

func (s *Service) DeleteDictDatas(ctx context.Context, ids []int64) error {
	return s.repo.DeleteDictDatas(ctx, ids)
}
