package api

import (
	"fmt"

	db "sourabHere/ticketBooking/internal/db/sqlc"
	"sourabHere/ticketBooking/internal/token"
	"sourabHere/ticketBooking/internal/util"

	"github.com/gofiber/fiber/v2"
)

type Server struct {
	Config     util.Config
	Store      *db.Store
	TokenMaker token.Maker
	App        *fiber.App
}

func NewServer(config util.Config, store *db.Store) (*Server, error) {
	tokenMaker, err := token.NewJWTMaker(config.TOKENSECRETKEY)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	app := fiber.New()

	server := &Server{
		Config:     config,
		Store:      store,
		TokenMaker: tokenMaker,
		App:        app,
	}

	return server, nil

}

func (server *Server) Start(address string) error {
	return server.App.Listen(fmt.Sprintf(":%s", address))
}
