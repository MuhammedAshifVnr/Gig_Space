package repo

import "github.com/MuhammedAshifVnr/Gig_Space/Payment_Svc/pkg/model"

func (r *PaymentRepo)CreateOrderPayment(data model.OrderPayment)error{
	err:=r.DB.Create(&data).Error
	return err
}