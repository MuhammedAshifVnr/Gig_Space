package repo

import (
	"fmt"
	"strings"

	"github.com/MuhammedAshifVnr/Gig_Space/Payment_Svc/pkg/model"
)

func (r *PaymentRepo) CreateWallet(data model.Wallet) error {
	if err := r.DB.Create(&data).Error; err != nil {

		if strings.Contains(err.Error(), "duplicate key value") {
			return fmt.Errorf("wallet already exists for this user")
		}
		return err
	}
	return nil
}

func (r *PaymentRepo) GetWallet(ID uint) (model.Wallet, error) {
	query := `SELECT * FROM wallets WHERE user_id = ?`
	var wallet model.Wallet
	err := r.DB.Raw(query, ID).Scan(&wallet)
	return wallet, err.Error
}

func (r *PaymentRepo) AddFundAccID(FundID string, userID uint) error {
	query := "UPDATE wallets SET fund_account_id = ? WHERE user_id = ?"
	err := r.DB.Exec(query, FundID, userID)
	return err.Error
}
