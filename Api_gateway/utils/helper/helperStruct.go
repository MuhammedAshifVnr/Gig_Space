package helper

import "time"

type SignupData struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"useremail"`
	Password  string `json:"userpassword"`
	Country   string `json:"country"`
	Phone     string `json:"phone"`
}

type LoginData struct {
	Email    string `json:"useremail"`
	Password string `json:"userpassword"`
}

type ADLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AddCategory struct {
	Name string `json:"category_name"`
}

type AddSkill struct {
	Name string `json:"skill_name"`
}

type Skill struct {
	SkillName        string `json:"skill_name"`
	ProficiencyLevel int    `json:"proficiency_level"`
}

type AddSkillsRequest struct {
	Skills []Skill `json:"skills"`
}

type Message struct {
	SenderID    int32
	RecipientID int32
	MessageText string
	CreatedAt   time.Time
}
