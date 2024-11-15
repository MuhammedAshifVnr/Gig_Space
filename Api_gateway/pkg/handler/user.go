package handler

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/MuhammedAshifVnr/Gig_Space/api_gateway/utils/helper"
	"github.com/MuhammedAshifVnr/Gig_Space_Proto/proto"
	"github.com/gofiber/fiber/v2"
)

// UserHandler handles user-related requests
type UserHandler struct {
	userClient proto.UserServiceClient
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(userConn proto.UserServiceClient) *UserHandler {
	return &UserHandler{
		userClient: userConn,
	}
}

// Signup godoc
// @Summary Sign up a new user
// @Description Create a new user account
// @Tags User
// @Accept  json
// @Produce  json
// @Param  request body helper.SignupData true "Signup Request"
// @Param role query string true "User role" Enums(freelancer, client)
// @Router /user/signup [post]
func (h *UserHandler) Signup(c *fiber.Ctx) error {
	var req helper.SignupData
	role := c.Query("role")
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse JSON",
		})
	}

	if validationErrors, err := helper.ValidateRequest(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":  "validation failed",
			"fields": validationErrors,
		})
	}

	hashPassword, err := helper.HashPassword(req.Password)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	resp, err := h.userClient.UserSignup(context.Background(), &proto.SignupReq{
		Firstname: req.Firstname,
		Lastname:  req.Lastname,
		Email:     req.Email,
		Password:  hashPassword,
		Country:   req.Country,
		Phone:     req.Phone,
		Role:      role,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create user",
			"e":     err.Error(),
		})
	}

	return c.Status(200).JSON(resp)
}

func (h *UserHandler) VerifyingEmail(c *fiber.Ctx) error {
	otp := c.Query("otp")
	email := c.Query("email")
	res, err := h.userClient.VerifyingEmail(context.Background(), &proto.VerifyReq{
		Otp:   otp,
		Email: email,
	})
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(200).JSON(res)
}

// Login godoc
// @Summary      User login
// @Description  Authenticate user with email and password
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        credentials  body      helper.LoginData  true  "Login credentials"
// @Router      /user/login [post]
func (h *UserHandler) Login(c *fiber.Ctx) error {
	var credentials helper.LoginData
	if err := c.BodyParser(&credentials); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse JSON",
		})
	}
	if validationErrors, err := helper.ValidateRequest(&credentials); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":  "validation failed",
			"fields": validationErrors,
		})
	}

	res, err := h.userClient.Login(context.Background(), &proto.LoginReq{
		Email:    credentials.Email,
		Password: credentials.Password,
	})
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	c.Cookie(&fiber.Cookie{
		Name:     "UserToken",
		Value:    res.Token,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
	})
	return c.Status(int(res.Status)).JSON(res)
}

// @Summary Add freelancer skill
// @Description This endpoint allows a freelancer to add a skill and set their proficiency level.
// @Tags User
// @Produce application/json
// @Param skillName formData string true "Skill name" example "Go Programming"
// @Param proficency formData int true "Proficiency level" example 5
// @Router /user/skills [post]
func (h *UserHandler) FreeAddSkills(c *fiber.Ctx) error {
	skill, err := helper.GetRequiredFormValue(c, "skillName")
	fmt.Println("Skill", skill)
	if err != nil {
		return err
	}
	proficency := c.FormValue("proficency")
	userIDLocal := c.Locals("userID")
	userid, ok := userIDLocal.(uint)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "userID is of invalid type",
		})
	}
	proficiencyLevel, _ := strconv.Atoi(proficency)

	res, err := h.userClient.FreelacerAddSkill(context.Background(), &proto.FreeAddSkillsReq{
		UserId:           uint32(userid),
		SkillName:        skill,
		ProficiencyLevel: int32(proficiencyLevel),
	})
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(int(res.Status)).JSON(res)
}

// @Summary      Update user profile
// @Description  Update the user's bio and title in their profile
// @Tags         User
// @Accept       multipart/form-data
// @Produce      json
// @Param        Bio    formData  string  true  "User bio"
// @Param        Title  formData  string  true  "User title"
// @Router       /user/profile [put]
func (h *UserHandler) UpdateProfile(c *fiber.Ctx) error {
	bio, err := helper.GetRequiredFormValue(c, "Bio")
	if err != nil {
		return err
	}
	title, err := helper.GetRequiredFormValue(c, "Title")
	if err != nil {
		return err
	}
	userIDLocal := c.Locals("userID")
	userid, ok := userIDLocal.(uint)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "userID is of invalid type",
		})
	}
	res, err := h.userClient.UpdateBio(context.Background(), &proto.UpdateProfileReq{
		UserId:      uint32(userid),
		Title:       title,
		Description: bio,
	})
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(200).JSON(res)
}

// @Summary Get user profile
// @Description Retrieves the profile details of the user based on their user ID
// @Tags User
// @Accept json
// @Produce json
// @Router /user/profile [get]
func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	userIDLocal := c.Locals("userID")
	userid, ok := userIDLocal.(uint)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "userID is of invalid type",
		})
	}
	res, err := h.userClient.GetUserProfile(context.Background(), &proto.ProfileReq{UserId: uint32(userid)})
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(200).JSON(res)
}

