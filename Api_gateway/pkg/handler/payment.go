package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/MuhammedAshifVnr/Gig_Space/api_gateway/utils/helper"
	"github.com/MuhammedAshifVnr/Gig_Space_Proto/proto"
	"github.com/gofiber/fiber/v2"
)

type PaymentHandler struct {
	PaymentClient proto.PaymentServiceClient
}

func NewPaymentHnadller(PayClient proto.PaymentServiceClient) *PaymentHandler {
	return &PaymentHandler{
		PaymentClient: PayClient,
	}
}

// @Summary Create Subscription Payment
// @Description Create a subscription payment for a user based on the selected plan.
// @Tags Payments
// @Accept  json
// @Produce  json
// @Param PlanID path int true "Plan ID"
// @Success 200 {string} string "HTML rendered with subscription details"
// @Router /payments/subscription/{PlanID} [post]
func (h *PaymentHandler) CreateSubscriptionPayment(c *fiber.Ctx) error {
	userID, _ := c.Locals("userID").(uint)
	fmt.Println(userID)
	planID, _ := strconv.Atoi(c.Params("PlanID"))

	res, err := h.PaymentClient.CreateSubscription(context.Background(), &proto.CreateSubscriptionRequest{
		UserId: uint32(userID),
		PlanId: uint32(planID),
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(200).JSON(res.Message)

}

func (h *PaymentHandler) RenderPaymentPage(c *fiber.Ctx) error {
	// Get the subscription_id from the query parameters
	subscriptionID := c.Query("subscription_id")

	// Render the payment page with the subscription ID passed to the frontend
	return c.Render("temp/tests.html", fiber.Map{
		"subscription_id": subscriptionID,
		"message":         "Please complete your payment for the subscription.",
	})
}

func (h *PaymentHandler) UpdatePaymentStatus(c *fiber.Ctx) error {
	var payload struct {
		RazorpayPaymentID string `json:"razorpay_payment_id"`
		RazorpayOrderID   string `json:"razorpay_order_id"`
		RazorpaySignature string `json:"razorpay_signature"`
	}

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	res, err := h.PaymentClient.UpdatePaymentStatus(context.Background(), &proto.UpdatePaymentReq{
		PaymentId: payload.RazorpayPaymentID,
		OrderId:   payload.RazorpayOrderID,
		Signature: payload.RazorpaySignature,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(res)
}

// @Summary      Renew user subscription
// @Description  This endpoint renews the subscription for a given user and plan.
// @Tags         Subscription
// @Accept       json
// @Produce      json
// @Param        PlanID   path      int    true  "Plan ID"
// @Router       /payments/subscription-renew/{PlanID} [post]
func (h *PaymentHandler) RenewSubscription(c *fiber.Ctx) error {
	userID, _ := c.Locals("userID").(uint)
	fmt.Println(userID)
	planID, _ := strconv.Atoi(c.Params("PlanID"))

	res, err := h.PaymentClient.RenewSubscription(context.Background(), &proto.RenewSubReq{
		UserId: uint32(userID),
		PlanId: uint32(planID),
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(200).JSON(res)
}

// @Summary Create a new wallet for the user
// @Description This endpoint allows users to create a wallet with a secure PIN.
// @Tags Wallet
// @Accept  application/x-www-form-urlencoded
// @Produce application/json
// @Param Pin formData int true "4-digit secure PIN"
// @Param RePin formData int true "Repeat PIN to confirm"
// @Router /payments/wallet/create [post]
func (h *PaymentHandler) CreateWallet(c *fiber.Ctx) error {
	pin1 := c.FormValue("Pin")
	pin2 := c.FormValue("RePin")

	if len(pin1) != 4 || pin1 != pin2 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"Message": "PIN must be 4 digits long and both entries must match.",
		})
	}

	PinHash, err := helper.HashPassword(pin1)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"Message": "Error while hashing the PIN.",
			"Error":   err.Error(),
		})
	}

	userID, _ := c.Locals("userID").(uint)

	res, err := h.PaymentClient.CreateWallet(context.Background(), &proto.CreateWalletReq{
		UserId: uint64(userID),
		Pin:    PinHash,
	})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"Message": "Failed to create wallet.",
			"Error":   err.Error(),
		})
	}

	return c.Status(int(res.Status)).JSON(res)
}

// @Summary Get Wallet Information
// @Description Retrieve the wallet balance by providing the correct PIN for security
// @Tags Wallet
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param Pin formData string true "User's wallet PIN"
// @Router /payments/wallet [post]
func (h *PaymentHandler) GetWallet(c *fiber.Ctx) error {
	pin := c.FormValue("Pin")
	userID, _ := c.Locals("userID").(uint)
	fmt.Println("user :", userID, "passw:", pin)
	res, err := h.PaymentClient.GetWallet(context.Background(), &proto.GetwalletReq{
		UserId: uint32(userID),
		Pin:    pin,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"Balance": res.Balance,
		"Status":  200,
	})
}

// @Summary Add Bank Account
// @Description Adds a bank account for the user, ensuring account numbers match and valid IFSC is provided
// @Tags Payment
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param Account1 formData int true "Account Number"
// @Param Account2 formData int true "Confirm Account Number"
// @Param IFSC formData string true "IFSC Code"
// @Param Name formData string true "Beneficiary Name"
// @Router /payments/bank [post]
func (h *PaymentHandler) AddBankAccount(c *fiber.Ctx) error {
	userID, _ := c.Locals("userID").(uint)
	Account1 := c.FormValue("Account1")
	Account2 := c.FormValue("Account2")
	if Account1 != Account2 {
		return c.Status(404).JSON(fiber.Map{
			"error": "Account numbers do not match",
		})
	}
	Ifsc := c.FormValue("IFSC")
	name := c.FormValue("Name")
	fmt.Println("name : ", name, "Account: ", Account1)
	res, err := h.PaymentClient.CreateBankAccount(context.Background(), &proto.CreaBankReq{
		UserId:        uint32(userID),
		Ifsc:          Ifsc,
		AccountNumber: Account1,
		Name:          name,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(int(res.Status)).JSON(res)
}

// @Summary Withdraw amount from user's wallet
// @Description This endpoint allows the user to withdraw a specific amount from their wallet by providing a valid PIN.
// @Tags Payment
// @Accept application/x-www-form-urlencoded
// @Produce application/json
// @Param Pin formData string true "PIN for wallet withdrawal"
// @Param Amount formData number true "Amount to withdraw"
// @Router /payments/withdrawal [post]
func (h *PaymentHandler) Withdrawal(c *fiber.Ctx) error {
	pin := c.FormValue("Pin")
	if pin == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Pin is required",
		})
	}

	amountStr := c.FormValue("Amount")
	amount, err := strconv.ParseFloat(amountStr, 32)
	if err != nil || amount <= 0 {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid amount",
		})
	}

	userID, _ := c.Locals("userID").(uint)
	res, err := h.PaymentClient.Withdrawal(context.Background(), &proto.WithdrawalReq{
		UserId: uint32(userID),
		Pin:    pin,
		Amount: float32(amount),
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(int(res.Status)).JSON(res)
}
