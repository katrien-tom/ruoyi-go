package job

import (
	"context"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindJobByID(ctx context.Context, id int64) (*SysJob, error) {
	var job SysJob
	if err := r.db.WithContext(ctx).Where("job_id = ?", id).First(&job).Error; err != nil {
		return nil, err
	}
	return &job, nil
}

func (r *Repository) FindAllJobs(ctx context.Context) ([]SysJob, error) {
	var jobs []SysJob
	if err := r.db.WithContext(ctx).Find(&jobs).Error; err != nil {
		return nil, err
	}
	return jobs, nil
}

func (r *Repository) FindJobLogs(ctx context.Context, offset, limit int) ([]SysJobLog, int64, error) {
	var logs []SysJobLog
	var total int64

	if err := r.db.WithContext(ctx).Model(&SysJobLog{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).Order("job_log_id DESC").Offset(offset).Limit(limit).Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}
