package user

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserNameDuplicate  = errors.New("user name already exists")
	ErrPhoneDuplicate     = errors.New("phone number already exists")
	ErrEmailDuplicate     = errors.New("email already exists")
	ErrSelfOperation      = errors.New("cannot operate on yourself")
	ErrPasswordEmpty      = errors.New("password cannot be empty")
)

type Service struct {
	repo *Repository
}

func NewService(r *Repository) *Service {
	return &Service{repo: r}
}

func (s *Service) FindByUserName(ctx context.Context, userName string) (*SysUser, error) {
	return s.repo.FindByUserName(ctx, userName)
}

func (s *Service) FindPermissionsByUserID(ctx context.Context, userID int64) ([]string, error) {
	return s.repo.FindPermissionsByUserID(ctx, userID)
}

func (s *Service) FindRoleKeysByUserID(ctx context.Context, userID int64) ([]string, error) {
	return s.repo.FindRoleKeysByUserID(ctx, userID)
}

func (s *Service) ListUsers(ctx context.Context, req UserListRequest) (*UserListResponse, error) {
	if req.PageNum <= 0 {
		req.PageNum = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	offset := (req.PageNum - 1) * req.PageSize

	rows, total, err := s.repo.FindAll(ctx, offset, req.PageSize, userListFilter{
		UserName:    req.UserName,
		Phonenumber: req.Phonenumber,
		Status:      req.Status,
		DeptID:      req.DeptID,
		BeginTime:   req.BeginTime,
		EndTime:     req.EndTime,
	})
	if err != nil {
		return nil, err
	}

	return &UserListResponse{Rows: rows, Total: total}, nil
}

func (s *Service) GetUserDetail(ctx context.Context, userID int64) (*UserResponse, error) {
	u, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	resp := &UserResponse{
		UserID:      u.UserID,
		DeptID:      u.DeptID,
		UserName:    u.UserName,
		NickName:    u.NickName,
		UserType:    u.UserType,
		Email:       u.Email,
		Phonenumber: u.Phonenumber,
		Sex:         u.Sex,
		Avatar:      u.Avatar,
		Status:      u.Status,
		DelFlag:     u.DelFlag,
		LoginIP:     u.LoginIP,
		LoginDate:   u.LoginDate,
		CreateBy:    u.CreateBy,
		CreateTime:  u.CreateTime,
		UpdateBy:    u.UpdateBy,
		UpdateTime:  u.UpdateTime,
		Remark:      u.Remark,
	}

	// dept info
	if u.DeptID != nil {
		dept, err := s.repo.FindDeptByID(ctx, *u.DeptID)
		if err == nil {
			resp.DeptName = dept.DeptName
			resp.Dept = &DeptInfo{DeptID: dept.DeptID, DeptName: dept.DeptName}
		}
	}

	// role info
	roleIDs, _ := s.repo.FindRoleIDsByUserID(ctx, userID)
	if len(roleIDs) > 0 {
		allRoles, _ := s.repo.FindAllRoles(ctx)
		roleMap := make(map[int64]RoleInfo)
		for _, r := range allRoles {
			roleMap[r.RoleID] = r
		}
		for _, rid := range roleIDs {
			if ri, ok := roleMap[rid]; ok {
				resp.Roles = append(resp.Roles, ri)
			}
		}
	}

	// post info
	postIDs, _ := s.repo.FindPostIDsByUserID(ctx, userID)
	resp.PostIDs = postIDs

	return resp, nil
}

func (s *Service) AddUser(ctx context.Context, req AddUserRequest, loginUserID int64) error {
	userName := strings.TrimSpace(req.UserName)
	password := strings.TrimSpace(req.Password)

	if userName == "" {
		return fmt.Errorf("user name cannot be empty")
	}
	if password == "" {
		return ErrPasswordEmpty
	}

	exists, err := s.repo.ExistsByUserNameExcluding(ctx, userName, 0)
	if err != nil {
		return err
	}
	if exists {
		return ErrUserNameDuplicate
	}

	if req.Phonenumber != "" {
		exists, err := s.repo.ExistsByPhonenumberExcluding(ctx, req.Phonenumber, 0)
		if err != nil {
			return err
		}
		if exists {
			return ErrPhoneDuplicate
		}
	}

	if req.Email != "" {
		exists, err := s.repo.ExistsByEmailExcluding(ctx, req.Email, 0)
		if err != nil {
			return err
		}
		if exists {
			return ErrEmailDuplicate
		}
	}

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u := &SysUser{
		DeptID:      req.DeptID,
		UserName:    userName,
		NickName:    req.NickName,
		Email:       req.Email,
		Phonenumber: req.Phonenumber,
		Sex:         req.Sex,
		Status:      req.Status,
		Password:    string(hashedPwd),
		Remark:      req.Remark,
	}

	if u.Status == "" {
		u.Status = "0"
	}

	if err := s.repo.Create(ctx, u); err != nil {
		return err
	}

	if len(req.RoleIDs) > 0 {
		_ = s.repo.ReplaceUserRoles(ctx, u.UserID, req.RoleIDs)
	}
	if len(req.PostIDs) > 0 {
		_ = s.repo.ReplaceUserPosts(ctx, u.UserID, req.PostIDs)
	}

	return nil
}

func (s *Service) EditUser(ctx context.Context, req EditUserRequest, loginUserID int64) error {
	if req.UserID == loginUserID {
		return ErrSelfOperation
	}

	u, err := s.repo.FindByID(ctx, req.UserID)
	if err != nil {
		return ErrUserNotFound
	}

	if req.Phonenumber != "" {
		exists, err := s.repo.ExistsByPhonenumberExcluding(ctx, req.Phonenumber, req.UserID)
		if err != nil {
			return err
		}
		if exists {
			return ErrPhoneDuplicate
		}
	}

	if req.Email != "" {
		exists, err := s.repo.ExistsByEmailExcluding(ctx, req.Email, req.UserID)
		if err != nil {
			return err
		}
		if exists {
			return ErrEmailDuplicate
		}
	}

	u.DeptID = req.DeptID
	u.NickName = req.NickName
	u.Email = req.Email
	u.Phonenumber = req.Phonenumber
	u.Sex = req.Sex
	u.Status = req.Status
	u.Remark = req.Remark

	if err := s.repo.Update(ctx, u); err != nil {
		return err
	}

	_ = s.repo.ReplaceUserRoles(ctx, req.UserID, req.RoleIDs)
	_ = s.repo.ReplaceUserPosts(ctx, req.UserID, req.PostIDs)

	return nil
}

func (s *Service) DeleteUsers(ctx context.Context, userIDs []int64, loginUserID int64) error {
	for _, uid := range userIDs {
		if uid == loginUserID {
			return ErrSelfOperation
		}
	}
	return s.repo.SoftDelete(ctx, userIDs)
}

func (s *Service) ResetPassword(ctx context.Context, req ResetPwdRequest) error {
	if strings.TrimSpace(req.Password) == "" {
		return ErrPasswordEmpty
	}

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(strings.TrimSpace(req.Password)), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.repo.UpdatePassword(ctx, req.UserID, string(hashedPwd))
}

func (s *Service) ChangeStatus(ctx context.Context, req ChangeStatusRequest) error {
	return s.repo.UpdateStatus(ctx, req.UserID, req.Status)
}

func (s *Service) GetAuthRole(ctx context.Context, userID int64) (*AuthRoleResponse, error) {
	allRoles, err := s.repo.FindAllRoles(ctx)
	if err != nil {
		return nil, err
	}

	roleIDs, err := s.repo.FindRoleIDsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &AuthRoleResponse{
		Roles:       allRoles,
		UserRoleIDs: roleIDs,
	}, nil
}

func (s *Service) SaveAuthRole(ctx context.Context, req AuthRoleRequest) error {
	_, err := s.repo.FindByID(ctx, req.UserID)
	if err != nil {
		return ErrUserNotFound
	}
	return s.repo.ReplaceUserRoles(ctx, req.UserID, req.RoleIDs)
}

// SelectPostList returns all available posts.
func (s *Service) SelectPostList(ctx context.Context) ([]map[string]any, error) {
	return s.repo.SelectPostList(ctx)
}
