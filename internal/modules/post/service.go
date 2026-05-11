package post

import (
	"context"
	"errors"
)

var (
	errPostCodeDuplicate = errors.New("post code already exists")
	errPostNameDuplicate = errors.New("post name already exists")
)

type Service struct {
	repo *Repository
}

func NewService(r *Repository) *Service {
	return &Service{repo: r}
}

func (s *Service) ListPosts(ctx context.Context, req PostListRequest, pageNum, pageSize int) (*PostListResponse, error) {
	if pageNum <= 0 {
		pageNum = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	offset := (pageNum - 1) * pageSize

	posts, total, err := s.repo.FindAll(ctx, offset, pageSize, postListFilter{
		PostCode: req.PostCode,
		PostName: req.PostName,
		Status:   req.Status,
	})
	if err != nil {
		return nil, err
	}

	rows := make([]PostResponse, len(posts))
	for i, p := range posts {
		rows[i] = PostResponse{
			PostID:   p.PostID,
			PostCode: p.PostCode,
			PostName: p.PostName,
			PostSort: p.PostSort,
			Status:   p.Status,
			Remark:   p.Remark,
		}
	}
	return &PostListResponse{Rows: rows, Total: total}, nil
}

func (s *Service) GetByID(ctx context.Context, id int64) (*PostResponse, error) {
	p, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &PostResponse{
		PostID:   p.PostID,
		PostCode: p.PostCode,
		PostName: p.PostName,
		PostSort: p.PostSort,
		Status:   p.Status,
		Remark:   p.Remark,
	}, nil
}

func (s *Service) AddPost(ctx context.Context, req AddPostRequest) error {
	dup, err := s.repo.ExistsByCodeExcluding(ctx, req.PostCode, 0)
	if err != nil {
		return err
	}
	if dup {
		return errPostCodeDuplicate
	}

	dup, err = s.repo.ExistsByNameExcluding(ctx, req.PostName, 0)
	if err != nil {
		return err
	}
	if dup {
		return errPostNameDuplicate
	}

	p := &SysPost{
		PostCode: req.PostCode,
		PostName: req.PostName,
		PostSort: req.PostSort,
		Status:   req.Status,
		Remark:   req.Remark,
	}
	if p.Status == "" {
		p.Status = "0"
	}
	return s.repo.Create(ctx, p)
}

func (s *Service) EditPost(ctx context.Context, req EditPostRequest) error {
	dup, err := s.repo.ExistsByCodeExcluding(ctx, req.PostCode, req.PostID)
	if err != nil {
		return err
	}
	if dup {
		return errPostCodeDuplicate
	}

	dup, err = s.repo.ExistsByNameExcluding(ctx, req.PostName, req.PostID)
	if err != nil {
		return err
	}
	if dup {
		return errPostNameDuplicate
	}

	p := &SysPost{
		PostID:   req.PostID,
		PostCode: req.PostCode,
		PostName: req.PostName,
		PostSort: req.PostSort,
		Status:   req.Status,
		Remark:   req.Remark,
	}
	return s.repo.Update(ctx, p)
}

func (s *Service) DeletePosts(ctx context.Context, ids []int64) error {
	return s.repo.Delete(ctx, ids)
}

func (s *Service) OptionSelect(ctx context.Context) ([]map[string]any, error) {
	posts, err := s.repo.FindAllForSelect(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]map[string]any, len(posts))
	for i, p := range posts {
		result[i] = map[string]any{
			"postId":   p.PostID,
			"postCode": p.PostCode,
			"postName": p.PostName,
		}
	}
	return result, nil
}
