package role

import (
	"context"
	"errors"
)

var (
	errRoleNameDuplicate = errors.New("role name already exists")
	errRoleKeyDuplicate  = errors.New("role key already exists")
)

type Service struct {
	repo *Repository
}

func NewService(r *Repository) *Service {
	return &Service{repo: r}
}

func (s *Service) ListRoles(ctx context.Context, req RoleListRequest, pageNum, pageSize int) (*RoleListResponse, error) {
	if pageNum <= 0 {
		pageNum = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	offset := (pageNum - 1) * pageSize

	roles, total, err := s.repo.FindAll(ctx, offset, pageSize, roleListFilter{
		RoleName: req.RoleName,
		RoleKey:  req.RoleKey,
		Status:   req.Status,
	})
	if err != nil {
		return nil, err
	}

	rows := make([]RoleResponse, len(roles))
	for i, r := range roles {
		rows[i] = RoleResponse{
			RoleID:            r.RoleID,
			RoleName:          r.RoleName,
			RoleKey:           r.RoleKey,
			RoleSort:          r.RoleSort,
			DataScope:         r.DataScope,
			MenuCheckStrictly: r.MenuCheckStrictly,
			DeptCheckStrictly: r.DeptCheckStrictly,
			Status:            r.Status,
			DelFlag:           r.DelFlag,
			Remark:            r.Remark,
		}
	}
	return &RoleListResponse{Rows: rows, Total: total}, nil
}

func (s *Service) GetByID(ctx context.Context, id int64) (*RoleResponse, error) {
	r, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &RoleResponse{
		RoleID:            r.RoleID,
		RoleName:          r.RoleName,
		RoleKey:           r.RoleKey,
		RoleSort:          r.RoleSort,
		DataScope:         r.DataScope,
		MenuCheckStrictly: r.MenuCheckStrictly,
		DeptCheckStrictly: r.DeptCheckStrictly,
		Status:            r.Status,
		DelFlag:           r.DelFlag,
		Remark:            r.Remark,
	}, nil
}

func (s *Service) AddRole(ctx context.Context, req AddRoleRequest) error {
	dup, err := s.repo.ExistsByNameExcluding(ctx, req.RoleName, 0)
	if err != nil {
		return err
	}
	if dup {
		return errRoleNameDuplicate
	}

	dup, err = s.repo.ExistsByKeyExcluding(ctx, req.RoleKey, 0)
	if err != nil {
		return err
	}
	if dup {
		return errRoleKeyDuplicate
	}

	r := &SysRole{
		RoleName:          req.RoleName,
		RoleKey:           req.RoleKey,
		RoleSort:          req.RoleSort,
		DataScope:         req.DataScope,
		MenuCheckStrictly: req.MenuCheckStrictly,
		DeptCheckStrictly: req.DeptCheckStrictly,
		Status:            req.Status,
		Remark:            req.Remark,
	}
	if r.Status == "" {
		r.Status = "0"
	}
	if r.DataScope == "" {
		r.DataScope = "1"
	}
	if r.MenuCheckStrictly == 0 {
		r.MenuCheckStrictly = 1
	}
	if r.DeptCheckStrictly == 0 {
		r.DeptCheckStrictly = 1
	}

	if err := s.repo.Create(ctx, r); err != nil {
		return err
	}

	if len(req.MenuIDs) > 0 {
		_ = s.repo.ReplaceRoleMenus(ctx, r.RoleID, req.MenuIDs)
	}

	return nil
}

func (s *Service) EditRole(ctx context.Context, req EditRoleRequest) error {
	dup, err := s.repo.ExistsByNameExcluding(ctx, req.RoleName, req.RoleID)
	if err != nil {
		return err
	}
	if dup {
		return errRoleNameDuplicate
	}

	dup, err = s.repo.ExistsByKeyExcluding(ctx, req.RoleKey, req.RoleID)
	if err != nil {
		return err
	}
	if dup {
		return errRoleKeyDuplicate
	}

	r := &SysRole{
		RoleID:            req.RoleID,
		RoleName:          req.RoleName,
		RoleKey:           req.RoleKey,
		RoleSort:          req.RoleSort,
		DataScope:         req.DataScope,
		MenuCheckStrictly: req.MenuCheckStrictly,
		DeptCheckStrictly: req.DeptCheckStrictly,
		Status:            req.Status,
		Remark:            req.Remark,
	}

	if err := s.repo.Update(ctx, r); err != nil {
		return err
	}

	_ = s.repo.ReplaceRoleMenus(ctx, req.RoleID, req.MenuIDs)

	return nil
}

func (s *Service) DeleteRoles(ctx context.Context, ids []int64) error {
	return s.repo.SoftDelete(ctx, ids)
}

func (s *Service) ChangeStatus(ctx context.Context, req ChangeRoleStatusRequest) error {
	return s.repo.UpdateStatus(ctx, req.RoleID, req.Status)
}

func (s *Service) UpdateDataScope(ctx context.Context, req DataScopeRequest) error {
	_, err := s.repo.FindByID(ctx, req.RoleID)
	if err != nil {
		return err
	}

	if err := s.repo.ReplaceRoleDepts(ctx, req.RoleID, req.DeptIDs); err != nil {
		return err
	}

	r := &SysRole{
		RoleID:    req.RoleID,
		DataScope: req.DataScope,
	}
	return s.repo.Update(ctx, r)
}

func (s *Service) GetAllocatedUsers(ctx context.Context, roleID int64, userName, phonenumber string, pageNum, pageSize int) (*AuthUserListResponse, error) {
	if pageNum <= 0 {
		pageNum = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	offset := (pageNum - 1) * pageSize

	rows, total, err := s.repo.FindAllocatedUsers(ctx, roleID, userName, phonenumber, offset, pageSize)
	if err != nil {
		return nil, err
	}
	return &AuthUserListResponse{Rows: rows, Total: total}, nil
}

func (s *Service) GetUnallocatedUsers(ctx context.Context, roleID int64, userName, phonenumber string, pageNum, pageSize int) (*AuthUserListResponse, error) {
	if pageNum <= 0 {
		pageNum = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	offset := (pageNum - 1) * pageSize

	rows, total, err := s.repo.FindUnallocatedUsers(ctx, roleID, userName, phonenumber, offset, pageSize)
	if err != nil {
		return nil, err
	}
	return &AuthUserListResponse{Rows: rows, Total: total}, nil
}

func (s *Service) CancelAuthUser(ctx context.Context, roleID, userID int64) error {
	return s.repo.CancelAuthUser(ctx, roleID, userID)
}

func (s *Service) CancelAuthUsers(ctx context.Context, roleID int64, userIDs []int64) error {
	return s.repo.CancelAuthUsers(ctx, roleID, userIDs)
}

func (s *Service) SelectAuthUsers(ctx context.Context, roleID int64, userIDs []int64) error {
	return s.repo.SelectAuthUsers(ctx, roleID, userIDs)
}

func (s *Service) GetDeptTreeForRole(ctx context.Context, roleID int64) (map[string]any, error) {
	checkedKeys, err := s.repo.FindRoleDeptIDs(ctx, roleID)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"checkedKeys": checkedKeys,
	}, nil
}
