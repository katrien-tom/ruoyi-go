package user

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/banyejiu/ruoyi-go/pkg/constants"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindByUserName(ctx context.Context, userName string) (*SysUser, error) {
	var u SysUser
	if err := r.db.WithContext(ctx).
		Where("user_name = ? AND del_flag = ?", userName, constants.DelFlagNormal).
		First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *Repository) FindByID(ctx context.Context, userID int64) (*SysUser, error) {
	var u SysUser
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND del_flag = ?", userID, constants.DelFlagNormal).
		First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *Repository) FindPermissionsByUserID(ctx context.Context, userID int64) ([]string, error) {
	var permissions []string
	err := r.db.WithContext(ctx).
		Table("sys_menu m").
		Select("DISTINCT m.perms").
		Joins("JOIN sys_role_menu rm ON rm.menu_id = m.menu_id").
		Joins("JOIN sys_user_role ur ON ur.role_id = rm.role_id").
		Where("ur.user_id = ? AND m.status = ? AND m.perms <> '' AND m.perms IS NOT NULL", userID, constants.StatusNormal).
		Scan(&permissions).Error
	if err != nil {
		return nil, err
	}
	return permissions, nil
}

func (r *Repository) FindRoleKeysByUserID(ctx context.Context, userID int64) ([]string, error) {
	var roleKeys []string
	err := r.db.WithContext(ctx).
		Table("sys_role r").
		Select("DISTINCT r.role_key").
		Joins("JOIN sys_user_role ur ON ur.role_id = r.role_id").
		Where("ur.user_id = ? AND r.status = ? AND r.del_flag = ? AND r.role_key <> ''", userID, constants.StatusNormal, constants.DelFlagNormal).
		Scan(&roleKeys).Error
	if err != nil {
		return nil, err
	}
	return roleKeys, nil
}

type userListFilter struct {
	UserName    string
	Phonenumber string
	Status      string
	DeptID      *int64
	BeginTime   string
	EndTime     string
}

func (r *Repository) FindAll(ctx context.Context, offset, limit int, filter userListFilter) ([]UserResponse, int64, error) {
	var total int64

	base := r.db.WithContext(ctx).
		Table("sys_user u").
		Select("u.*, d.dept_name").
		Joins("LEFT JOIN sys_dept d ON u.dept_id = d.dept_id").
		Where("u.del_flag = ?", constants.DelFlagNormal)

	if filter.UserName != "" {
		base = base.Where("u.user_name LIKE ?", "%"+filter.UserName+"%")
	}
	if filter.Phonenumber != "" {
		base = base.Where("u.phonenumber LIKE ?", "%"+filter.Phonenumber+"%")
	}
	if filter.Status != "" {
		base = base.Where("u.status = ?", filter.Status)
	}
	if filter.DeptID != nil {
		base = base.Where("u.dept_id = ?", *filter.DeptID)
	}
	if filter.BeginTime != "" {
		base = base.Where("u.create_time >= ?", filter.BeginTime)
	}
	if filter.EndTime != "" {
		base = base.Where("u.create_time <= ?", filter.EndTime)
	}

	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	type rawRow struct {
		SysUser
		DeptName string `gorm:"column:dept_name"`
	}

	var rows []rawRow
	if err := base.Order("u.create_time DESC").Offset(offset).Limit(limit).Scan(&rows).Error; err != nil {
		return nil, 0, err
	}

	result := make([]UserResponse, len(rows))
	for i, row := range rows {
		result[i] = UserResponse{
			UserID:      row.UserID,
			DeptID:      row.DeptID,
			DeptName:    row.DeptName,
			UserName:    row.UserName,
			NickName:    row.NickName,
			UserType:    row.UserType,
			Email:       row.Email,
			Phonenumber: row.Phonenumber,
			Sex:         row.Sex,
			Avatar:      row.Avatar,
			Status:      row.Status,
			DelFlag:     row.DelFlag,
			LoginIP:     row.LoginIP,
			LoginDate:   row.LoginDate,
			CreateBy:    row.CreateBy,
			CreateTime:  row.CreateTime,
			UpdateBy:    row.UpdateBy,
			UpdateTime:  row.UpdateTime,
			Remark:      row.Remark,
		}
	}

	return result, total, nil
}

func (r *Repository) Create(ctx context.Context, u *SysUser) error {
	return r.db.WithContext(ctx).Create(u).Error
}

func (r *Repository) Update(ctx context.Context, u *SysUser) error {
	return r.db.WithContext(ctx).Model(u).Where("user_id = ?", u.UserID).Updates(u).Error
}

func (r *Repository) SoftDelete(ctx context.Context, userIDs []int64) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&SysUser{}).
		Where("user_id IN ?", userIDs).
		Updates(map[string]any{
			"del_flag":    constants.DelFlagDeleted,
			"update_time": now,
		}).Error
}

func (r *Repository) FindUserRoles(ctx context.Context, userID int64) ([]SysUserRole, error) {
	var userRoles []SysUserRole
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&userRoles).Error; err != nil {
		return nil, err
	}
	return userRoles, nil
}

func (r *Repository) FindUserPosts(ctx context.Context, userID int64) ([]SysUserPost, error) {
	var userPosts []SysUserPost
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&userPosts).Error; err != nil {
		return nil, err
	}
	return userPosts, nil
}

