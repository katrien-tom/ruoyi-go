package dept

import (
	"context"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindAll(ctx context.Context) ([]SysDept, error) {
	var depts []SysDept
	err := r.db.WithContext(ctx).Where("del_flag = ?", "0").Order("parent_id, order_num").Find(&depts).Error
	return depts, err
}

func (r *Repository) FindByID(ctx context.Context, id int64) (*SysDept, error) {
	var dept SysDept
	if err := r.db.WithContext(ctx).Where("dept_id = ? AND del_flag = ?", id, "0").First(&dept).Error; err != nil {
		return nil, err
	}
	return &dept, nil
}

func (r *Repository) Create(ctx context.Context, d *SysDept) error {
	return r.db.WithContext(ctx).Create(d).Error
}

func (r *Repository) Update(ctx context.Context, d *SysDept) error {
	return r.db.WithContext(ctx).Model(d).Where("dept_id = ?", d.DeptID).Updates(d).Error
}

func (r *Repository) UpdateChildrenAncestors(ctx context.Context, oldAncestor, newAncestor string) error {
	return r.db.WithContext(ctx).
		Model(&SysDept{}).
		Where("ancestors LIKE ?", oldAncestor+"%").
		UpdateColumn("ancestors", gorm.Expr(
			"CASE WHEN ancestors = ? THEN ? ELSE REPLACE(ancestors, ?, ?) END",
			oldAncestor, newAncestor, oldAncestor, newAncestor,
		)).Error
}

func (r *Repository) Delete(ctx context.Context, deptID int64) error {
	return r.db.WithContext(ctx).Model(&SysDept{}).Where("dept_id = ?", deptID).Update("del_flag", "2").Error
}

func (r *Repository) HasChildren(ctx context.Context, deptID int64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&SysDept{}).Where("parent_id = ? AND del_flag = ?", deptID, "0").Count(&count).Error
	return count > 0, err
}

func (r *Repository) ExistsUsersInDept(ctx context.Context, deptID int64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Table("sys_user").
		Where("dept_id = ? AND del_flag = ?", deptID, "0").
		Count(&count).Error
	return count > 0, err
}

func (r *Repository) ExistsByNameExcluding(ctx context.Context, parentID int64, name string, excludeID int64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&SysDept{}).
		Where("parent_id = ? AND dept_name = ? AND del_flag = ? AND dept_id <> ?", parentID, name, "0", excludeID).
		Count(&count).Error
	return count > 0, err
}

func (r *Repository) FindAncestors(ctx context.Context, parentID int64) (string, error) {
	if parentID == 0 {
		return "0", nil
	}
	var dept SysDept
	if err := r.db.WithContext(ctx).Select("ancestors, dept_id").Where("dept_id = ?", parentID).First(&dept).Error; err != nil {
		return "", err
	}
	return fmt.Sprintf("%s,%d", dept.Ancestors, dept.DeptID), nil
}

func (r *Repository) BatchUpdateDeptAncestors(ctx context.Context, dept *SysDept, tx *gorm.DB) error {
	var children []SysDept
	if err := tx.WithContext(ctx).Where("parent_id = ?", dept.DeptID).Find(&children).Error; err != nil {
		return err
	}
	newAncestors := fmt.Sprintf("%s,%d", dept.Ancestors, dept.DeptID)
	for i := range children {
		children[i].Ancestors = newAncestors
		if err := r.BatchUpdateDeptAncestors(ctx, &children[i], tx); err != nil {
			return err
		}
	}
	return tx.Model(dept).Where("dept_id = ?", dept.DeptID).Update("ancestors", newAncestors).Error
}

func (r *Repository) WithTransaction(fn func(tx *Repository) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		txRepo := &Repository{db: tx}
		return fn(txRepo)
	})
}

func (r *Repository) DB() *gorm.DB {
	return r.db
}

// FindAllForExport returns all depts for data export (no pagination).
func (r *Repository) FindAllForExport(ctx context.Context) ([]SysDept, error) {
	return r.FindAll(ctx)
}

// NormalizeAncestors converts "0,100,200" -> "0".
// When dept parent is 0, ancestors should just be "0".
func NormalizeAncestors(parentID int64, parentAncestors string) string {
	if parentID == 0 {
		return "0"
	}
	return strings.TrimRight(parentAncestors+","+fmt.Sprint(parentID), ",")
}
