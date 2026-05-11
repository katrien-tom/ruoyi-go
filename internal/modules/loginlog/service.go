package loginlog

import "context"

type Service struct {
	repo *Repository
}

func NewService(r *Repository) *Service {
	return &Service{repo: r}
}

func (s *Service) ListLogs(ctx context.Context, req LoginLogListRequest, pageNum, pageSize int) (*LoginLogListResponse, error) {
	if pageNum <= 0 {
		pageNum = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	offset := (pageNum - 1) * pageSize

	logs, total, err := s.repo.FindAll(ctx, offset, pageSize, loginLogFilter{
		UserName:  req.UserName,
		IPAddr:    req.IPAddr,
		Status:    req.Status,
		BeginTime: req.BeginTime,
		EndTime:   req.EndTime,
	})
	if err != nil {
		return nil, err
	}

	rows := make([]LoginLogResponse, len(logs))
	for i, l := range logs {
		rows[i] = LoginLogResponse{
			InfoID:        l.InfoID,
			UserName:      l.UserName,
			IPAddr:        l.IPAddr,
			LoginLocation: l.LoginLocation,
			Browser:       l.Browser,
			OS:            l.OS,
			Status:        l.Status,
			Msg:           l.Msg,
			LoginTime:     l.LoginTime,
		}
	}
	return &LoginLogListResponse{Rows: rows, Total: total}, nil
}

func (s *Service) DeleteLogs(ctx context.Context, ids []int64) error {
	return s.repo.Delete(ctx, ids)
}

func (s *Service) Clean(ctx context.Context) error {
	return s.repo.Clean(ctx)
}
