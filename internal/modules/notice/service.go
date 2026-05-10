package notice

import (
	"context"
	"time"
)

type Service struct {
	repo *Repository
	now  func() time.Time
}

func NewService(r *Repository) *Service {
	return &Service{repo: r, now: time.Now}
}

func (s *Service) ListNotices(ctx context.Context, req NoticeListRequest, pageNum, pageSize int) (*NoticeListResponse, error) {
	if pageNum <= 0 {
		pageNum = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	offset := (pageNum - 1) * pageSize

	notices, total, err := s.repo.FindAll(ctx, offset, pageSize, req.NoticeTitle, req.NoticeType, req.CreateBy)
	if err != nil {
		return nil, err
	}

	rows := make([]NoticeResponse, len(notices))
	for i, n := range notices {
		createTime := ""
		if n.CreateTime != nil {
			createTime = n.CreateTime.Format(time.DateTime)
		}
		rows[i] = NoticeResponse{
			NoticeID:      n.NoticeID,
			NoticeTitle:   n.NoticeTitle,
			NoticeType:    n.NoticeType,
			NoticeContent: string(n.NoticeContent),
			Status:        n.Status,
			CreateBy:      n.CreateBy,
			CreateTime:    createTime,
			Remark:        n.Remark,
		}
	}
	return &NoticeListResponse{Rows: rows, Total: total}, nil
}

func (s *Service) GetByID(ctx context.Context, id int) (*NoticeResponse, error) {
	n, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	createTime := ""
	if n.CreateTime != nil {
		createTime = n.CreateTime.Format(time.DateTime)
	}
	return &NoticeResponse{
		NoticeID:      n.NoticeID,
		NoticeTitle:   n.NoticeTitle,
		NoticeType:    n.NoticeType,
		NoticeContent: string(n.NoticeContent),
		Status:        n.Status,
		CreateBy:      n.CreateBy,
		CreateTime:    createTime,
		Remark:        n.Remark,
	}, nil
}

func (s *Service) AddNotice(ctx context.Context, req AddNoticeRequest) error {
	now := s.now()
	n := &SysNotice{
		NoticeTitle:   req.NoticeTitle,
		NoticeType:    req.NoticeType,
		NoticeContent: []byte(req.NoticeContent),
		Status:        req.Status,
		Remark:        req.Remark,
		CreateTime:    &now,
	}
	if n.Status == "" {
		n.Status = "0"
	}
	return s.repo.Create(ctx, n)
}

func (s *Service) EditNotice(ctx context.Context, req EditNoticeRequest) error {
	n := &SysNotice{
		NoticeID:      req.NoticeID,
		NoticeTitle:   req.NoticeTitle,
		NoticeType:    req.NoticeType,
		NoticeContent: []byte(req.NoticeContent),
		Status:        req.Status,
		Remark:        req.Remark,
	}
	return s.repo.Update(ctx, n)
}

func (s *Service) DeleteNotices(ctx context.Context, ids []int) error {
	return s.repo.Delete(ctx, ids)
}
