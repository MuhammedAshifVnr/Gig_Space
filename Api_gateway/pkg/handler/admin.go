package handler

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/MuhammedAshifVnr/Gig_Space/api_gateway/utils/helper"
	"github.com/MuhammedAshifVnr/Gig_Space_Proto/proto"
	"github.com/gofiber/fiber/v2"
)

// AdminLogin
// @Summary Admin login
// @Description Log in as an admin using email and password.
// @Tags Admin
// @Accept json
// @Produce json
// @Param request body helper.ADLogin true "Admin login credentials"
// @Router /admin/login [post]
func (h *UserHandler) AdminLogin(c *fiber.Ctx) error {
	var req helper.ADLogin
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse JSON",
		})
	}
	res, err := h.userClient.AdminLogin(context.Background(), &proto.AdLoginReq{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to login",
			"error":   err,
		})
	}
	tokenValue := res.Data["token"].GetStringValue()

	c.Cookie(&fiber.Cookie{
		Name:     "AdminToken",
		Value:    tokenValue,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
	})
	return c.Status(int(res.Status)).JSON(res)
}

// AddCategory
// @Summary Add a new category
// @Description Allows admin to add a new category
// @Tags Admin
// @Accept json
// @Produce json
// @Param category_name body helper.AddCategory true "Category Name"
// @Router /admin/category [post]
func (h *UserHandler) AddCategory(c *fiber.Ctx) error {
	var req helper.AddCategory
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse JSON",
		})
	}
	res, err := h.userClient.AddCategory(context.Background(), &proto.CategoryReq{
		Name: req.Name,
	})
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(int(res.Status)).JSON(res)
}

// AddSkill
// @Summary Add a new skill to a user's profile
// @Description This endpoint allows a admin to add a new skill by providing the skill name.
// @Tags Admin
// @Accept  json
// @Produce  json
// @Param   skill  body  helper.AddSkill  true  "Skill information"
// @Router /admin/skill [post]
func (h *UserHandler) AddSkill(c *fiber.Ctx) error {
	var req helper.AddSkill
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse JSON",
		})
	}
	res, err := h.userClient.AddSkill(context.Background(), &proto.AddSkillReq{
		SkillName: req.Name,
	})
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(int(res.Status)).JSON(res)
}

// @Summary Get all categories
// @Description Retrieves a list of all categories
// @Tags Admin
// @Accept json
// @Produce json
// @Router /admin/category [get]
func (h *UserHandler) GetCategory(c *fiber.Ctx) error {
	res, err := h.userClient.GetCategory(context.Background(), &proto.EmtpyReq{})
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"data":   res,
		"status": 200,
	})
}

// @Summary Get all skills
// @Description Retrieves a list of all skills
// @Tags Admin
// @Accept json
// @Produce json
// @Router /admin/skills [get]
func (h *UserHandler) GetSkills(c *fiber.Ctx) error {
	res, err := h.userClient.GetSkill(context.Background(), &proto.EmtpyReq{})
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"data":   res,
		"status": 200,
	})
}

// @Summary Admin delete skill
// @Description Deletes a skill by ID, used by admin
// @Tags Admin
// @Param skillID path int true "Skill ID"
// @Accept json
// @Produce json
// @Router /admin/skills/{skillID} [delete]
func (h *UserHandler) AdDeleteSkill(c *fiber.Ctx) error {
	id := c.Params("skillID")
	skillID, _ := strconv.Atoi(id)
	_, err := h.userClient.AdDeleteSkill(context.Background(), &proto.ADeleteSkillReq{Id: uint32(skillID)})
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"message": "Skill successfully deleted",
		"status":  200,
	})

}

// @Summary Admin delete skill
// @Description Deletes a category by ID, used by admin
// @Tags Admin
// @Param CatID path int true "Category ID"
// @Accept json
// @Produce json
// @Router /admin/category/{CatID} [delete]
func (h *UserHandler) AdDeleteCategory(c *fiber.Ctx) error {
	id := c.Params("CatID")
	skillID, _ := strconv.Atoi(id)
	_, err := h.userClient.DeleteCategory(context.Background(), &proto.DeleteCatReq{Id: uint32(skillID)})
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"message": "Category successfully deleted",
		"status":  200,
	})

}

// @Summary Get all users
// @Description Retrieve a list of all users
// @Tags Admin
// @Accept json
// @Produce json
// @Router /admin/users [get]
func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	res, err := h.userClient.GetAllUsers(context.Background(), &proto.EmtpyReq{})
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"data":   res,
		"status": 200,
	})

}

