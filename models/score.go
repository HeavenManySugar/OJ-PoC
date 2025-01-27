package models

import "gorm.io/gorm"

type Score struct {
	gorm.Model
	GitRepo string `json:"git_repo" example:"user_name/repo_name" validate:"required"`
	Score int `json:"score" example:"100" validate:"required"`
}