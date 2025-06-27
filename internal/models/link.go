package models

// Link /**
type Link struct {
	ID        uint   `gorm:"primaryKey"`
	shortcode string `gorm:"uniqueIndex;unique;size:10;not null"`
	longURL   string `gorm:"not null"`
	createdAt string `gorm:"autoCreateTime;not null"`
}
