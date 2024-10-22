package repo

import (
	"fmt"
	"log"
	"time"

	"github.com/MuhammedAshifVnr/Gig_Space/Payment_Svc/pkg/model"
	"gorm.io/gorm"
)

type PaymentRepo struct {
	DB *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) *PaymentRepo {
	return &PaymentRepo{
		DB: db,
	}
}

func (r *PaymentRepo) GetActiveSubscription(userID uint) (model.Subscription, error) {
	var userSub model.Subscription
	query := "SELECT * FROM subscriptions WHERE user_id = ?"
	err := r.DB.Raw(query, userID).Scan(&userSub)
	if err.Error != nil {
		log.Println("Failed to find user : ", err.Error)
		return userSub, err.Error
	}
	return userSub, nil
}

func (r *PaymentRepo) GetPlanByID(id uint) (model.Plan, error) {
	var plan model.Plan
	query := `SELECT * FROM plans WHERE id =?`
	err := r.DB.Raw(query, id).Scan(&plan)
	if err.Error != nil {
		log.Println("Faild to find plan : ", err)
		return model.Plan{}, err.Error
	}
	return plan, nil
}

func (r *PaymentRepo) CreateSubscription(sub model.Subscription) error {
	err := r.DB.Create(&sub).Error
	if err != nil {
		log.Println("Faild to create subscription: ", err)
		return err
	}
	return nil
}

func (r *PaymentRepo) CreatePayment(payment model.Payment) error {
	err := r.DB.Create(&payment).Error
	if err != nil {
		log.Println("Faild to create payment: ", err)
		return err
	}
	return nil
}

func (r *PaymentRepo) UpdateStatus(orderID, transaction_id, status string) error {
	query := `UPDATE order_payments  SET status = ?, transaction_id =? WHERE order_id = ?`
	err := r.DB.Exec(query, status, transaction_id, orderID).Error
	if err != nil {
		fmt.Println("--")
		return err
	}
	// startDate := time.Now()
	// query = `UPDATE subscriptions SET active = ? ,start_date =?,end_date = ? WHERE subscription_id = ?`
	// err = r.DB.Exec(query, "Active", startDate, startDate, startDate.AddDate(0, 0, 30), 1).Error
	 return err
}

// func (r *PaymentRepo) UpdateSubscription(sub model.Subscription) error {
// 	err := r.DB.Save(&sub).Error
// 	if err != nil {
// 		log.Println("Faild to renew subscripton: ", err)
// 	}
// 	return err
// }

func (r *PaymentRepo) GetSubDetails(userID uint) (time.Time, error) {
	var endDate time.Time
	query := `select end_date from subscriptions where user_id = ? `
	err := r.DB.Raw(query, userID).Scan(&endDate).Error
	if err != nil {
		log.Println("Faild to Find Subscripton Details: ", err)
	}
	return endDate, err
}

func (r *PaymentRepo) UpdatePayment(payment model.Payment) error {
	return r.DB.Model(&model.Payment{}).Where("transaction_id = ?", payment.TransactionID).Updates(payment).Error
}

// UpdateSubscription updates the subscription record in the database
func (r *PaymentRepo) UpdateSubscription(subscription model.Subscription) error {
	return r.DB.Model(&model.Subscription{}).Where("subscription_id = ?", subscription.SubscriptionID).Updates(subscription).Error
}

// GetSubscriptionByID retrieves a subscription by its ID
func (r *PaymentRepo) GetSubscriptionByID(subscriptionID string) (*model.Subscription, error) {
	var subscription model.Subscription
	err := r.DB.Where("subscription_id = ?", subscriptionID).First(&subscription).Error
	if err != nil {
		return nil, err
	}
	return &subscription, nil
}
