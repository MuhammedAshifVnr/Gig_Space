package di

import (
	"github.com/MuhammedAshifVnr/Gig_Space/Payment_Svc/pkg/config"
	"github.com/MuhammedAshifVnr/Gig_Space/Payment_Svc/pkg/db"
	"github.com/MuhammedAshifVnr/Gig_Space/Payment_Svc/pkg/internal/repo"
	"github.com/MuhammedAshifVnr/Gig_Space/Payment_Svc/pkg/internal/service"
	"github.com/MuhammedAshifVnr/Gig_Space/Payment_Svc/utils/client"
	"github.com/spf13/viper"
)

func InitializeService() *service.PaymentService {
	db, rdb := db.InitializeDB(viper.GetString("DBUrl"))
	repo := repo.NewPaymentRepository(db, rdb)
	user := client.NewUserClinet()
	gig := client.NewGigClinet()
	kafka := config.InitKafkaWriters([]string{viper.GetString("Broker")})
	service := service.NewPaymentService(repo, user, gig, kafka)
	return service
}
