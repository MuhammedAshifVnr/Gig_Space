package helper

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/razorpay/razorpay-go"
	"github.com/spf13/viper"
)

const (
	razorpayBaseURL = "https://api.razorpay.com/v1"
)

func RandString() string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	bytes := make([]byte, 8)
	_, err := rand.Read(bytes)
	if err != nil {
		return ""
	}

	for i := 0; i < 8; i++ {
		bytes[i] = letters[bytes[i]%byte(len(letters))]
	}
	return string(bytes)
}

func AddFundAccount(contactID string, bankAccount map[string]interface{}) (string, error) {
	client := razorpay.NewClient(viper.GetString("XApiKey"), viper.GetString("XApiSecret"))
	fundAccountData := map[string]interface{}{
		"contact_id":   contactID,
		"account_type": "bank_account",
		"bank_account": bankAccount,
	}

	body, err := client.FundAccount.Create(fundAccountData, nil)
	
	if err != nil {
		return "", err
	}

	return body["id"].(string), nil
}

func CreateContact(userID, name, email, phone string) (string, error) {
	contactData := map[string]interface{}{
		"name":         name,
		"email":        email,
		"contact":      phone,
		"type":         "customer",
		"reference_id": userID,
	}

	jsonData, err := json.Marshal(contactData)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/contacts", razorpayBaseURL), bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(viper.GetString("XApiKey"), viper.GetString("XApiSecret"))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return "", errors.New("failed to create contact")
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}

	contactID, ok := result["id"].(string)
	if !ok {
		return "", errors.New("invalid response from Razorpay")
	}

	return contactID, nil
}

func Payout(payoutData map[string]interface{}) (string, error) {
	jsonData, err := json.Marshal(payoutData)
	if err != nil {
		fmt.Println("Error marshaling payout data:", err)
		return "", err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/payouts", razorpayBaseURL), bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating new request:", err)
		return "", err
	}

	// Set Razorpay API credentials
	req.SetBasicAuth(viper.GetString("XApiKey"), viper.GetString("XApiSecret"))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return "", err
	}

	// Print response body for debugging
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		fmt.Printf("Payout failed with status code %d, response: %s\n", resp.StatusCode, string(body))
		return "", errors.New("failed to create payout")
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}

	payoutID, ok := result["id"].(string)
	if !ok {
		return "", errors.New("invalid response from Razorpay")
	}

	return payoutID, nil
}
