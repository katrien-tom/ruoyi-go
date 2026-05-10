package dept

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

var (
	errDeptHasChildren  = errors.New("dept has children")
	errDeptHasUsers     = errors.New("dept has assigned users")
	errDeptNameDuplicate = errors.New("dept name already exists")
)

type Service struct {
	repo *Repository
}

func NewService(r *Repository) *Service {
	return &Service{repo: r}
}

func (s *Service) GetDeptList(ctx context.Context, req DeptListRequest) ([]DeptResponse, error) {
	depts, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	// Apply filters
	var filtered []SysDept
	for _, d := range depts {
		if req.DeptName != "" && d.DeptName != req.DeptName {
			continue
		}
		if req.Status != "" && d.Status != req.Status {
			continue
		}
		filtered = append(filtered, d)
	}
	return buildDeptTree(filtered, 0), nil
}

func (s *Service) GetDeptTreeSelect(ctx context.Context) ([]DeptTreeSelectResponse, error) {
	depts, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return buildTreeSelect(depts, 0), nil
}

func (s *Service) GetByID(ctx context.Context, deptID int64) (*SysDept, error) {
	return s.repo.FindByID(ctx, deptID)
}

func (s *Service) AddDept(ctx context.Context, req AddDeptRequest) error {
	// Check duplicate name under same parent
	dup, err := s.repo.ExistsByNameExcluding(ctx, req.ParentID, req.DeptName, 0)
	if err != nil {
		return err
	}
	if dup {
		return errDeptNameDuplicate
	}

	ancestors := "0"
	if req.ParentID != 0 {
		parent, err := s.repo.FindByID(ctx, req.ParentID)
		if err != nil {
			return fmt.Errorf("parent dept not found")
		}
		ancestors = fmt.Sprintf("%s,%d", parent.Ancestors, parent.DeptID)
	}

	d := &SysDept{
		ParentID:  req.ParentID,
		Ancestors: ancestors,
		DeptName:  req.DeptName,
		OrderNum:  req.OrderNum,
		Leader:    req.Leader,
		Phone:     req.Phone,
		Email:     req.Email,
		Status:    req.Status,
	}
	if d.Status == "" {
		d.Status = "0"
	}
	return s.repo.Create(ctx, d)
}

func (s *Service) EditDept(ctx context.Context, req EditDeptRequest) error {
	// Check duplicate name
	dup, err := s.repo.ExistsByNameExcluding(ctx, req.ParentID, req.DeptName, req.DeptID)
	if err != nil {
		return err
	}
	if dup {
		return errDeptNameDuplicate
	}

	dept, err := s.repo.FindByID(ctx, req.DeptID)
	if err != nil {
		return err
	}

	newAncestors := "0"
	if req.ParentID != 0 {
		parent, err := s.repo.FindByID(ctx, req.ParentID)
		if err != nil {
			return fmt.Errorf("parent dept not found")
		}
		newAncestors = fmt.Sprintf("%s,%d", parent.Ancestors, parent.DeptID)
	}

	oldAncestors := fmt.Sprintf("%s,%d", dept.Ancestors, dept.DeptID)

	dept.ParentID = req.ParentID
	dept.Ancestors = newAncestors
	dept.DeptName = req.DeptName
	dept.OrderNum = req.OrderNum
	dept.Leader = req.Leader
	dept.Phone = req.Phone
	dept.Email = req.Email
	dept.Status = req.Status

	return s.repo.DB().Transaction(func(tx *gorm.DB) error {
		txRepo := &Repository{db: tx}
		if err := txRepo.Update(ctx, dept); err != nil {
			return err
		}
		if newAncestors != oldAncestors {
			if err := txRepo.UpdateChildrenAncestors(ctx, oldAncestors, fmt.Sprintf("%s,%d", newAncestors, dept.DeptID)); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *Service) DeleteDept(ctx context.Context, deptID int64) error {
	hasChildren, err := s.repo.HasChildren(ctx, deptID)
	if err != nil {
		return err
	}
	if hasChildren {
		return errDeptHasChildren
	}

	hasUsers, err := s.repo.ExistsUsersInDept(ctx, deptID)
	if err != nil {
		return err
	}
	if hasUsers {
		return errDeptHasUsers
	}

	return s.repo.Delete(ctx, deptID)
}

func buildDeptTree(depts []SysDept, parentID int64) []DeptResponse {
	var result []DeptResponse
	for _, d := range depts {
		if d.ParentID != parentID {
			continue
		}
		node := DeptResponse{
			DeptID:    d.DeptID,
			ParentID:  d.ParentID,
			Ancestors: d.Ancestors,
			DeptName:  d.DeptName,
			OrderNum:  d.OrderNum,
			Leader:    d.Leader,
			Phone:     d.Phone,
			Email:     d.Email,
			Status:    d.Status,
			DelFlag:   d.DelFlag,
			Children:  buildDeptTree(depts, d.DeptID),
		}
		if node.Children == nil {
			node.Children = []DeptResponse{}
		}
		result = append(result, node)
	}
	if result == nil {
		result = []DeptResponse{}
	}
	return result
}

func buildTreeSelect(depts []SysDept, parentID int64) []DeptTreeSelectResponse {
	var result []DeptTreeSelectResponse
	for _, d := range depts {
		if d.ParentID != parentID {
			continue
		}
		node := DeptTreeSelectResponse{
			ID:       d.DeptID,
			Label:    d.DeptName,
			Children: buildTreeSelect(depts, d.DeptID),
		}
		if node.Children == nil {
			node.Children = []DeptTreeSelectResponse{}
		}
		result = append(result, node)
	}
	if result == nil {
		result = []DeptTreeSelectResponse{}
	}
	return result
}
