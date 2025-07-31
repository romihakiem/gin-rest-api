package models

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	Title      string   `gorm:"not null" json:"title"`
	Body       string   `gorm:"type:text" json:"body"`
	CategoryID uint     `gorm:"foreignkey:CategoryID" json:"categoryID"`
	UserID     uint     `gorm:"foreignkey:UserID" json:"userID"`
	Category   Category `gorm:"foreignkey:CategoryID"`
	User       User     `gorm:"foreignkey:UserID"`
	Comments   []Comment
}
