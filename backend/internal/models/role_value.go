package models

// RoleAndValue represents organizational roles and personal values
type RoleAndValue struct {
	BaseModel
	UserID      string `gorm:"type:uuid;not null" json:"user_id"`
	EntryType   string `gorm:"type:varchar(20);check:entry_type IN ('role', 'value');not null" json:"entry_type"`
	EntryText   string `gorm:"type:text;not null" json:"entry_text"`
	CreatedDate string `gorm:"autoCreateTime" json:"created_date"`
	DayNumber   *int   `json:"day_number,omitempty"`
}

func (RoleAndValue) TableName() string {
	return "roles_and_values"
}
