package models

import "time"

type Task struct {
	ID uint `json:"id" gorm:"primaryKey"`
	Title string `json:"title" gorm:"not null"`
	Statements string `json:"statements" gorm:"not null"`
	TimeLimit int `json:"time_limit_ms"`
	MemoryLimit int `json:"memory_limit_mb"`
	TestCases []TestCase `json:"test_cases"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TestCase struct {
	ID uint `json:"id" gorm:"primaryKey"`
	TaskID uint `json:"task_id" gorm:"index;constraint:OnDelete:CASCADE;"`
	InputData string `json:"input_data" gorm:"type:text"`
	ExpectedOutput string `json:"expected_output" gorm:"type:text"`
	IsHidden bool `json:"is_hidden"`
}
