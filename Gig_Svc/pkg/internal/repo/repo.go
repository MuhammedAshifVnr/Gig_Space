package repo

import (
	"fmt"

	"github.com/MuhammedAshifVnr/Gig_Space/Gig_Svc/pkg/internal/model"
	"github.com/MuhammedAshifVnr/Gig_Space_Proto/proto"

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

func (r *GigRepo) GetAllGigByID(id uint) ([]*proto.Gig, error) {
	var gigs []*proto.Gig
	query := "SELECT id, title, description, frelancer_id, price, category, delivery_days, revisions FROM gigs WHERE frelancer_id = $1"
	err := r.DB.Raw(query, id).Scan(&gigs).Error
	if err != nil {
		return gigs, err
	}
	fmt.Println("gig", gigs,"err",err)
	return gigs, nil
}
