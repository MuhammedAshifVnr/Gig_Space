package repo

import (
	"github.com/MuhammedAshifVnr/Gig_Space/User_svc/pkg/model"
	"github.com/MuhammedAshifVnr/Gig_Space/User_svc/utils/helper"
	"github.com/MuhammedAshifVnr/Gig_Space_Proto/proto"
)

type RepoInter interface {
	CheckingExist(email, phone string) error
	SignupData(data model.User, otp string) error
	VerifyingEmail(otp, email string) (string, error)
	CreateUser(user model.User) error
	GetUser(email string) (model.User, error)
	GetAdmin(email string) (model.Admin, error)
	AddCategory(category model.Category) error
	AddSkill(skill model.Skills) error
	GetSkill(name string) (model.Skills, error)
	FreelacerAddSkill(userID uint, skillID uint, proficiency int) error
	GetBio(userID uint) (model.UserProfile, error)
	CreateBio(bio model.UserProfile) error
	UpdateBio(bio model.UserProfile) error
	GetProfile(userID uint) (helper.UserWithProfile, error)
	GetSkillByuserID(userID uint) ([]*proto.UserSkill, error)
	DeleteSkillByID(userID, skillID uint) error
	GetCategory() ([]*proto.Category, error)
	GetSkills() ([]*proto.Skill, error)
	AdminDeleteSkill(id uint) error
	AdminDeleteCategory(id uint) error
	GetCategoryID(name string) (uint, error)
	GetAllUsers() ([]*proto.Profile, error)
	CreatPhoto(url string, id uint) error
	UpdatePhoto(url string, id uint) error
	GetPhoto(id uint) (model.ProfilePhoto, error)
	UnBlockUser(id uint)error
	BlockUser(id uint)error
	GetUserByID(id uint) (model.User, error)
	ResetPassword(email, password string) error
	DeleteOtp(otp string)
	UpdateAddress(add model.Address) error
	CreateAddress(State, District, City string, id int) error
	GetAddress(id uint) model.Address
	RoleChange(role string, id uint) error
	GetUserEmail(userid uint)(string,error)
}
