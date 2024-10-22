package repo

import (
	"time"

	"github.com/MuhammedAshifVnr/Gig_Space/Payment_Svc/pkg/model"
	"github.com/MuhammedAshifVnr/Gig_Space_Proto/proto"
)

type RepoInter interface {
	GetActiveSubscription(userID uint) (model.Subscription, error)
	GetPlanByID(id uint) (model.Plan, error)
	CreateSubscription(sub model.Subscription) error
	CreatePayment(payment model.Payment) error
	UpdateStatus(orderID, transaction_id, status string) error
	//UpdateSubscription(sub model.Subscription) error
	GetSubDetails(userID uint) (time.Time, error)
	CreatePlan(planData model.Plan) error
	GetAllPlans() ([]*proto.Plan, error)
	DeletePlan(planID string) error
	UpdatePayment(payment model.Payment) error
    UpdateSubscription(subscription model.Subscription) error
    GetSubscriptionByID(subscriptionID string) (*model.Subscription, error)
	CreateOrderPayment(data model.OrderPayment)error
	CreateWallet(data model.Wallet) error
	GetWallet(ID uint) (model.Wallet, error)
	AddFundAccID(FundID string, userID uint) error
	UpdateWallet(wallet model.Wallet)error
}
