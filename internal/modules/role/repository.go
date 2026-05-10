package role

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

type roleListFilter struct {
	RoleName string
	RoleKey  string
	Status   string
}

func (r *Repository) FindAll(ctx context.Context, offset, limit int, filter roleListFilter) ([]SysRole, int64, error) {
	var total int64
	base := r.db.WithContext(ctx).Model(&SysRole{}).Where("del_flag = ?", "0")

	if filter.RoleName != "" {
		base = base.Where("role_name LIKE ?", "%"+filter.RoleName+"%")
	}
	if filter.RoleKey != "" {
		base = base.Where("role_key LIKE ?", "%"+filter.RoleKey+"%")
	}
	if filter.Status != "" {
		base = base.Where("status = ?", filter.Status)
	}

	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var roles []SysRole
	if err := base.Order("role_sort").Offset(offset).Limit(limit).Find(&roles).Error; err != nil {
		return nil, 0, err
	}
	return roles, total, nil
}

func (r *Repository) FindByID(ctx context.Context, id int64) (*SysRole, error) {
	var role SysRole
	if err := r.db.WithContext(ctx).Where("role_id = ? AND del_flag = ?", id, "0").First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *Repository) Create(ctx context.Context, role *SysRole) error {
	return r.db.WithContext(ctx).Create(role).Error
}

func (r *Repository) Update(ctx context.Context, role *SysRole) error {
	return r.db.WithContext(ctx).Model(role).Where("role_id = ?", role.RoleID).Updates(role).Error
}

func (r *Repository) SoftDelete(ctx context.Context, ids []int64) error {
	return r.db.WithContext(ctx).Model(&SysRole{}).Where("role_id IN ?", ids).Update("del_flag", "2").Error
}

func (r *Repository) UpdateStatus(ctx context.Context, roleID int64, status string) error {
	return r.db.WithContext(ctx).Model(&SysRole{}).Where("role_id = ?", roleID).Update("status", status).Error
}

func (r *Repository) ExistsByNameExcluding(ctx context.Context, name string, excludeID int64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&SysRole{}).
		Where("role_name = ? AND del_flag = ? AND role_id <> ?", name, "0", excludeID).
		Count(&count).Error
	return count > 0, err
}

func (r *Repository) ExistsByKeyExcluding(ctx context.Context, key string, excludeID int64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&SysRole{}).
		Where("role_key = ? AND del_flag = ? AND role_id <> ?", key, "0", excludeID).
		Count(&count).Error
	return count > 0, err
}

func (r *Repository) FindRoleMenuIDs(ctx context.Context, roleID int64) ([]int64, error) {
	var ids []int64
	err := r.db.WithContext(ctx).
		Table("sys_role_menu").
		Select("menu_id").
		Where("role_id = ?", roleID).
		Scan(&ids).Error
	return ids, err
}

func (r *Repository) ReplaceRoleMenus(ctx context.Context, roleID int64, menuIDs []int64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("role_id = ?", roleID).Delete(&SysRoleMenu{}).Error; err != nil {
			return err
		}
		if len(menuIDs) == 0 {
			return nil
		}
		rows := make([]SysRoleMenu, len(menuIDs))
		for i, mid := range menuIDs {
			rows[i] = SysRoleMenu{RoleID: roleID, MenuID: mid}
		}
		return tx.Create(&rows).Error
	})
}

func (r *Repository) FindRoleDeptIDs(ctx context.Context, roleID int64) ([]int64, error) {
	var ids []int64
	err := r.db.WithContext(ctx).
		Table("sys_role_dept").
		Select("dept_id").
		Where("role_id = ?", roleID).
		Scan(&ids).Error
	return ids, err
}

