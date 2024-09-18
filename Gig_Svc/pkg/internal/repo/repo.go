package repo

import (
	"fmt"

	"github.com/MuhammedAshifVnr/Gig_Space/Gig_Svc/pkg/internal/model"

	"gorm.io/gorm"
)

type GigRepo struct {
	DB *gorm.DB
}

func NewGigRepository(db *gorm.DB) *GigRepo {
	return &GigRepo{
		DB: db,
	}
}

func (r *GigRepo) CreateGgi(gig model.Gig) error {
	err := r.DB.Create(&gig).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *GigRepo) GetGigsByFreelancerID(freelancerID uint) ([]model.Gig, error) {
	var gigs []model.Gig
	query := r.DB.Where("freelancer_id = ?", freelancerID).Preload("Images")

	err := query.Find(&gigs).Error
	if err != nil {
		return nil, err
	}
	return gigs, nil
}

func (r *GigRepo) GetGigByID(Id uint) (model.Gig, error) {
	var gig model.Gig
	// query := `select * from gigs where id =?`
	// err := r.DB.Raw(query, Id).Scan(&gig)
	err := r.DB.Preload("Images").First(&gig, Id).Error
	return gig, err
}

func (r *GigRepo) UpdateGig(gig model.Gig) error {
	err := r.DB.Save(&gig).Error
	return err
}

func (r *GigRepo) DeleteImages(id uint) error {
	fmt.Println("id = ", id)
	query := `DELETE FROM images WHERE gig_id = ?`
	err := r.DB.Exec(query, id).Error
	return err
}

func (r *GigRepo)DeleteGig(id,user_id uint)error{
	return r.DB.Transaction(func(tx *gorm.DB) error {

		if err := tx.Exec(`DELETE FROM images WHERE gig_id = ?`, id).Error; err != nil {
			return err
		}

		if err := tx.Exec(`DELETE FROM gigs WHERE id = ? AND freelancer_id = ?`, id,user_id).Error; err != nil {
			return err
		}

		return nil
	})
}