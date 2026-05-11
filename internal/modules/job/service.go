package job

import "context"

type Service struct {
	repo *Repository
}

func NewService(r *Repository) *Service {
	return &Service{repo: r}
}

func (s *Service) FindJobByID(ctx context.Context, id int64) (*SysJob, error) {
	return s.repo.FindJobByID(ctx, id)
}

func (s *Service) FindAllJobs(ctx context.Context) ([]SysJob, error) {
	return s.repo.FindAllJobs(ctx)
}

func (s *Service) FindJobLogs(ctx context.Context, offset, limit int) ([]SysJobLog, int64, error) {
	return s.repo.FindJobLogs(ctx, offset, limit)
}
