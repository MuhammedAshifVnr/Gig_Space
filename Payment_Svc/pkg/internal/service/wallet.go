package service

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/MuhammedAshifVnr/Gig_Space/Payment_Svc/pkg/model"
	"github.com/MuhammedAshifVnr/Gig_Space/Payment_Svc/utils/helper"
	"github.com/MuhammedAshifVnr/Gig_Space_Proto/proto"
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
		log.Println("Faild to create wallet :", err.Error())
		return nil, err
	}
	return &proto.PaymentCommonRes{
		Message: "Wallet created Successfully",
		Status:  200,
	}, nil
}

func (s *PaymentService) GetWallet(ctx context.Context, req *proto.GetwalletReq) (*proto.WalletRes, error) {
	wallet, err := s.Repo.GetWallet(uint(req.UserId))
	if err != nil {
		log.Println("Faild to find wallet : ", err.Error())
		return nil, err
	}

	fmt.Println("pin", wallet.Pin_hash)
	err = bcrypt.CompareHashAndPassword([]byte(wallet.Pin_hash), []byte(req.Pin))
	if err != nil {
		log.Println("Invalid Password: ", err.Error())
		return nil, errors.New("invalid password")
	}

	return &proto.WalletRes{
		Balance: float32(wallet.Balance),
	}, nil

}

func (s *PaymentService) CreateBankAccount(ctx context.Context, req *proto.CreaBankReq) (*proto.PaymentCommonRes, error) {
	user, err := s.UserClient.GetUserProfile(context.Background(), &proto.ProfileReq{UserId: req.UserId})
	if err != nil {
		log.Println("Faild to Find User: ", err.Error())
		return nil, err
	}
	contactID, err := helper.CreateContact(string(req.UserId), user.Firstname, user.Email, user.Phone)
	if err != nil {
		log.Println("Faild to Create Contact: ", err.Error())
		return nil, err
	}

	bankDetails := map[string]interface{}{
		"name":           req.Name,
		"account_number": req.AccountNumber,
		"ifsc":           req.Ifsc,
	}
	FundID, err := helper.AddFundAccount(contactID, bankDetails)
	if err != nil {
		log.Println("Faild to Add Fund Accoutn: ", err.Error())
		return nil, err
	}
	err = s.Repo.AddFundAccID(FundID, uint(req.UserId))
	if err != nil {
		log.Println("Faild to Save Fund Accoutn ID: ", err.Error())
		return nil, err
	}
	return &proto.PaymentCommonRes{Message: FundID, Status: 200}, nil
}

func (s *PaymentService) Withdrawal(ctx context.Context, req *proto.WithdrawalReq) (*proto.PaymentCommonRes, error) {
	wallet, err := s.Repo.GetWallet(uint(req.UserId))
	if err != nil {
		log.Println("Faild to get wallet: ", err.Error())
		return nil, err
	}
	fmt.Println(wallet)
	err = bcrypt.CompareHashAndPassword([]byte(wallet.Pin_hash), []byte(req.Pin))
	if err != nil {
		log.Println("Invalid Password: ", err.Error())
		return nil, errors.New("invalid password")
	}
	if wallet.Balance < int64(req.Amount) {
		log.Println("Insufficient balance")
		return nil, errors.New("insufficient balance")
	}

	data := map[string]interface{}{
		"account_number":       "2323230063217806",
		"fund_account_id":      wallet.Fund_account_id,
		"amount":               req.Amount*100,
		"currency":             "INR",
		"mode":                 "IMPS",
		"purpose":              "payout",
		"queue_if_low_balance": true,
		"reference_id":         "Payout123", //.........
		"notes": map[string]interface{}{
			"custom_note": "This is a custom note",//..............
		},
	}
	res, err := helper.Payout(data)
	if err != nil {
		log.Println("Faild to payout: ", err.Error())
		return nil, err
	}
	return &proto.PaymentCommonRes{
		Message: res,
		Status:  200,
	}, nil
}
