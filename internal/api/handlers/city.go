package handlers

import (
	"net/http"

	db "sourabHere/ticketBooking/internal/db/sqlc"
	"sourabHere/ticketBooking/internal/token"
	"sourabHere/ticketBooking/internal/util"

	"github.com/gofiber/fiber/v2"
)

type CityHandler struct {
	store      *db.Store
	tokenMaker token.Maker
	config     util.Config
}

func NewCityHandler(Store *db.Store, tokenMaker token.Maker, Config util.Config) *CityHandler {
	return &CityHandler{
		Store,
		tokenMaker,
		Config,
	}
}

func (h *CityHandler) ListCities(c *fiber.Ctx) error {
	cities, err := h.store.GetAllCities(c.Context())
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch cities",
		})
	}

	return c.Status(http.StatusOK).JSON(cities)
}