// @Summary Block a user by userID
// @Description This endpoint blocks a user by their user ID. The userID is retrieved from the URL path and must be a valid integer.
// @Tags Admin
// @Accept json
// @Produce json
// @Param userID path int true "User ID"
// @Router /admin/block/{userID} [post]
func (h *UserHandler) UserBlock(c *fiber.Ctx) error {
	userID := c.Params("userID")
	ID, _ := strconv.Atoi(userID)
	res, err := h.userClient.AdminUserBlock(context.Background(), &proto.BlockReq{Id: uint32(ID)})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(int(res.Status)).JSON(res)
}

// @Summary Create a new payment plan
// @Description This endpoint allows you to create a new payment plan by providing a name, price, and period.
// @Tags Payment
// @Accept application/x-www-form-urlencoded
// @Produce application/json
// @Param name formData string true "Name of the plan"
// @Param price formData int true "Price of the plan in cents"
// @Param period query string true "Period of the plan" Enums(monthly, yearly)
// @Router /admin/create-plan [post]
func (h *PaymentHandler) CreatePlan(c *fiber.Ctx) error {
	Name := c.FormValue("name")
	Price, _ := strconv.Atoi(c.FormValue("price"))
	Period := c.FormValue("period")

	if Name == "" || Price == 0 || Period == "" {
		return c.Status(400).JSON("eroor")
	}
	res, err := h.PaymentClient.CreatePlan(context.Background(), &proto.CreatePlanReq{
		Name:     Name,
		Amount:   int64(Price),
		Period:   Period,
		Interval: 1,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(int(res.Status)).JSON(res)
}

// @Summary Get all payment plans
// @Description This endpoint retrieves a list of all available payment plans.
// @Tags Payment
// @Produce application/json
// @Router /admin/plans [get]
func (h *PaymentHandler) GetPlans(c *fiber.Ctx) error {
	res, err := h.PaymentClient.GetPlan(context.Background(), &proto.EmptyReq{})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(200).JSON(res)
}

// @Summary Delete a payment plan
// @Description This endpoint deletes a payment plan based on the provided PlanID.
// @Tags Payment
// @Param PlanID path string true "ID of the plan to delete"
// @Produce application/json
// @Router /admin/plan/{PlanID} [delete]
func (h *PaymentHandler) DeletPlan(c *fiber.Ctx) error {
	id := c.Params("PlanID")
	res, err := h.PaymentClient.DeletePlan(context.Background(), &proto.DeletePlanReq{
		PlanId: id,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(200).JSON(res)
}

// @Summary Get all refund orders
// @Description Retrieves a list of all refund orders in the system.
// @Tags GetAllRefund
// @Accept json
// @Produce json
// @Router /admin/orders/refund [get]
func (h *GigHandler) GetAllRefundOrders(c *fiber.Ctx) error {
	res, err := h.GigClinet.AdminOrderController(context.Background(), &proto.EmptyGigReq{})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(200).JSON(res)
}

// @Summary Process a refund for an order
// @Description Initiates a refund for a specified order by its OrderID.
// @Tags GetAllRefund
// @Accept json
// @Produce json
// @Param OrderID path string true "ID of the order to be refunded"
// @Router /admin/orders/refund/{OrderID} [post]
func (h *GigHandler) Refund(c *fiber.Ctx) error {
	OrderID := c.Params("OrderID")
	if OrderID == "" {
		fmt.Println(OrderID)
		return c.Status(400).JSON(fiber.Map{
			"error": "All fields are requered",
		})
	}
	res, err := h.GigClinet.AdOrderRefund(context.Background(), &proto.AdRefundReq{
		OrderId: OrderID,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(200).JSON(res)
}

// @Summary Get all complet orders
// @Description Retrieves a list of all complet orders in the system.
// @Tags GetAllCompleted
// @Accept json
// @Produce json
// @Router /admin/orders/completed [get]
func (h *GigHandler) GetAllCompletedOrders(c *fiber.Ctx) error {
	res, err := h.GigClinet.AdOrderCheck(context.Background(), &proto.EmptyGigReq{})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(200).JSON(res)
}

// @Summary Process a complet for an order
// @Description Initiates a complet for a specified order by its OrderID.
// @Tags GetAllComplet
// @Accept json
// @Produce json
// @Param OrderID path string true "ID of the order to be completed"
// @Router /admin/orders/complet/{OrderID} [post]
func (h *GigHandler) Payment(c *fiber.Ctx) error {
	OrderID := c.Params("OrderID")
	if OrderID == "" {
		fmt.Println(OrderID)
		return c.Status(400).JSON(fiber.Map{
			"error": "All fields are requered",
		})
	}
	res, err := h.GigClinet.AdOrderRefund(context.Background(), &proto.AdRefundReq{
		OrderId: OrderID,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(200).JSON(res)
}