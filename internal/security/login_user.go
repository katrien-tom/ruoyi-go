package security

import "time"

type LoginUser struct {
	UserID      int64     `json:"userId"`
	DeptID      *int64    `json:"deptId,omitempty"`
	Token       string    `json:"token"`
	LoginTime   int64     `json:"loginTime"`
	ExpireTime  int64     `json:"expireTime"`
	IPAddr      string    `json:"ipaddr,omitempty"`
	LoginIP     string    `json:"loginIp,omitempty"`
	Browser     string    `json:"browser,omitempty"`
	OS          string    `json:"os,omitempty"`
	Roles       []string  `json:"roles,omitempty"`
	Permissions []string  `json:"permissions,omitempty"`
	User        *UserInfo `json:"user"`
}

type UserInfo struct {
	UserID      int64   `json:"userId"`
	DeptID      *int64  `json:"deptId,omitempty"`
	UserName    string  `json:"userName"`
	NickName    string  `json:"nickName"`
	UserType    string  `json:"userType"`
	Email       string  `json:"email"`
	Phonenumber string  `json:"phonenumber"`
	Sex         string  `json:"sex"`
	Avatar      string  `json:"avatar"`
	Status      string  `json:"status"`
	DelFlag     string  `json:"delFlag"`
	Remark      *string `json:"remark,omitempty"`
}

func (u *LoginUser) LoginAt() time.Time {
	return time.UnixMilli(u.LoginTime)
}

func (u *LoginUser) ExpireAt() time.Time {
	return time.UnixMilli(u.ExpireTime)
}
