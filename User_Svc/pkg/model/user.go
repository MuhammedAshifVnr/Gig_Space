package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	FirstName string
	LastName  string
	Email     string
	Password  string
	Country   string
	Phone     string
	Role      string
	IsActive  bool
}

type UserProfile struct {
	gorm.Model
	UserID uint
	User   User
	Bio    string
	Title  string
}

type Freelancer_Skills struct {
	gorm.Model
	UserID           uint   `gorm:"not null;uniqueIndex:freelancer_skills_unique"`
	User             User   `gorm:"foreignKey:UserID"`
	SkillID          uint   `gorm:"not null;uniqueIndex:freelancer_skills_unique"`
	Skills           Skills `gorm:"foreignKey:SkillID"`
	ProficiencyLevel int
}
type ProfilePhoto struct {
	gorm.Model
	UserID uint
	User   User
	Photo  string
}

type Address struct {
	gorm.Model
	User_id  uint
	User     User
	State    string
	District string
	City     string
}
