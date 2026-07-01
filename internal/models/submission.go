package models

import "time"

type Submission struct {
	ID uint `json:"id" gorm:"primaryKey"`
	TaskID uint `json:"task_id" gorm:"index"`
	Language string `json:"language" gorm:"not null"`
	SourceCode string `json:"source_code" gorm:"type:text;not null"`
	Status string `json:"status" gorm:"default:'Pending'"`
	CreatedAt time.Time `json:"created_at"`
}
