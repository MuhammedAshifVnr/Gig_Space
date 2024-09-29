package handler

import (
	"context"
	"fmt"
	"strconv"

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
