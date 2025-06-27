package models

// Link /**
type Link struct {
	ID        uint   `gorm:"primaryKey"`
	ShortCode string `gorm:"uniqueIndex;unique;size:10;not null"`
	LongURL   string `gorm:"not null"`
	CreatedAt string `gorm:"autoCreateTime;not null"`
}
