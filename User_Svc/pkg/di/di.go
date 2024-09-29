package di

import (
	"github.com/MuhammedAshifVnr/Gig_Space/User_svc/pkg/config"
	"github.com/MuhammedAshifVnr/Gig_Space/User_svc/pkg/db"
	"github.com/MuhammedAshifVnr/Gig_Space/User_svc/pkg/internal/repo"
	"github.com/MuhammedAshifVnr/Gig_Space/User_svc/pkg/internal/service"
	"github.com/MuhammedAshifVnr/Gig_Space/User_svc/utils/client"
)

func InitializeService(cfg *config.Config) *service.UserService {
	db ,rdb,s3:= db.InitializeDB(cfg.DBUrl)
	repo := repo.NewUserRepository(db,rdb)
	payment:=client.NewPaymentClinet()
	service := service.NewUserService(repo,s3,payment)

	return service
}
