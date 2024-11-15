package handler

import (
	"context"
	"fmt"
	"strconv"

	"github.com/MuhammedAshifVnr/Gig_Space/api_gateway/utils/helper"
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
	title, err := helper.GetRequiredFormValue(c, "title")
	if err != nil {
		return err
	}
	description, err := helper.GetRequiredFormValue(c, "description")
	if err != nil {
		return err
	}
	category, err := helper.GetRequiredFormValue(c, "category")
	if err != nil {
		return err
	}
	deliveryStr, err := helper.GetRequiredFormValue(c, "delivery")
	if err != nil {
		return err
	}
	revisionsStr, err := helper.GetRequiredFormValue(c, "revisions")
	if err != nil {
		return err
	}
	priceS, err := helper.GetRequiredFormValue(c, "price")
	if err != nil {
		return err
	}
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
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
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
func (h *GigHandler) GetGigByUserID(c *fiber.Ctx) error {
	user := c.Locals("userID")
	userid, ok := user.(uint)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "userID is of invalid type",
		})
	}
	res, err := h.GigClinet.GetGigsByFreelancerID(context.Background(), &proto.GetGigsByFreelancerIDRequest{FreelancerId: uint64(userid)})
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

// @Summary Update an existing gig
// @Description Updates the details of an existing gig by ID, including title, description, category, delivery days, revisions, price, and images.
// @Tags Gigs
// @Accept multipart/form-data
// @Produce application/json
// @Param id path string true "Gig ID"
// @Param title formData string false "Title of the gig"
// @Param description formData string false "Description of the gig"
// @Param category formData string false "Category of the gig"
// @Param delivery formData int false "Delivery days"
// @Param revisions formData int false "Number of revisions"
// @Param price formData int false "Price of the gig"
// @Param images formData file false "Images for the gig (multiple files allowed)"
// @Router /gig/{id} [put]
func (h *GigHandler) UpdaeteGig(c *fiber.Ctx) error {
	Id := c.Params("id")
	title, err := helper.GetRequiredFormValue(c, "title")
	if err != nil {
		return err
	}

	description, err := helper.GetRequiredFormValue(c, "description")
	if err != nil {
		return err
	}

	category, err := helper.GetRequiredFormValue(c, "category")
	if err != nil {
		return err
	}

	deliveryStr, err := helper.GetRequiredFormValue(c, "delivery")
	if err != nil {
		return err
	}

	revisionsStr, err := helper.GetRequiredFormValue(c, "revisions")
	if err != nil {
		return err
	}

	priceS, err := helper.GetRequiredFormValue(c, "price")
	if err != nil {
		return err
	}
	sellerID := c.Locals("userID")
	userid, ok := sellerID.(uint)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "userID is of invalid type",
		})
	}
	price, _ := strconv.Atoi(priceS)
	GigID, _ := strconv.Atoi(Id)
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
	res, err := h.GigClinet.UpdateGigByID(context.Background(), &proto.UpdateGigRequest{
		Id:                uint64(GigID),
		Title:             title,
		Description:       description,
		Category:          category,
		UserId:            uint32(userid),
		Price:             float64(price),
		DeliveryDays:      int64(delivery),
		NumberOfRevisions: int64(revisions),
		Images:            imageBytesList,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(int(res.Status)).JSON(res)

}

// @Summary Delete an existing gig
// @Description Deletes a gig by its ID, ensuring the user is authorized to delete the gig.
// @Tags Gigs
// @Produce application/json
// @Param GigID path string true "Gig ID"
// @Router /gig/{GigID} [delete]
func (h *GigHandler) DeleteGig(c *fiber.Ctx) error {
	id := c.Params("GigID")
	gigID, _ := strconv.Atoi(id)
	userid, ok := c.Locals("userID").(uint)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "invalid userID type",
		})
	}

	res, err := h.GigClinet.DeleteGigByID(context.Background(), &proto.DeleteReq{
		GigId:  uint32(gigID),
		UserId: uint32(userid),
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(int(res.Status)).JSON(res)
}

// @Summary Create an order
// @Description Creates a new order for a specified gig using the gig ID.
// @Tags Orders
// @Accept json
// @Produce json
// @Param GigID path int true "Gig ID"
// @Router /gig/order/{GigID} [post]
func (h *GigHandler) CreateOrder(c *fiber.Ctx) error {
	id := c.Params("GigID")
	gigID, _ := strconv.Atoi(id)
	userid, _ := c.Locals("userID").(uint)
	res, err := h.GigClinet.CreateOrder(context.Background(), &proto.CreateOrderReq{
		ClinetId: uint32(userid),
		GigId:    uint32(gigID),
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(int(res.Status)).JSON(res)
}

// @Summary Create a quote for a Gig
// @Description Allows a client to create a quote for a specific gig by providing necessary details like price, message, and delivery days.
// @Tags Quotes
// @Accept multipart/form-data
// @Produce json
// @Param GigID path int true "Gig ID"
// @Param Message formData string true "Description or message for the quote"
// @Param Price formData int true "Price offered by the client"
// @Param DeliveryDays formData int true "Delivery days requested by the client"
// @Router /gig/quotes/{GigID} [post]
func (h *GigHandler) CreateQuote(c *fiber.Ctx) error {
	id := c.Params("GigID")
	gigID, _ := strconv.Atoi(id)
	fmt.Println(gigID)
	describe, err := helper.GetRequiredFormValue(c, "Message")
	if err != nil {
		return nil
	}
	price, _ := strconv.Atoi(c.FormValue("Price"))
	delivery, _ := strconv.Atoi(c.FormValue("DeliveryDays"))
	userid, _ := c.Locals("userID").(uint)
	fmt.Println(price, " ", describe)
	res, err := h.GigClinet.RequestQuote(context.Background(), &proto.QuoteReq{
		GigId:        uint64(gigID),
		ClientId:     uint64(userid),
		Describe:     describe,
		Price:        float64(price),
		DeliveryDays: int64(delivery),
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(int(res.Status)).JSON(res)
}

// @Summary Get all quotes for a user
// @Description This endpoint fetches all quotes for the authenticated user, either as a client or as a freelancer.
// @Tags Quotes
// @Accept  application/json
// @Produce application/json
// @Router /gig/quotes [get]
func (h *GigHandler) GetAllQuote(c *fiber.Ctx) error {
	userid, _ := c.Locals("userID").(uint)
	res, err := h.GigClinet.GetAllQuotes(context.Background(), &proto.GetAllQuoteReq{
		UserId: uint32(userid),
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(200).JSON(res)
}

// @Summary Create a custom gig
// @Description This endpoint allows freelancers to create a custom gig based on a client's request.
// @Tags Custom Gig
// @Accept  multipart/form-data
// @Produce application/json
// @Param gig_request_id formData int true "Gig Request ID"
// @Param client_id formData int true "Client ID"
// @Param title formData string true "Gig Title"
// @Param description formData string true "Gig Description"
// @Param price formData number true "Gig Price"
// @Param delivery_days formData int true "Delivery Days"
// @Router /gig/custom [post]
func (h *GigHandler) CreateCustomGig(c *fiber.Ctx) error {
	freelancerID, ok := c.Locals("userID").(uint)
	gigRequestID, err1 := strconv.Atoi(c.FormValue("gig_request_id"))
	clientID, err2 := strconv.Atoi(c.FormValue("client_id"))
	title, err := helper.GetRequiredFormValue(c, "title")
	if err != nil {
		return err
	}
	description, err := helper.GetRequiredFormValue(c, "description")
	if err != nil {
		return err
	}
	price, err3 := strconv.ParseFloat(c.FormValue("price"), 64)
	deliveryDays, err4 := strconv.Atoi(c.FormValue("delivery_days"))

	fmt.Println(gigRequestID, " ", clientID, " ", title)
	// Check if any of the required fields are missing or invalid in a single if condition
	if !ok || err1 != nil || err2 != nil || title == "" || description == "" || err3 != nil || err4 != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "All fields are required and must be valid",
		})
	}

	res, err := h.GigClinet.CreateOffer(context.Background(), &proto.CreateOfferReq{
		GigRequestId: uint64(gigRequestID),
		FreelancerId: uint64(freelancerID),
		ClientId:     uint64(clientID),
		Title:        title,
		Descripton:   description,
		Price:        float32(price),
		DeliveryDays: int64(deliveryDays),
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(200).JSON(res)
}

// @Summary Get all offers for a client
// @Description This endpoint fetches all offers sent to the authenticated client.
// @Tags Custom Gig
// @Accept  application/json
// @Produce application/json
// @Router /gig/offers [get]
func (h *GigHandler) GetAllOffers(c *fiber.Ctx) error {
	userid, _ := c.Locals("userID").(uint)
	res, err := h.GigClinet.GetAllOffers(context.Background(), &proto.GetAllOfferReq{
		ClientId: uint64(userid),
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(200).JSON(res)
}

// @Summary      Create an Offer Order
// @Description  Creates an order for a gig offer using the provided GigID and user ID from context.
// @Tags         Orders
// @Accept       json
// @Produce      json
// @Param        GigID   path      int   true  "Gig ID"
// @Router       /gig/offers/{GigID}/ [post]
func (h *GigHandler) CreateOfferOrder(c *fiber.Ctx) error {
	id := c.Params("GigID")
	gigID, _ := strconv.Atoi(id)
	userid, _ := c.Locals("userID").(uint)
	res, err := h.GigClinet.CreateOfferOrder(context.Background(), &proto.CreateOrderReq{
		ClinetId: uint32(userid),
		GigId:    uint32(gigID),
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(int(res.Status)).JSON(res)
}

// @Summary Get all order requests
// @Description Retrieve all requests for orders associated with the user
// @Tags Request
// @Accept json
// @Produce json
// @Router /gig/requests [get]
func (h *GigHandler) GetAllOrdersRequest(c *fiber.Ctx) error {
	userid, _ := c.Locals("userID").(uint)
	res, err := h.GigClinet.GetAllRequest(context.Background(), &proto.GetAllRequestReq{
		UserId: uint64(userid),
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(200).JSON(res)
}

// @Summary      Accept an order
// @Description  Accepts an order by updating its status to "Accepted" in the system
// @Tags         Request
// @Accept       json
// @Produce      json
// @Param        order_id   path      string  true  "Order ID"
// @Router       /gig/{order_id}/accept [post]
func (h *GigHandler) AccepteOrder(c *fiber.Ctx) error {
	orderID := c.Params("order_id")
	res, err := h.GigClinet.AcceptRequest(context.Background(), &proto.AcceptReq{OrderId: orderID})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(200).JSON(res)
}

// @Summary      Reject an order
// @Description  Reject an order by updating its status to "Reject" in the system
// @Tags         Request
// @Accept       json
// @Produce      json
// @Param        order_id   path      string  true  "Order ID"
// @Router       /gig/{order_id}/reject [post]
func (h *GigHandler) RejectOrder(c *fiber.Ctx) error {
	orderID := c.Params("order_id")
	res, err := h.GigClinet.RejectRequest(context.Background(), &proto.RejectReq{OrderId: orderID})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(200).JSON(res)
}

// @Summary      Get all orders for a freelancer
// @Description  Retrieve all orders associated with the authenticated freelancer
// @Tags         Orders
// @Accept       json
// @Produce      json
// @Router       /gig/orders/freelancer [get]
func (h *GigHandler) GetAllOrdersFreelancer(c *fiber.Ctx) error {
	userid, _ := c.Locals("userID").(uint)
	res, err := h.GigClinet.GetAllOrders(context.Background(), &proto.AllOrdersReq{UserId: uint64(userid)})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(200).JSON(res)
}

// @Summary      Get order by ID
// @Description  Retrieve details of a specific order using its unique ID
// @Tags         Orders
// @Accept       json
// @Produce      json
// @Param        order_id   path      string   true   "Order ID"
// @Router       /gig/orders/{order_id} [get]
func (h *GigHandler) GetOrderByID(c *fiber.Ctx) error {
	orderID := c.Params("order_id")
	res, err := h.GigClinet.GetOrderByID(context.Background(), &proto.OrderByIDReq{
		OrderId: orderID,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(200).JSON(res)
}

// @Summary      Update order status
// @Description  Update the status of a specific order by the client
// @Tags         Orders
// @Accept       application/x-www-form-urlencoded
// @Produce      json
// @Param        order_id   path      string   true    "Order ID"
// @Router       /gig/orders/{order_id}/pending [put]
func (h *GigHandler) ClientPendingUpdate(c *fiber.Ctx) error {
	orderID := c.Params("order_id")
	client := c.Locals("userID").(uint)
	res, err := h.GigClinet.ClientUpdatePendingStatus(context.Background(), &proto.OrderIDReq{
		OrderId:  orderID,
		ClientId: uint64(client),
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(200).JSON(res)
}

// @Summary      Update order status
// @Description  Update the status of a specific order by the client
// @Tags         Orders
// @Accept       application/x-www-form-urlencoded
// @Produce      json
// @Param        order_id   path      string   true    "Order ID"
// @Router       /gig/orders/{order_id}/done [put]
func (h *GigHandler) ClientDoneUpdate(c *fiber.Ctx) error {
	orderID := c.Params("order_id")
	client := c.Locals("userID").(uint)
	res, err := h.GigClinet.ClientUpdateDoneStatus(context.Background(), &proto.OrderIDReq{
		OrderId:  orderID,
		ClientId: uint64(client),
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(200).JSON(res)
}

// @Summary      Get All Gigs
// @Description  Retrieve all gigs excluding the specified user ID.
// @Tags         Gigs
// @Accept       json
// @Produce      json
// @Router       /gig/client/all    [get]
func (h *GigHandler) ClientGetAllGig(c *fiber.Ctx) error {
	client := c.Locals("userID").(uint)
	res, err := h.GigClinet.GetAllGig(context.Background(), &proto.GigReq{UserId: uint64(client)})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(200).JSON(res)
}

// @Summary      Get Gig by ID
// @Description  Retrieve a specific gig by its ID.
// @Tags         Gigs
// @Accept       json
// @Produce      json
// @Param        gig_id     path      int     true    "Gig ID"
// @Router       /gig/client/{gig_id} [get]
func (h *GigHandler) ClientGetGigByID(c *fiber.Ctx) error {
	gig_id, _ := strconv.Atoi(c.Params("gig_id"))
	res, err := h.GigClinet.GetGigByID(context.Background(), &proto.GigIDreq{
		GigId: uint64(gig_id),
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(200).JSON(res)
}
