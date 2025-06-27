package models

import "time"

// Link /**
type Link struct {
	ID        uint      `gorm:"primaryKey"`
	ShortCode string    `gorm:"uniqueIndex;unique;size:10;not null"`
	LongURL   string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime;not null"`
}
