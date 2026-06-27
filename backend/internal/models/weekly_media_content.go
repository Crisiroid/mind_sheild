package models

type WeeklyMediaContent struct {
	BaseModel
	FileName       string `gorm:"size:255;not null" json:"file_name"`
	FileType       string `gorm:"size:50;not null" json:"file_type"`
	FileSize       int64  `gorm:"not null" json:"file_size"`
	FileURL        string `gorm:"size:500;not null" json:"file_url"`
	WeekNumber     int    `gorm:"not null;check:week_number BETWEEN 1 AND 52" json:"week_number"`
	Description    string `gorm:"type:text" json:"description,omitempty"`
	ContentType    string `gorm:"size:100;not null" json:"content_type"`
	OriginalName   string `gorm:"size:255;not null" json:"original_name"`
	StoragePath    string `gorm:"size:500;not null" json:"storage_path"`
	IsActive       bool   `gorm:"default:true" json:"is_active"`
	DownloadCount  int    `gorm:"default:0" json:"download_count"`
	CreatedByAdmin string `gorm:"type:uuid" json:"created_by_admin,omitempty"`
}

func (WeeklyMediaContent) TableName() string {
	return "weekly_media_contents"
}
