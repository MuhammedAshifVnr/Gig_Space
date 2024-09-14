package di

import (
	"github.com/MuhammedAshifVnr/Gig_Space/Gig_Svc/pkg/db"
	"github.com/MuhammedAshifVnr/Gig_Space/Gig_Svc/pkg/internal/repo"
	"github.com/MuhammedAshifVnr/Gig_Space/Gig_Svc/pkg/internal/service"
	"github.com/MuhammedAshifVnr/Gig_Space/Gig_Svc/utils/client"
	"github.com/spf13/viper"
)

func InitializeService() *service.GigService {
	db,s3Svc := db.InitializeDB(viper.GetString("DSN"))
	repo := repo.NewGigRepository(db)
	userClient:=client.NewUserClinet()
	service := service.NewGigService(repo,s3Svc,userClient)

	return service
}
