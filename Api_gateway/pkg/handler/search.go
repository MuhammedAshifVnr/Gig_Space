package handler

import (
	"context"
	"log"
	"strconv"

	"github.com/MuhammedAshifVnr/Gig_Space_Proto/proto"
	"github.com/gofiber/fiber/v2"
)

type SearchHandler struct {
	SearchClient proto.SearchServiceClient
}

func NewSearchHandler(sClient proto.SearchServiceClient) *SearchHandler {
	return &SearchHandler{
		SearchClient: sClient,
	}
}

// @Summary Search gigs
// @Description This endpoint allows searching for gigs based on a query string, price range, revisions, and delivery days.
// @Tags Search
// @Accept  json
// @Produce application/json
// @Param query query string false "Search query to match title, description, or category"
// @Param price_upto query float64 false "Maximum price filter"
// @Param revisions_min query int false "Minimum number of revisions filter"
// @Param delivery_days query int false "Maximum delivery days filter"
// @Router /search/gigs [get]
func (h *SearchHandler) SearchGig(c *fiber.Ctx) error {
	query := c.Query("query", "")

	price_upto, _ := strconv.ParseFloat(c.Query("price_upto", "0"), 64)
	revisions_min, _ := strconv.Atoi(c.Query("revisions_min", "0"))
	deliveryDays, _ := strconv.Atoi(c.Query("delivery_days", "0"))

	res, err := h.SearchClient.SearchGig(context.Background(), &proto.SearchGigReq{
		Query:            query,
		PriceUpto:       float32(price_upto),
		RevisionsMin:    int32(revisions_min),
		DeliveryDaysMax: int32(deliveryDays),
	})
	if err != nil {
		log.Println("Error Calling Search Service: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Search service failed"})
	}
	return c.Status(200).JSON(res)
}
