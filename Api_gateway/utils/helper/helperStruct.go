package helper

import "time"

type SignupData struct {
	Firstname string `json:"firstname" validate:"required,alpha,min=4,max=50"`
	Lastname  string `json:"lastname" validate:"required,alpha,min=4,max=50"`
	Email     string `json:"useremail" validate:"required,email"`
	Password  string `json:"userpassword" validate:"required,min=6"`
	Country   string `json:"country" validate:"required,alpha"`
	Phone     string `json:"phone" validate:"required,numeric,len=10"` 
}

type LoginData struct {
	Email    string `json:"useremail" validate:"required,email"`
	Password string `json:"userpassword" validate:"required,min=4"`
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

type PaymentAuthorizedEvent struct {
	Event  string `json:"event"`
	Entity struct {
		ID       string `json:"id"`
		Amount   int    `json:"amount"`
		Currency string `json:"currency"`
		// Add other required fields from the webhook payload
	} `json:"entity"`
	// Add other required fields from the webhook payload
}