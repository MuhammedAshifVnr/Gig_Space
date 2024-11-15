package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/MuhammedAshifVnr/Gig_Space/Payment_Svc/pkg/model"
	"github.com/MuhammedAshifVnr/Gig_Space/Payment_Svc/utils/helper"
	"github.com/MuhammedAshifVnr/Gig_Space/Payment_Svc/utils/notification"
	"github.com/MuhammedAshifVnr/Gig_Space_Proto/proto"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

func (s *PaymentService) CreateWallet(ctx context.Context, req *proto.CreateWalletReq) (*proto.PaymentCommonRes, error) {
	wallet := model.Wallet{
		UserID:   uint(req.UserId),
		Balance:  0,
		Pin_hash: req.Pin,
	}
	err := s.Repo.CreateWallet(wallet)

	if err != nil {
		s.Log.WithFields(logrus.Fields{
			"user_id": req.UserId,
			"error":   err.Error(),
		}).Error("Failed to create wallet")
		return nil, err
	}

	s.Log.WithFields(logrus.Fields{
		"user_id": req.UserId,
	}).Info("Wallet created successfully")

	return &proto.PaymentCommonRes{
		Message: "Wallet created Successfully",
		Status:  200,
	}, nil
}

func (s *PaymentService) GetWallet(ctx context.Context, req *proto.GetwalletReq) (*proto.WalletRes, error) {
	wallet, err := s.Repo.GetWallet(uint(req.UserId))
	if err != nil {
		s.Log.WithFields(logrus.Fields{
			"user_id": req.UserId,
			"error":   err.Error(),
		}).Error("Failed to find wallet")
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(wallet.Pin_hash), []byte(req.Pin))
	if err != nil {
		s.Log.WithFields(logrus.Fields{
			"user_id": req.UserId,
			"error":   err.Error(),
		}).Warn("Invalid password attempt")
		return nil, errors.New("invalid password")
	}

	s.Log.WithFields(logrus.Fields{
		"user_id": req.UserId,
	}).Info("Wallet retrieved successfully")

	return &proto.WalletRes{
		Balance: float32(wallet.Balance),
	}, nil

}

func (s *PaymentService) CreateBankAccount(ctx context.Context, req *proto.CreaBankReq) (*proto.PaymentCommonRes, error) {
	user, err := s.UserClient.GetUserProfile(context.Background(), &proto.ProfileReq{UserId: req.UserId})
	if err != nil {
		s.Log.WithFields(logrus.Fields{
			"user_id": req.UserId,
		}).Error("Failed to find user profile:", err)
		return nil, err
	}
	contactID, err := helper.CreateContact(string(req.UserId), user.Firstname, user.Email, user.Phone)
	if err != nil {
		s.Log.WithField("user_id", req.UserId).Error("Failed to create contact:", err)
		return nil, err
	}

	bankDetails := map[string]interface{}{
		"name":           req.Name,
		"account_number": req.AccountNumber,
		"ifsc":           req.Ifsc,
	}
	FundID, err := helper.AddFundAccount(contactID, bankDetails)
	if err != nil {
		s.Log.WithField("user_id", req.UserId).Error("Failed to add fund account:", err)
		return nil, err
	}
	err = s.Repo.AddFundAccID(FundID, uint(req.UserId))
	if err != nil {
		s.Log.WithField("user_id", req.UserId).Error("Failed to save fund account ID:", err)
		return nil, err
	}
	s.Log.WithField("user_id", req.UserId).Info("Bank account created successfully")
	return &proto.PaymentCommonRes{Message: FundID, Status: 200}, nil
}

func (s *PaymentService) Withdrawal(ctx context.Context, req *proto.WithdrawalReq) (*proto.PaymentCommonRes, error) {
	wallet, err := s.Repo.GetWallet(uint(req.UserId))
	if err != nil {
		s.Log.WithField("user_id", req.UserId).Error("Failed to get wallet:", err)
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(wallet.Pin_hash), []byte(req.Pin))
	if err != nil {
		s.Log.WithField("user_id", req.UserId).Warn("Invalid PIN provided")
		return nil, errors.New("invalid password")
	}
	if wallet.Balance < int64(req.Amount) {
		s.Log.WithField("user_id", req.UserId).Warn("Insufficient balance")
		return nil, errors.New("insufficient balance")
	}

	data := map[string]interface{}{
		"account_number":       "2323230063217806", //env
		"fund_account_id":      wallet.Fund_account_id,
		"amount":               req.Amount * 100,
		"currency":             "INR",
		"mode":                 "IMPS",
		"purpose":              "payout",
		"queue_if_low_balance": true,
		"reference_id":         "Payout123", //.........
		"notes": map[string]interface{}{
			"custom_note": "This is a custom note", //..............
		},
	}
	res, err := helper.Payout(data)
	if err != nil {
		s.Log.WithField("user_id", req.UserId).Error("Failed to process payout:", err)
		return nil, err
	}
	wallet.Balance -= int64(req.Amount)
	err = s.Repo.UpdateWallet(wallet)
	if err != nil {
		s.Log.WithField("user_id", req.UserId).Error("Failed to update wallet balance:", err)
		return nil, err
	}
	s.Log.WithField("user_id", req.UserId).Info("Withdrawal successful")
	return &proto.PaymentCommonRes{
		Message: res,
		Status:  200,
	}, nil
}

