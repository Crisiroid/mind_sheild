package models

type AdminRole struct {
	BaseModel
	RoleName    string `gorm:"type:varchar(50);uniqueIndex;not null" json:"role_name"`
	Description string `gorm:"type:text" json:"description,omitempty"`
	Permissions string `gorm:"type:jsonb" json:"permissions,omitempty"`
}

func (AdminRole) TableName() string {
	return "admin_panel.roles"
}
