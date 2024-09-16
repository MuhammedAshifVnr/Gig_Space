package repo

import (
	"fmt"

	"github.com/MuhammedAshifVnr/Gig_Space/User_svc/pkg/model"
)

func (r *UserRepo) GetAdmin(email string) (model.Admin, error) {
	var admin model.Admin
	query := `SELECT * FROM admins WHERE email = $1`
	err := r.DB.Raw(query, email).Scan(&admin).Error

	if err != nil {
		return model.Admin{}, err
	}

	if admin.Email == "" {
		return model.Admin{}, fmt.Errorf("admin with email %s not found", email)
	}
	fmt.Println("admin = ", admin)
	return admin, nil
}

func (r *UserRepo) AddCategory(category model.Category) error {
	if err := r.DB.Create(&category).Error; err != nil {
		return err
	}
	return nil
}

func (r *UserRepo) AddSkill(skill model.Skills) error {
	if err := r.DB.Create(&skill).Error; err != nil {
		return err
	}
	return nil
}

func (r *UserRepo) AdminDeleteSkill(id uint) error {
	query := `Delete FROM skills WHERE id = ?`
	err := r.DB.Exec(query, id)
	if err.RowsAffected == 0 {
		return fmt.Errorf("no matching skill found to delete")
	}
	return err.Error
}

func (r *UserRepo) AdminDeleteCategory(id uint) error {
	fmt.Println("id ", id)
	query := `Delete FROM categories WHERE id = ?`
	err := r.DB.Exec(query, id)

	if err.RowsAffected == 0 {
		return fmt.Errorf("no matching category found to delete")
	}
	return err.Error
}

func (r *UserRepo) GetUserByID(id uint) (model.User, error) {
	var user model.User
	err := r.DB.First(&user, "id=?", id)
	if err.Error != nil {
		return model.User{}, err.Error
	}
	fmt.Println("user=", user)
	return user, nil
}

func (r *UserRepo) BlockUser(id uint) error {
	query := `UPDATE users SET is_active =? WHERE id = ?`
	err := r.DB.Exec(query, false, id)
	return err.Error
}

func (r *UserRepo) UnBlockUser(id uint) error {
	query := `UPDATE users SET is_active =? WHERE id = ?`
	err := r.DB.Exec(query, true, id)
	return err.Error
}