func (r *Repository) ReplaceUserRoles(ctx context.Context, userID int64, roleIDs []int64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ?", userID).Delete(&SysUserRole{}).Error; err != nil {
			return err
		}
		if len(roleIDs) == 0 {
			return nil
		}
		userRoles := make([]SysUserRole, len(roleIDs))
		for i, rid := range roleIDs {
			userRoles[i] = SysUserRole{UserID: userID, RoleID: rid}
		}
		return tx.Create(&userRoles).Error
	})
}

func (r *Repository) ReplaceUserPosts(ctx context.Context, userID int64, postIDs []int64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ?", userID).Delete(&SysUserPost{}).Error; err != nil {
			return err
		}
		if len(postIDs) == 0 {
			return nil
		}
		userPosts := make([]SysUserPost, len(postIDs))
		for i, pid := range postIDs {
			userPosts[i] = SysUserPost{UserID: userID, PostID: pid}
		}
		return tx.Create(&userPosts).Error
	})
}

func (r *Repository) ExistsByUserNameExcluding(ctx context.Context, userName string, excludeID int64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&SysUser{}).
		Where("user_name = ? AND del_flag = ? AND user_id <> ?", userName, constants.DelFlagNormal, excludeID).
		Count(&count).Error
	return count > 0, err
}

func (r *Repository) ExistsByPhonenumberExcluding(ctx context.Context, phonenumber string, excludeID int64) (bool, error) {
	if phonenumber == "" {
		return false, nil
	}
	var count int64
	err := r.db.WithContext(ctx).
		Model(&SysUser{}).
		Where("phonenumber = ? AND del_flag = ? AND user_id <> ?", phonenumber, constants.DelFlagNormal, excludeID).
		Count(&count).Error
	return count > 0, err
}

func (r *Repository) ExistsByEmailExcluding(ctx context.Context, email string, excludeID int64) (bool, error) {
	if email == "" {
		return false, nil
	}
	var count int64
	err := r.db.WithContext(ctx).
		Model(&SysUser{}).
		Where("email = ? AND del_flag = ? AND user_id <> ?", email, constants.DelFlagNormal, excludeID).
		Count(&count).Error
	return count > 0, err
}

func (r *Repository) UpdateStatus(ctx context.Context, userID int64, status string) error {
	return r.db.WithContext(ctx).
		Model(&SysUser{}).
		Where("user_id = ?", userID).
		Update("status", status).Error
}

func (r *Repository) UpdatePassword(ctx context.Context, userID int64, hashedPwd string) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&SysUser{}).
		Where("user_id = ?", userID).
		Updates(map[string]any{
			"password":        hashedPwd,
			"pwd_update_date": now,
			"update_time":     now,
		}).Error
}

func (r *Repository) FindDeptByID(ctx context.Context, deptID int64) (*struct {
	DeptID   int64  `gorm:"column:dept_id"`
	DeptName string `gorm:"column:dept_name"`
}, error) {
	var result struct {
		DeptID   int64  `gorm:"column:dept_id"`
		DeptName string `gorm:"column:dept_name"`
	}
	err := r.db.WithContext(ctx).
		Table("sys_dept").
		Select("dept_id, dept_name").
		Where("dept_id = ? AND del_flag = ?", deptID, constants.DelFlagNormal).
		First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *Repository) FindAllRoles(ctx context.Context) ([]RoleInfo, error) {
	var roles []RoleInfo
	err := r.db.WithContext(ctx).
		Table("sys_role").
		Select("role_id, role_name, role_key").
		Where("status = ? AND del_flag = ?", constants.StatusNormal, constants.DelFlagNormal).
		Order("role_sort").
		Scan(&roles).Error
	return roles, err
}

func (r *Repository) FindRoleIDsByUserID(ctx context.Context, userID int64) ([]int64, error) {
	var ids []int64
	err := r.db.WithContext(ctx).
		Table("sys_user_role").
		Select("role_id").
		Where("user_id = ?", userID).
		Scan(&ids).Error
	return ids, err
}

func (r *Repository) FindPostIDsByUserID(ctx context.Context, userID int64) ([]int64, error) {
	var ids []int64
	err := r.db.WithContext(ctx).
		Table("sys_user_post").
		Select("post_id").
		Where("user_id = ?", userID).
		Scan(&ids).Error
	return ids, err
}

// SelectPostList returns all available posts as {postId, postName} pairs.
func (r *Repository) SelectPostList(ctx context.Context) ([]map[string]any, error) {
	var result []map[string]any
	err := r.db.WithContext(ctx).
		Table("sys_post").
		Select("post_id as value, post_name as label").
		Where("status = ?", constants.StatusNormal).
		Order("post_sort").
		Scan(&result).Error
	return result, err
}

// CheckUserDataScope checks if the login user has permission to manage the target user's data.
// For system:user:manage scope, allow all. Otherwise only allow managing own created users.
func (r *Repository) CheckUserDataScope(ctx context.Context, targetUserID, loginUserID int64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&SysUser{}).
		Where("user_id = ? AND del_flag = ?", targetUserID, constants.DelFlagNormal).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	if count == 0 {
		return false, fmt.Errorf("user not found: %d", targetUserID)
	}

	// Disallow operating on self
	if targetUserID == loginUserID {
		return false, fmt.Errorf("cannot operate on yourself")
	}

	return true, nil
}
