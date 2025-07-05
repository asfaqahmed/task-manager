package models

import "gorm.io/gorm"

type Task struct {
	gorm.Model
	Title    string
	Note     string
	Reminder string // ISO8601 string
	UserID   uint
}