func (s *PaymentService) AddWalletRefund(ctx context.Context, req *proto.AddRefund) (*proto.PaymentCommonRes, error) {
	err := s.Repo.AddRefundAmount(uint(req.UserId), int(req.Amount))
	if err != nil {
		s.Log.WithFields(logrus.Fields{
			"user_id": req.UserId,
			"amount":  req.Amount,
		}).Error("Failed to update balance:", err)
		return nil, err
	}
	s.Log.WithField("user_id", req.UserId).Info("Refund added to wallet successfully")
	return nil, nil
}

func (s *PaymentService) ChangeWalletPin(ctx context.Context, req *proto.ChangePinReq) (*proto.PaymentCommonRes, error) {
	wallet, err := s.Repo.GetWallet(uint(req.UserId))
	if err != nil {
		s.Log.WithField("user_id", req.UserId).Error("Failed to find wallet:", err)
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(wallet.Pin_hash), []byte(req.CurrentPin))
	if err != nil {
		s.Log.WithField("user_id", req.UserId).Warn("Invalid current PIN provided")
		return nil, errors.New("invalid password")
	}

	err = s.Repo.UpdatePin(uint(req.UserId), req.NewPin)
	if err != nil {
		s.Log.WithField("user_id", req.UserId).Error("Failed to update PIN:", err)
		return nil, err
	}

	s.Log.WithField("user_id", req.UserId).Info("Wallet PIN changed successfully")
	return &proto.PaymentCommonRes{
		Message: "Pin Changed Successfully",
		Status:  200,
	}, nil
}

func (s *PaymentService) ForgotWalletPin(ctx context.Context, req *proto.ForgotPinReq) (*proto.PaymentCommonRes, error) {
	res, err := s.UserClient.GetUserEmail(ctx, &proto.ProfileReq{UserId: uint32(req.UserId)})
	if err != nil {
		s.Log.WithFields(logrus.Fields{
			"user_id": req.UserId,
		}).Errorf("failed to retrieve user email: %v", err)
		return nil, fmt.Errorf("failed to retrieve user email: %w", err)
	}
	otp := helper.GenerateOtp()
	err = s.Repo.ForgotPinOtp(model.Wallet{UserID: uint(req.UserId)}, otp)
	if err != nil {
		s.Log.WithFields(logrus.Fields{
			"user_id": req.UserId,
			"otp":     otp,
		}).Errorf("failed to save OTP to repository: %v", err)
		return nil, fmt.Errorf("failed to save OTP: %w", err)
	}
	forgotTopic := viper.GetString("ForgotTopic")
	writer, ok := s.kafkaWriter[forgotTopic]
	if ok {
		err = notification.SendNotification(ctx, writer, model.ForgotEvent{
			Otp:   otp,
			Email: res.Email,
			Event: "Wallet",
		}, forgotTopic)
		if err != nil {
			s.Log.WithFields(logrus.Fields{
				"user_id": req.UserId,
				"email":   res.Email,
			}).Errorf("failed to publish notification: %v", err)
			return nil, fmt.Errorf("failed to publish notification: %w", err)
		}
	} else {
		s.Log.WithFields(logrus.Fields{
			"topic": "wallet forgot",
		}).Warn("Kafka writer not found for forgot pin topic")
	}

	s.Log.WithFields(logrus.Fields{
		"user_id": req.UserId,
		"email":   res.Email,
	}).Info("OTP sent successfully to user's email")

	return &proto.PaymentCommonRes{
		Message: "OTP sended into yout mail.",
		Status:  200,
	}, nil
}

func (s *PaymentService) ResetWalletPin(ctx context.Context, req *proto.PinResetReq) (*proto.PaymentCommonRes, error) {
	val, err := s.Repo.VerifyOtp(req.Otp)
	if err == redis.Nil {
		s.Log.WithField("otp", req.Otp).Warn("Invalid or expired OTP")
		return nil, errors.New("this OTP was invalid or expired")
	}

	var user model.Wallet
	err = json.Unmarshal([]byte(val), &user)
	if err != nil {
		s.Log.WithField("otp", req.Otp).Error("Failed to unmarshal user data:", err)
		return nil, errors.New("Could not unmarshal user: " + err.Error())
	}

	err = s.Repo.UpdatePin(user.UserID, req.Pin)
	if err != nil {
		s.Log.WithField("user_id", user.UserID).Error("Failed to update wallet PIN:", err)
		return nil, err
	}

	s.Log.WithField("user_id", user.UserID).Info("Wallet PIN reset successfully")
	return &proto.PaymentCommonRes{
		Message: "Wallet Pin Updated",
		Status:  200,
	}, nil
}
