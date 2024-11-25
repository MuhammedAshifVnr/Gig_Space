package handler

import (
	"context"
	"fmt"
	"log"
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
	pin, err := helper.GetRequiredFormValue(c, "Pin")
	if err != nil {
		return err
	}
	userID, _ := c.Locals("userID").(uint)
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
	Account1, err := helper.GetRequiredFormValue(c, "Account1")
	if err != nil {
		return err
	}
	Account2 := c.FormValue("Account2")
	if Account1 != Account2 {
		return c.Status(404).JSON(fiber.Map{
			"error": "Account numbers do not match",
		})
	}
	Ifsc, err := helper.GetRequiredFormValue(c, "IFSC")
	if err != nil {
		return err
	}
	name, err := helper.GetRequiredFormValue(c, "Name")
	if err != nil {
		return err
	}
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
	pin, err := helper.GetRequiredFormValue(c, "Pin")
	if err != nil {
		return err
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

// @Summary Change Wallet PIN
// @Description Allows the user to change their wallet PIN by providing the current PIN and a new PIN.
// @Tags Wallet
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param current_pin formData int true "Current PIN"
// @Param new_pin1 formData int true "New PIN"
// @Param new_pin2 formData int true "Confirm New PIN"
// @Router /payments/wallet/change-pin [post]
func (h *PaymentHandler) ChangeWalletPin(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(uint)
	if !ok {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	currentPin, err := strconv.Atoi(c.FormValue("current_pin"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid format for current PIN"})
	}

	newPin1, err := strconv.Atoi(c.FormValue("new_pin1"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid format for new PIN"})
	}

	newPin2, err := strconv.Atoi(c.FormValue("new_pin2"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid format for confirm new PIN"})
	}

	if newPin1 != newPin2 {
		return c.Status(400).JSON(fiber.Map{"error": "New PIN entries do not match"})
	}

	hashPin, err := helper.HashPassword(strconv.Itoa(newPin2))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error hashing PIN"})
	}

	res, err := h.PaymentClient.ChangeWalletPin(context.Background(), &proto.ChangePinReq{
		UserId:     uint64(userID),
		CurrentPin: strconv.Itoa(currentPin),
		NewPin:     hashPin,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(int(res.Status)).JSON(res)
}

// @Summary      Request OTP for resetting wallet PIN
// @Description  Generates and sends an OTP to the user's registered email to reset their wallet PIN.
// @Tags         Wallet
// @Accept       json
// @Produce      json
// @Router       /payments/wallet/forgot-pin [post]
func (h *PaymentHandler) ForgotWalletPin(c *fiber.Ctx) error {
	userID, _ := c.Locals("userID").(uint)
	res, err := h.PaymentClient.ForgotWalletPin(context.Background(), &proto.ForgotPinReq{
		UserId: uint64(userID),
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(int(res.Status)).JSON(res)
}

// @Summary Reset Wallet PIN
// @Description Resets the wallet PIN for the user after OTP verification
// @Tags Wallet
// @Accept multipart/form-data
// @Produce json
// @Param OTP formData string true "One-Time Password (OTP)"
// @Param new_pin1 formData int true "New PIN"
// @Param new_pin2 formData int true "Confirm New PIN"
// @Router /payments/wallet/reset-pin [post]
func (h *PaymentHandler) ResetWalletPin(c *fiber.Ctx) error {
	otp, err := helper.GetRequiredFormValue(c, "OTP")
	if err != nil {
		return err
	}
	newPin1, err := strconv.Atoi(c.FormValue("new_pin1"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid format for new PIN"})
	}

	newPin2, err := strconv.Atoi(c.FormValue("new_pin2"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid format for confirm new PIN"})
	}

	if newPin1 != newPin2 {
		return c.Status(400).JSON(fiber.Map{"error": "New PIN entries do not match"})
	}

	hashPin, err := helper.HashPassword(strconv.Itoa(newPin2))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error hashing PIN"})
	}
	res, err := h.PaymentClient.ResetWalletPin(context.Background(), &proto.PinResetReq{
		Otp: otp,
		Pin: hashPin,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(int(res.Status)).JSON(res)
}

func (h *PaymentHandler) UpdateWebhook(c *fiber.Ctx) error {
    // Log the raw payload for debugging
    log.Printf("Webhook payload: %s", string(c.Body()))

    // Parse the payload into a generic map
    var payload map[string]interface{}
    if err := c.BodyParser(&payload); err != nil {
        log.Println("Failed to parse webhook payload:", err)
        return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
    }

    log.Println("Parsed payload:", payload)

    // Safely extract fields
    event, ok := payload["event"].(string)
    if !ok {
        return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "missing or invalid 'event'"})
    }

    // Access nested structures
    subscriptionData, subExists := payload["payload"].(map[string]interface{})["subscription"].(map[string]interface{})
    if !subExists {
        return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "missing 'subscription' data in payload"})
    }

    subscriptionEntity, entityExists := subscriptionData["entity"].(map[string]interface{})
    if !entityExists {
        return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "missing 'entity' in subscription data"})
    }

    subscriptionID, idOk := subscriptionEntity["id"].(string)
    if !idOk {
        return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "missing or invalid 'id' in subscription entity"})
    }

    // Access other nested fields like "amount", handling null or missing fields
    var amount string
    if val, ok := subscriptionEntity["amount"].(float64); ok {
        amount = fmt.Sprintf("%.0f", val) // Convert float64 to string
    } else {
        log.Println("Missing 'amount', setting as '0'")
        amount = "0" // Default value
    }

    // Call the gRPC client
    res, err := h.PaymentClient.HandleWebhook(context.Background(), &proto.WebhookRequest{
        Payload: map[string]string{
            "event":         event,
            "entity.id":     subscriptionID,
            "entity.amount": amount,
        },
    })
    if err != nil {
        return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }

    return c.Status(http.StatusOK).JSON(res)
}
