package helper

import (
	"context"

	"github.com/MuhammedAshifVnr/Gig_Space/Notificaton_svc/utils/client"
	"github.com/MuhammedAshifVnr/Gig_Space_Proto/proto"
)

func GetUserID(orderID string)(uint64,error){
	res,err:=client.NewGigClinet().GetFreelancerIDByOrder(context.Background(),&proto.OrderByIDReq{OrderId: orderID})
	return res.UserId,err
}