// @Summary Delete a skill from a user
// @Description Deletes a specific skill for a user based on the user ID and skill ID
// @Tags User
// @Accept json
// @Produce json
// @Param Skill path int true "Skill ID to delete"
// @Router /user/skill/{Skill} [delete]
func (h *UserHandler) DeleteSkill(c *fiber.Ctx) error {
	userIDLocal := c.Locals("userID")
	userid, ok := userIDLocal.(uint)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "userID is of invalid type",
		})
	}
	skill := c.Params("Skill")
	skillId, _ := strconv.Atoi(skill)
	res, err := h.userClient.DeleteSkill(context.Background(), &proto.DeleteSkillRes{
		UserId:  uint32(userid),
		SkillId: uint32(skillId),
	})
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(int(res.Status)).JSON(res)
}

// @Summary Uploads a profile photo for the user
// @Description Uploads a profile photo for the user based on the userID.
// @Tags User
// @Accept multipart/form-data
// @Produce json
// @Param photo formData file true "Profile photo"
// @Router /user/profile-photo [post]
func (h *UserHandler) UploadProfilePhoto(c *fiber.Ctx) error {
	user := c.Locals("userID")
	userid, ok := user.(uint)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "userID is of invalid type",
		})
	}
	file, err := c.FormFile("photo")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "file upload error",
		})
	}
	openfile, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "cannot open uploaded file",
		})
	}
	defer openfile.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, openfile)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "file read error",
		})
	}
	photo := buf.Bytes()
	res, err := h.userClient.UploadProfilePhoto(context.Background(), &proto.UpProilePicReq{
		UserId: uint32(userid),
		Pic:    photo,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to upload photo")
	}
	return c.Status(int(res.Status)).JSON(res)
}

// @Summary Forgot Password
// @Description Sends a forgot password request to the user's email address. The user will receive an OTP to reset the password.
// @Tags User
// @Accept application/x-www-form-urlencoded
// @Produce application/json
// @Param email formData string true "User Email"
// @Router /user/forgot-password [post]
func (h *UserHandler) ForgotPasswordReq(c *fiber.Ctx) error {
	email, err := helper.GetRequiredFormValue(c, "email")
	if err != nil {
		return err
	}
	res, err := h.userClient.ForgotPassword(context.Background(), &proto.FP_Req{
		Email: email,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(int(res.Status)).JSON(res)
}

// @Summary Reset Password
// @Description Allows the user to reset their password using the OTP and new password values.
// @Tags User
// @Accept application/x-www-form-urlencoded
// @Produce application/json
// @Param otp formData string true "OTP sent to user email"
// @Param pwd1 formData string true "New password"
// @Param pwd2 formData string true "Confirm new password"
// @Router /user/reset-password [post]
func (h *UserHandler) ResetPassword(c *fiber.Ctx) error {
	otp, err := helper.GetRequiredFormValue(c, "otp")
	if err != nil {
		return err
	}
	pwd1, err := helper.GetRequiredFormValue(c, "pwd1")
	if err != nil {
		return err
	}
	pwd2, err := helper.GetRequiredFormValue(c, "pwd2")
	if err != nil {
		return err
	}
	if pwd1 != pwd2 {
		return c.Status(400).JSON(fiber.Map{
			"error": "Password Doesn't Match.",
		})
	}
	hashPassword, err := helper.HashPassword(pwd2)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	res, err := h.userClient.ResetPassword(context.Background(), &proto.ResetPwdReq{
		Otp:      otp,
		Password: hashPassword,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(int(res.Status)).JSON(res)
}

// @Summary Update user address
// @Tags User
// @Description Update the user's address (city, district, state)
// @Param city formData string false "City"
// @Param district formData string false "District"
// @Param state formData string false "State"
// @Router /user/address [post]
func (h *UserHandler) UpdateAddress(c *fiber.Ctx) error {
	city, err := helper.GetRequiredFormValue(c, "city")
	if err != nil {
		return err
	}
	district, err := helper.GetRequiredFormValue(c, "district")
	if err != nil {
		return err
	}
	state, err := helper.GetRequiredFormValue(c, "state")
	if err != nil {
		return err
	}
	userid, ok := c.Locals("userID").(uint)
	fmt.Println("userid= ", userid)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "invalid userID type",
		})
	}
	res, err := h.userClient.UpdateAddress(context.Background(), &proto.AddressReq{
		Id:       uint32(userid),
		City:     city,
		State:    state,
		District: district,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(int(res.Status)).JSON(res)
}

// @Summary Change the role of a user
// @Description Updates the role of the user based on the user ID in the request context.
// @Tags User
// @Accept json
// @Produce json
// @Router /user/role [put]
func (h *UserHandler) RoleChange(c *fiber.Ctx) error {
	userID, _ := c.Locals("userID").(uint)
	res, err := h.userClient.UpdatRole(context.Background(), &proto.RoleReq{
		Id: uint32(userID),
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(int(res.Status)).JSON(res)
}
