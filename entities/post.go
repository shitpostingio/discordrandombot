package entities

import (
	"time"
)

//Post Entity
//Each `Post` belongs to a `User` (User, UserID)
//Each `Post` belongs to a `Type` (Type, TypeID)
//Each `Post` can have one `Fingerprint` (Fingerprint)
type Post struct {
	ID          uint   `gorm:"primary_key"`
	UserID      uint   `gorm:"not null"`
	MessageID   int    `gorm:"not null"`
	TypeID      uint   `gorm:"not null"`
	FileID      string `gorm:"type:varchar(190);unique;not null"`
	Caption     string
	PostedAt    *time.Time `gorm:"default:null"`
	HasError    bool       `gorm:"default:false"`
	CreatedAt   time.Time  `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt   *time.Time `gorm:"default:null"`
	DeletedAt   *time.Time `gorm:"default:null"`
}