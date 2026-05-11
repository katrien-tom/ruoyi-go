package operlog

import "context"

type Service struct {
	repo *Repository
}

func NewService(r *Repository) *Service {
	return &Service{repo: r}
}

func (s *Service) ListLogs(ctx context.Context, req OperLogListRequest, pageNum, pageSize int) (*OperLogListResponse, error) {
	if pageNum <= 0 {
		pageNum = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	offset := (pageNum - 1) * pageSize

	logs, total, err := s.repo.FindAll(ctx, offset, pageSize, operLogFilter{
		Title:        req.Title,
		BusinessType: req.BusinessType,
		OperName:     req.OperName,
		Status:       req.Status,
		BeginTime:    req.BeginTime,
		EndTime:      req.EndTime,
	})
	if err != nil {
		return nil, err
	}

	rows := make([]OperLogResponse, len(logs))
	for i, l := range logs {
		rows[i] = OperLogResponse{
			OperID:        l.OperID,
			Title:         l.Title,
			BusinessType:  l.BusinessType,
			Method:        l.Method,
			RequestMethod: l.RequestMethod,
			OperatorType:  l.OperatorType,
			OperName:      l.OperName,
			DeptName:      l.DeptName,
			OperURL:       l.OperURL,
			OperIP:        l.OperIP,
			OperLocation:  l.OperLocation,
			OperParam:     l.OperParam,
			JSONResult:    l.JSONResult,
			Status:        l.Status,
			ErrorMsg:      l.ErrorMsg,
			OperTime:      l.OperTime,
			CostTime:      l.CostTime,
		}
	}
	return &OperLogListResponse{Rows: rows, Total: total}, nil
}

func (s *Service) DeleteLogs(ctx context.Context, ids []int64) error {
	return s.repo.Delete(ctx, ids)
}

func (s *Service) Clean(ctx context.Context) error {
	return s.repo.Clean(ctx)
}
