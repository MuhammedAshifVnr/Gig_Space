package helper

type UserWithProfile struct {
	UserID       uint   `json:"user_id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	Country      string `json:"country"`
	Phone        string `json:"phone"`
	Role         string `json:"role"`
	Bio          string `json:"bio"`
	Title        string `json:"title"`
	ProfilePhoto string `json:"profile_photo"`
}

type UserSkill struct {
	SkillName        string
	ProficiencyLevel int
}
