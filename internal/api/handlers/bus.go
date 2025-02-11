package handlers

import (
	"database/sql"
	"net/http"

	db "sourabHere/ticketBooking/internal/db/sqlc"
	"sourabHere/ticketBooking/internal/token"
	"sourabHere/ticketBooking/internal/util"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type BusHandler struct {
	store      *db.Store
	tokenMaker token.Maker
	config     util.Config
}

type listAvailableSeatsRequest struct {
	BusID   int32 `params:"bus_id" validate:"required"`
	RouteID int32 `params:"route_id" validate:"required"`
}

type seatResponse struct {
	SeatID     int32  `json:"seat_id"`
	SeatNumber int32  `json:"seat_number"`
	Status     string `json:"status"`
}

func NewBusHandler(Store *db.Store, tokenMaker token.Maker, Config util.Config) *BusHandler {
	return &BusHandler{
		Store,
		tokenMaker,
		Config,
	}
}

func (h *BusHandler) ListAvailableSeats(c *fiber.Ctx) error {

	var req listAvailableSeatsRequest
	if err := c.ParamsParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request parameters"})
	}

	validate := validator.New()
	if err := validate.Struct(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	_, err := h.store.CheckBusRouteAssociation(c.Context(), db.CheckBusRouteAssociationParams{
		ID:   req.RouteID,
		ID_2: req.BusID,
	})

	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Bus or Route not found or they do not match"})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to validate bus and route association"})
	}

	seats, err := h.store.GetAvailableSeatsForBus(c.Context(), db.GetAvailableSeatsForBusParams{
		RouteID: req.RouteID,
		BusID:   req.BusID,
	})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch available seats"})
	}

	var response []seatResponse
	for _, seat := range seats {
		response = append(response, seatResponse{
			SeatID:     seat.SeatID,
			SeatNumber: seat.SeatNumber,
			Status:     "available",
		})
	}

	return c.Status(http.StatusOK).JSON(response)
}
