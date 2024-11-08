package service

import (
	"context"
	"fmt"

	"github.com/MuhammedAshifVnr/Gig_Space/Payment_Svc/pkg/logger"
	"github.com/MuhammedAshifVnr/Gig_Space/Payment_Svc/pkg/model"
	"github.com/MuhammedAshifVnr/Gig_Space_Proto/proto"
	"github.com/sirupsen/logrus"
)

func (s *PaymentService) CreatePlan(ctx context.Context, req *proto.CreatePlanReq) (*proto.PaymentCommonRes, error) {

	planData := map[string]interface{}{
		"period":   req.Period,
		"interval": req.Interval,
		"item": map[string]interface{}{
			"name":     req.Name,
			"amount":   req.Amount * 100,
			"currency": "INR",
		},
	}

	plan, err := s.RazorClient.Plan.Create(planData, nil)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"plan_name": req.Name,
			"amount":    req.Amount,
		}).Error("Failed to create plan via Razorpay SDK: ", err)
		return nil, fmt.Errorf("failed to create plan via Razorpay SDK: %w", err)
	}

	err = s.Repo.CreatePlan(model.Plan{
		Name:           req.Name,
		Price:          int(req.Amount),
		RazorpayPlanID: plan["id"].(string),
		Period:         req.Period,
		Interval:       int(req.Interval),
	})
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"plan_name":        req.Name,
			"razorpay_plan_id": plan["id"].(string),
		}).Error("Failed to save plan in the database: ", err)
		return nil, err
	}

	logger.Log.WithFields(logrus.Fields{
		"plan_name":        req.Name,
		"razorpay_plan_id": plan["id"].(string),
	}).Info("Plan created and saved successfully")

	return &proto.PaymentCommonRes{
		Message: "Plan added successfully",
		Status:  200,
	}, nil

}

func (s *PaymentService) GetPlan(ctx context.Context, req *proto.EmptyReq) (*proto.GetPlanRes, error) {
	plans, err := s.Repo.GetAllPlans()
	if err != nil {
		logger.Log.Error("Failed to retrieve plans from database: ", err)
		return nil, err
	}
	return &proto.GetPlanRes{
		Plans: plans,
	}, nil
}

func (s *PaymentService) DeletePlan(ctx context.Context, req *proto.DeletePlanReq) (*proto.PaymentCommonRes, error) {
	err := s.Repo.DeletePlan(req.PlanId)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"plan_id": req.PlanId,
		}).Error("Failed to delete plan: ", err)
		return nil, err
	}

	logger.Log.WithFields(logrus.Fields{
		"plan_id": req.PlanId,
	}).Info("Plan deleted successfully")

	return &proto.PaymentCommonRes{
		Message: "Plan deleted successfully",
		Status:  200,
	}, nil
}