func (r *Repository) ReplaceRoleDepts(ctx context.Context, roleID int64, deptIDs []int64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("role_id = ?", roleID).Delete(&struct {
			RoleID int64 `gorm:"column:role_id"`
			DeptID int64 `gorm:"column:dept_id"`
		}{}).Error; err != nil {
			return err
		}
		if len(deptIDs) == 0 {
			return nil
		}
		type row struct {
			RoleID int64 `gorm:"column:role_id"`
			DeptID int64 `gorm:"column:dept_id"`
		}
		rows := make([]row, len(deptIDs))
		for i, did := range deptIDs {
			rows[i] = row{RoleID: roleID, DeptID: did}
		}
		return tx.Table("sys_role_dept").Create(&rows).Error
	})
}

func (r *Repository) FindAllocatedUsers(ctx context.Context, roleID int64, userName, phonenumber string, offset, limit int) ([]AuthUserResponse, int64, error) {
	var total int64
	base := r.db.WithContext(ctx).
		Table("sys_user u").
		Select("u.user_id, u.user_name, u.nick_name, u.phonenumber, u.status, u.dept_id, d.dept_name").
		Joins("JOIN sys_user_role ur ON u.user_id = ur.user_id").
		Joins("LEFT JOIN sys_dept d ON u.dept_id = d.dept_id").
		Where("ur.role_id = ? AND u.del_flag = ?", roleID, "0")

	if userName != "" {
		base = base.Where("u.user_name LIKE ?", "%"+userName+"%")
	}
	if phonenumber != "" {
		base = base.Where("u.phonenumber LIKE ?", "%"+phonenumber+"%")
	}

	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var rows []AuthUserResponse
	if err := base.Offset(offset).Limit(limit).Scan(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

func (r *Repository) FindUnallocatedUsers(ctx context.Context, roleID int64, userName, phonenumber string, offset, limit int) ([]AuthUserResponse, int64, error) {
	var total int64
	allocatedSub := r.db.WithContext(ctx).
		Table("sys_user_role ur").
		Select("ur.user_id").
		Where("ur.role_id = ?", roleID)

	base := r.db.WithContext(ctx).
		Table("sys_user u").
		Select("u.user_id, u.user_name, u.nick_name, u.phonenumber, u.status, u.dept_id, d.dept_name").
		Joins("LEFT JOIN sys_dept d ON u.dept_id = d.dept_id").
		Where("u.del_flag = ? AND u.user_id NOT IN (?)", "0", allocatedSub)

	if userName != "" {
		base = base.Where("u.user_name LIKE ?", "%"+userName+"%")
	}
	if phonenumber != "" {
		base = base.Where("u.phonenumber LIKE ?", "%"+phonenumber+"%")
	}

	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var rows []AuthUserResponse
	if err := base.Offset(offset).Limit(limit).Scan(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

func (r *Repository) CancelAuthUser(ctx context.Context, roleID, userID int64) error {
	return r.db.WithContext(ctx).
		Where("role_id = ? AND user_id = ?", roleID, userID).
		Delete(&struct {
			RoleID int64 `gorm:"column:role_id"`
			UserID int64 `gorm:"column:user_id"`
		}{}).Error
}

func (r *Repository) CancelAuthUsers(ctx context.Context, roleID int64, userIDs []int64) error {
	return r.db.WithContext(ctx).
		Table("sys_user_role").
		Where("role_id = ? AND user_id IN ?", roleID, userIDs).
		Delete(nil).Error
}

func (r *Repository) SelectAuthUsers(ctx context.Context, roleID int64, userIDs []int64) error {
	type row struct {
		RoleID int64 `gorm:"column:role_id"`
		UserID int64 `gorm:"column:user_id"`
	}
	rows := make([]row, len(userIDs))
	for i, uid := range userIDs {
		rows[i] = row{RoleID: roleID, UserID: uid}
	}
	return r.db.WithContext(ctx).Table("sys_user_role").Create(&rows).Error
}

func (r *Repository) FindAllRoles(ctx context.Context) ([]SysRole, error) {
	var roles []SysRole
	err := r.db.WithContext(ctx).Where("status = ? AND del_flag = ?", "0", "0").Order("role_sort").Find(&roles).Error
	return roles, err
}
