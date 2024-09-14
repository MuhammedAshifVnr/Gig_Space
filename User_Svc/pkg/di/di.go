package di

import (
	"github.com/MuhammedAshifVnr/Gig_Space/User_svc/pkg/config"
	"github.com/MuhammedAshifVnr/Gig_Space/User_svc/pkg/db"
	"github.com/MuhammedAshifVnr/Gig_Space/User_svc/pkg/internal/repo"
	"github.com/MuhammedAshifVnr/Gig_Space/User_svc/pkg/internal/service"
)

func InitializeService(cfg *config.Config) *service.UserService {
	db ,rdb:= db.InitializeDB(cfg.DBUrl)
	repo := repo.NewUserRepository(db,rdb)
	service := service.NewUserService(repo)

	return service
}
