package handler

import (
	"context"
	"fmt"
	"strconv"

	"github.com/MuhammedAshifVnr/Gig_Space_Proto/proto"
	"github.com/gofiber/fiber/v2"
)

type GigHandler struct {
	GigClinet proto.GigServiceClient
}

func NewGigHandler(gigConn proto.GigServiceClient) *GigHandler {
	return &GigHandler{
		GigClinet: gigConn,
	}
}

// @Summary Create a new gig
// @Description This endpoint creates a new gig with a title, description, price, and images. Images are uploaded via multipart form.
// @Tags Gigs
// @Accept  multipart/form-data
// @Produce application/json
// @Param title formData string true "Title of the gig"
// @Param description formData string true "Description of the gig"
// @Param category formData string true "Category of the gig"
// @Param delivery formData string true "Number of delivery days"
// @Param revisions formData int true "Number of revisions"
// @Param price formData string true "Price of the gig"
// @Param images formData []file true "Images for the gig (can upload multiple images)"
// @Router /gig/add [post]
func (h *GigHandler) CreateGig(c *fiber.Ctx) error {
	title := c.FormValue("title")
	description := c.FormValue("description")
	category := c.FormValue("category")
	deliveryStr := c.FormValue("delivery")
	revisionsStr := c.FormValue("revisions")
	priceS := c.FormValue("price")
	sellerID := c.Locals("userID")
	userid, ok := sellerID.(uint)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "userID is of invalid type",
		})
	}
	price, _ := strconv.Atoi(priceS)
	delivery, _ := strconv.Atoi(deliveryStr)
	revisions, _ := strconv.Atoi(revisionsStr)
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Failed to upload images")
	}
	images := form.File["images"] // Get all the uploaded images

	var imageBytesList [][]byte

	for _, fileHeader := range images {
		file, err := fileHeader.Open()
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Failed to read image file")
		}
		defer file.Close()

		// Convert image to bytes
		imageBytes := make([]byte, fileHeader.Size)
		file.Read(imageBytes)

		// Append to the list
		imageBytesList = append(imageBytesList, imageBytes)
	}
	fmt.Println("price =", delivery)
	_, err = h.GigClinet.CreateGig(context.Background(), &proto.CreateGigReq{
		Title:             title,
		Description:       description,
		UserId:            uint32(userid),
		Price:             float64(price),
		Images:            imageBytesList,
		Category:          category,
		DeliveryDays:      int64(delivery),
		NumberOfRevisions: int64(revisions),
	})
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err,
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"message": "Gig add",
		"status":  200,
	})
}

// @Summary Get Gigs by User ID
// @Description Get all gigs created by the logged-in user
// @Tags Gigs
// @Accept json
// @Produce json
// @Router /gig/user [get]
func (h *GigHandler)GetGigByUserID(c *fiber.Ctx)error{
	user := c.Locals("userID")
	userid, ok := user.(uint)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "userID is of invalid type",
		})
	}
	res,err:=h.GigClinet.GetAllGigByID(context.Background(),&proto.GetAllGigsByIDReq{Id: uint32(userid)})
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