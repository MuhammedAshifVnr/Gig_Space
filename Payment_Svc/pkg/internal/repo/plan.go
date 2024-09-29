package repo

import (
	"github.com/MuhammedAshifVnr/Gig_Space/Payment_Svc/pkg/model"
	"github.com/MuhammedAshifVnr/Gig_Space_Proto/proto"
)

func (r *PaymentRepo) CreatePlan(planData model.Plan) error {
	return r.DB.Create(&planData).Error
}

func (r *PaymentRepo) GetAllPlans() ([]*proto.Plan, error) {
	var plans []*proto.Plan
	query := `SELECT name, price,  razorpay_plan_id, period from plans`
	err := r.DB.Raw(query).Scan(&plans).Error
	return plans, err
}

func (r *PaymentRepo) DeletePlan(planID string) error {
	query := `DELETE FROM plans WHERE razorpay_plan_id = ?`
	return r.DB.Exec(query, planID).Error
}
