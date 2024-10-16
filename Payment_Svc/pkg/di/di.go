package di

import (
	"github.com/MuhammedAshifVnr/Gig_Space/Payment_Svc/pkg/db"
	"github.com/MuhammedAshifVnr/Gig_Space/Payment_Svc/pkg/internal/repo"
	"github.com/MuhammedAshifVnr/Gig_Space/Payment_Svc/pkg/internal/service"
	"github.com/MuhammedAshifVnr/Gig_Space/Payment_Svc/utils/client"
	"github.com/spf13/viper"
)

func InitializeService() *service.PaymentService {
	db := db.InitializeDB(viper.GetString("DBUrl"))
	repo := repo.NewPaymentRepository(db)
	user:=client.NewUserClinet()
	service := service.NewPaymentService(repo,user)
	return service
}
