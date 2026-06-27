package models

// AdminUser represents admin panel user accounts
type AdminUser struct {
	BaseModel
	Username     string `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Email        string `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	PasswordHash string `gorm:"type:varchar(255);not null" json:"-"`
	FullName     string `gorm:"type:varchar(100)" json:"full_name,omitempty"`
	RoleID       string `gorm:"type:uuid" json:"role_id,omitempty"`
	IsActive     bool   `gorm:"default:true" json:"is_active"`
	LastLogin    string `json:"last_login,omitempty"`
}

func (AdminUser) TableName() string {
	return "admin_panel.admin_users"
}
