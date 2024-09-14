package handler

import (
	"context"
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
func (h *UserHandler)AddSkill(c *fiber.Ctx)error{
	var req helper.AddSkill
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse JSON",
		})
	}
	res, err:=h.userClient.AddSkill(context.Background(),&proto.AddSkillReq{
		SkillName: req.Name,
	})
	if err!=nil{
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
func (h *UserHandler)GetCategory(c *fiber.Ctx)error{
	res,err:=h.userClient.GetCategory(context.Background(),&proto.EmtpyReq{})
	if err !=nil{
		return c.Status(400).JSON(fiber.Map{
			"error":err.Error(),
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"data":res,
		"status":200,
	})
}

// @Summary Get all skills
// @Description Retrieves a list of all skills
// @Tags Admin
// @Accept json
// @Produce json
// @Router /admin/skills [get]
func (h *UserHandler)GetSkills(c *fiber.Ctx)error{
	res ,err :=h.userClient.GetSkill(context.Background(),&proto.EmtpyReq{})
	if err !=nil{
		return c.Status(400).JSON(fiber.Map{
			"error":err.Error(),
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"data":res,
		"status":200,
	})
}

// @Summary Admin delete skill
// @Description Deletes a skill by ID, used by admin
// @Tags Admin
// @Param skillID path int true "Skill ID"
// @Accept json
// @Produce json
// @Router /admin/skills/{skillID} [delete]
func (h *UserHandler)AdDeleteSkill(c *fiber.Ctx)error{
	id:=c.Params("skillID")
	skillID,_:=strconv.Atoi(id)
	_,err:=h.userClient.AdDeleteSkill(context.Background(),&proto.ADeleteSkillReq{Id: uint32(skillID)})
	if err != nil{
		return c.Status(400).JSON(fiber.Map{
			"error":err.Error(),
		})
	}
	return c.Status(200).JSON(fiber.Map{
		 "message": "Skill successfully deleted",
		 "status":200,
	})

}

// @Summary Admin delete skill
// @Description Deletes a category by ID, used by admin
// @Tags Admin
// @Param CatID path int true "Category ID"
// @Accept json
// @Produce json
// @Router /admin/category/{CatID} [delete]
func (h *UserHandler)AdDeleteCategory(c *fiber.Ctx)error{
	id:=c.Params("CatID")
	skillID,_:=strconv.Atoi(id)
	_,err:=h.userClient.DeleteCategory(context.Background(),&proto.DeleteCatReq{Id: uint32(skillID)})
	if err != nil{
		return c.Status(400).JSON(fiber.Map{
			"error":err.Error(),
		})
	}
	return c.Status(200).JSON(fiber.Map{
		 "message": "Category successfully deleted",
		 "status":200,
	})

}

// @Summary Get all users
// @Description Retrieve a list of all users
// @Tags Admin
// @Accept json
// @Produce json
// @Router /admin/users [get]
func (h *UserHandler)GetAllUsers(c *fiber.Ctx)error{
	res, err:=h.userClient.GetAllUsers(context.Background(),&proto.EmtpyReq{})
	if err !=nil{
		return c.Status(400).JSON(fiber.Map{
			"error":err.Error(),
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"data":res,
		"status":200,
	})

}