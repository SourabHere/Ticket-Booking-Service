package handlers

import (
	"database/sql"
	"net/http"
	"time"

	db "sourabHere/ticketBooking/internal/db/sqlc"
	"sourabHere/ticketBooking/internal/token"
	"sourabHere/ticketBooking/internal/util"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type TicketHandler struct {
	store      *db.Store
	tokenMaker token.Maker
	config     util.Config
}

type CancelTicketRequest struct {
	TicketID int32 `params:"ticket_id" validate:"required"`
}

type reserveSeatRequest struct {
	RouteID int32 `json:"route_id" validate:"required"`
	BusID   int32 `json:"bus_id" validate:"required"`
	SeatID  int32 `json:"seat_id" validate:"required"`
}

type reserveSeatResponse struct {
	TicketID   int32     `json:"ticket_id"`
	BusID      int32     `json:"bus_id"`
	SeatID     int32     `json:"seat_id"`
	ReservedAt time.Time `json:"reserved_at"`
}

type ListUserTicketsResponse struct {
	TicketID      int32     `json:"ticket_id"`
	BusID         int32     `json:"bus_id"`
	SeatID        int32     `json:"seat_id"`
	ReservedAt    time.Time `json:"reserved_at"`
	DepartureTime time.Time `json:"departure_time"`
	ArrivalTime   time.Time `json:"arrival_time"`
	Price         int       `json:"price"`
	SeatNumber    int       `json:"seat_number"`
	Status        string    `json:"status"`
}

func NewTicketHandler(Store *db.Store, tokenMaker token.Maker, Config util.Config) *TicketHandler {
	return &TicketHandler{
		Store,
		tokenMaker,
		Config,
	}
}

func (h *TicketHandler) ReserveSeat(c *fiber.Ctx) error {

	var req reserveSeatRequest
	if err := c.BodyParser(&req); err != nil {
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
			return fiber.NewError(http.StatusNotFound, "Bus or Route not found or they do not match")
		}
		return fiber.NewError(http.StatusInternalServerError, "Failed to validate bus and route association")
	}

	seat, err := h.store.CheckSeatAvailability(c.Context(), db.CheckSeatAvailabilityParams{
		ID:    req.SeatID,
		BusID: req.BusID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return fiber.NewError(http.StatusNotFound, "Seat not found")
		}
		return fiber.NewError(http.StatusInternalServerError, "Failed to fetch seat information")
	}

	if seat.Status != "available" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Seat is not available for purchase"})
	}

	payload := c.Locals("authorizationPayloadKey").(*token.Payload)

	user, err := h.store.GetUserByUsername(c.Context(), payload.Username)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	reservation, err := h.store.ReserveTicketTx(c.Context(), db.ReserveTicketTxParams{
		UserID: user.ID,
		BusID:  req.BusID,
		SeatID: req.SeatID,
	})
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "Failed to reserve seat: "+err.Error())
	}

	return c.Status(http.StatusOK).JSON(reservation)
}

func (h *TicketHandler) PurchaseTicket(c *fiber.Ctx) error {

	var req reserveSeatRequest
	if err := c.BodyParser(&req); err != nil {
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
			return fiber.NewError(http.StatusNotFound, "Bus or Route not found or they do not match")
		}
		return fiber.NewError(http.StatusInternalServerError, "Failed to validate bus and route association")
	}

	seat, err := h.store.CheckSeatAvailability(c.Context(), db.CheckSeatAvailabilityParams{
		ID:    req.SeatID,
		BusID: req.BusID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return fiber.NewError(http.StatusNotFound, "Seat not found")
		}
		return fiber.NewError(http.StatusInternalServerError, "Failed to fetch seat information")
	}

	if seat.Status != "available" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Seat is not available for purchase"})
	}

	payload := c.Locals("authorizationPayloadKey").(*token.Payload)
	user, err := h.store.GetUserByUsername(c.Context(), payload.Username)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	result, err := h.store.PurchaseTicketTx(c.Context(), db.PurchaseTicketTxParams{
		UserID: user.ID,
		BusID:  req.BusID,
		SeatID: req.SeatID,
	})
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "Failed to purchase ticket")
	}

	response := reserveSeatResponse{
		TicketID:   result.TicketID,
		BusID:      result.BusID,
		SeatID:     result.SeatID,
		ReservedAt: result.ReservedAt,
	}

	return c.Status(http.StatusOK).JSON(response)
}

func (h *TicketHandler) ListUserTickets(c *fiber.Ctx) error {

	payload := c.Locals("authorizationPayloadKey").(*token.Payload)

	user, err := h.store.GetUserByUsername(c.Context(), payload.Username)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	tickets, err := h.store.ListUserTickets(c.Context(), user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch user tickets"})
	}

	var response []ListUserTicketsResponse
	for _, ticket := range tickets {
		response = append(response, ListUserTicketsResponse{
			TicketID:      ticket.TicketID,
			BusID:         ticket.BusID,
			SeatID:        ticket.SeatID,
			ReservedAt:    ticket.ReservedAt.Time,
			DepartureTime: ticket.DepartureTime,
			ArrivalTime:   ticket.ArrivalTime,
			Price:         int(ticket.Price),
			SeatNumber:    int(ticket.SeatNumber),
			Status:        ticket.ReservationStatus,
		})
	}

	return c.Status(http.StatusOK).JSON(response)
}

func (h *TicketHandler) CancelTicket(c *fiber.Ctx) error {
	var req CancelTicketRequest
	if err := c.ParamsParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request parameters"})
	}

	payload := c.Locals("authorizationPayloadKey").(*token.Payload)
	user, err := h.store.GetUserByUsername(c.Context(), payload.Username)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	ticket, err := h.store.GetTicketByID(c.Context(), req.TicketID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Ticket not found"})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch ticket details"})
	}

	if ticket.UserID != user.ID {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{"error": "You do not have permission to cancel this ticket"})
	}

	if ticket.Status != "canceled" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Ticket cannot be canceled"})
	}

	err = h.store.CancelTicketTx(c.Context(), db.CancelTicketParams{
		TicketID: ticket.ID,
		SeatID:   ticket.SeatReservationID,
		UserID:   user.ID,
	})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to cancel ticket"})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Ticket canceled successfully"})
}
