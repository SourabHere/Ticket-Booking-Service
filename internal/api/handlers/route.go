package handlers

import (
	"net/http"
	"time"

	db "sourabHere/ticketBooking/internal/db/sqlc"
	"sourabHere/ticketBooking/internal/token"
	"sourabHere/ticketBooking/internal/util"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type RouteHandler struct {
	store      *db.Store
	tokenMaker token.Maker
	config     util.Config
}

type listRoutesRequest struct {
	OriginTerminalID      int32     `query:"origin_city_id" validate:"required"`
	DestinationTerminalID int32     `query:"destination_city_id" validate:"required"`
	DepartureTime         time.Time `query:"departure_time" validate:"required"`
}

type routeResponse struct {
	RouteID         int32     `json:"route_id"`
	BusID           int32     `json:"bus_id"`
	OriginCity      string    `json:"origin_city"`
	DestinationCity string    `json:"destination_city"`
	DepartureTime   time.Time `json:"departure_time"`
	ArrivalTime     time.Time `json:"arrival_time"`
	AvailableSeats  int64     `json:"available_seats"`
	Price           int32     `json:"price"`
}

func NewRouteHandler(store *db.Store, tokenMaker token.Maker, config util.Config) *RouteHandler {
	return &RouteHandler{
		store,
		tokenMaker,
		config,
	}
}

func (h *RouteHandler) SearchRoutes(c *fiber.Ctx) error {

	var req listRoutesRequest
	if err := c.QueryParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request parameters"})
	}

	var validate = validator.New()

	if err := validate.Struct(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	routes, err := h.store.ListRoutes(c.Context(), db.ListRoutesParams{
		OriginTerminalID:      req.OriginTerminalID,
		DestinationTerminalID: req.DestinationTerminalID,
		DepartureTime:         req.DepartureTime,
	})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch routes"})
	}

	var response []routeResponse
	for _, route := range routes {
		response = append(response, routeResponse{
			RouteID:         route.RouteID,
			BusID:           route.BusID,
			OriginCity:      route.OriginTerminalName,
			DestinationCity: route.DestinationTerminalName,
			DepartureTime:   route.DepartureTime,
			ArrivalTime:     route.ArrivalTime,
			AvailableSeats:  route.AvailableSeats,
			Price:           route.Price,
		})
	}

	return c.Status(http.StatusOK).JSON(response)
}
