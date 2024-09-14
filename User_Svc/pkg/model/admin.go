package model

import "gorm.io/gorm"

type Admin struct {
	gorm.Model
	Name     string
	Email    string
	Password string
}

type Category struct {
	gorm.Model
	Name     string `gorm:"unique;not null" json:"name"`
	IsActive bool
}

type Skills struct{
	gorm.Model
	SkillName string `gorm:"unique;not null" json:"name"`
}