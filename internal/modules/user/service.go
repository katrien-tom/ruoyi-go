package user

import "context"

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
