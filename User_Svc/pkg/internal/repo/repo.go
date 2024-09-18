package repo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/MuhammedAshifVnr/Gig_Space/User_svc/pkg/model"
	"github.com/MuhammedAshifVnr/Gig_Space/User_svc/utils/helper"
	"github.com/MuhammedAshifVnr/Gig_Space_Proto/proto"
	"github.com/redis/go-redis/v9"

	"gorm.io/gorm"
)

var ctx = context.Background()

type UserRepo struct {
	DB  *gorm.DB
	RDB *redis.Client
}

func NewUserRepository(db *gorm.DB, rdb *redis.Client) *UserRepo {
	return &UserRepo{
		DB:  db,
		RDB: rdb,
	}
}

func (r *UserRepo) CheckingExist(email, phone string) error {
	var count int
	if err := r.DB.Raw("Select count(*) from users where email = ? or phone = ?", email, phone).Scan(&count); err.Error != nil {
		return err.Error
	}
	if count != 0 {
		return errors.New("user alredy exist")
	}
	return nil
}

func (r *UserRepo) SignupData(data model.User, otp string) error {
	userJSON, err := json.Marshal(data)
	if err != nil {
		fmt.Println("====")
		return err
	}
	err = r.RDB.Set(ctx, otp, userJSON, 120*time.Second).Err()
	if err != nil {
		fmt.Println("=--=")
		return err
	}
	return nil
}

func (r *UserRepo) VerifyingEmail(otp, email string) (string, error) {
	val, err := r.RDB.Get(ctx, otp).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

func (r *UserRepo) DeleteOtp(otp string) {
	r.RDB.Del(ctx, otp)
}

func (r *UserRepo) CreateUser(user model.User) error {
	err := r.DB.Create(&user)
	return err.Error
}

func (r *UserRepo) GetUser(email string) (model.User, error) {
	var user model.User

	err := r.DB.First(&user, "email=?", email)
	if err.Error != nil {
		return model.User{}, err.Error
	}
	fmt.Println("user=", user)
	return user, nil
}

func (r *UserRepo) GetSkill(name string) (model.Skills, error) {
	var skill model.Skills
	query := `SELECT * FROM skills WHERE skill_name ILIKE ?`
	err := r.DB.Raw(query, name).Scan(&skill).Error
	if err != nil {
		return skill, err
	}
	return skill, nil
}

func (r *UserRepo) FreelacerAddSkill(userID uint, skillID uint, proficiency int) error {
	skill := model.Freelancer_Skills{
		UserID:           userID,
		SkillID:          skillID,
		ProficiencyLevel: proficiency,
	}
	err := r.DB.Create(&skill)
	if err != nil {
		return err.Error
	}
	return nil
}

func (r *UserRepo) UpdateBio(bio model.UserProfile) error {
	query := `UPDATE user_profiles SET bio = ?, title = ? WHERE user_id = ?`
	err := r.DB.Raw(query, bio.Bio, bio.Title, bio.UserID).Error
	return err
}

func (r *UserRepo) CreateBio(bio model.UserProfile) error {
	err := r.DB.Create(&bio).Error
	return err
}

func (r *UserRepo) GetBio(userID uint) (model.UserProfile, error) {
	query := `select * from userProfile where userId= ?`
	var user model.UserProfile
	err := r.DB.Raw(query, userID).Scan(&user)
	return user, err.Error
}

func (r *UserRepo) GetProfile(userID uint) (helper.UserWithProfile, error) {
	var user helper.UserWithProfile
	query := `SELECT u.id as user_id, u.first_name, u.last_name, u.email,u.country,u.phone, up.bio, up.title, pp.photo as profile_photo
        FROM users u
        LEFT JOIN user_profiles up ON u.id = up.user_id
        LEFT JOIN profile_photos pp ON u.id = pp.user_id
        WHERE u.id = ?`
	err := r.DB.Raw(query, userID).Scan(&user).Error
	return user, err
}

func (r *UserRepo) GetSkillByuserID(userID uint) ([]*proto.UserSkill, error) {
	var userSkills []*proto.UserSkill

	query := `
    SELECT fs.id ,s.skill_name, fs.proficiency_level
    FROM freelancer_skills fs
    JOIN skills s ON fs.skill_id = s.id
    WHERE fs.user_id = ?
`
	err := r.DB.Raw(query, userID).Scan(&userSkills).Error
	return userSkills, err
}

func (r *UserRepo) DeleteSkillByID(userID, skillID uint) error {
	query := `DELETE FROM freelancer_skills WHERE user_id = ? AND skill_id = ?`

	result := r.DB.Exec(query, userID, skillID)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("no matching skill found to delete")
	}

	return nil
}

func (r *UserRepo) GetCategory() ([]*proto.Category, error) {
	var category []*proto.Category
	query := `select id,name from categories`
	err := r.DB.Raw(query).Scan(&category).Error
	return category, err
}

func (r *UserRepo) GetSkills() ([]*proto.Skill, error) {
	var skills []*proto.Skill
	query := `select id ,  skill_name from skills`
	err := r.DB.Raw(query).Scan(&skills).Error
	return skills, err
}

func (r *UserRepo) GetCategoryID(name string) (uint, error) {
	var id uint
	query := `select id from categories where name =?`
	err := r.DB.Raw(query, name).Scan(&id)
	if err.Error != nil {
		return id, err.Error
	}

	if err.RowsAffected == 0 {
		return id, fmt.Errorf("enter a vaild category")
	}

	return id, nil
}

func (r *UserRepo) GetAllUsers() ([]*proto.Profile, error) {
	var users []*proto.Profile
	query := `SELECT 
    id ,
    first_name,
    last_name ,
    email ,
    country ,
    phone ,
	is_active
FROM users`
	err := r.DB.Raw(query).Scan(&users)
	return users, err.Error
}

func (r *UserRepo) CreatPhoto(url string, id uint) error {
	profile := model.ProfilePhoto{
		UserID: id,
		Photo:  url,
	}
	err := r.DB.Create(&profile)
	return err.Error
}

func (r *UserRepo) GetPhoto(id uint) (model.ProfilePhoto, error) {
	query := `select * from profile_photos where user_id =?`
	var profile model.ProfilePhoto
	err := r.DB.Raw(query, id).Scan(&profile)
	return profile, err.Error
}

func (r *UserRepo) UpdatePhoto(url string, id uint) error {
	query := `UPDATE user_profiles SET photo =? WHERE user_id = ?`
	err := r.DB.Raw(query, url, id)
	return err.Error
}

func (r *UserRepo) ResetPassword(email, password string) error {
	query := `UPDATE users SET password = ? WHERE email = ?`
	err := r.DB.Exec(query, password, email)
	return err.Error
}

func (r *UserRepo) GetAddress(id uint) model.Address {
	var add model.Address
	query := `select * from addresses where user_id = ?`
	r.DB.Raw(query, id).Scan(&add)
	return add
}

func (r *UserRepo) CreateAddress(State, District, City string, id int) error {
	query := `INSERT INTO addresses (state, district, city, user_id,created_at,updated_at) VALUES ($1, $2, $3, $4, $5,$6)`
	err := r.DB.Exec(query, State, District, City, id, time.Now(), time.Now()).Error
	return err
}

func (r *UserRepo) UpdateAddress(add model.Address) error {
	query := `UPDATE addresses SET state = $1,district =$2,city=$3,updated_at =$5 WHERE id = $4`
	err := r.DB.Exec(query, add.State, add.District, add.City, add.ID, time.Now()).Error
	return err
}
