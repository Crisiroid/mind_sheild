package models

import "time"

type AdminUser struct {
	BaseModel
	Username     string `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Email        string `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	PasswordHash string `gorm:"type:varchar(255);not null" json:"-"`
	FullName     string `gorm:"type:varchar(100)" json:"full_name,omitempty"`
	RoleID       string `gorm:"type:uuid" json:"role_id,omitempty"`
	IsActive     bool   `gorm:"default:true" json:"is_active"`
	LastLogin    string `json:"last_login,omitempty"`

	RefreshToken       string     `gorm:"type:varchar(500)" json:"-"`
	RefreshTokenExpiry *time.Time `json:"-"`
}

func (AdminUser) TableName() string {
	return "admin_panel.admin_users"
}
