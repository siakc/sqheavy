package main

import (
	"github.com/gofiber/fiber/v3"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
	. "sqheavy/db"
	. "sqheavy/settings"
	"sqheavy/routes"
)

func main() {
	InitHeavy()
	// Initialize a new Fiber app
	app := fiber.New()
	routes.MountRoutes(app)

	const LISTEN_ADDRESS string = SERVER_HOST + ":" + SERVER_PORT
	log.Info().Msg("Starting sqheavy server at " + LISTEN_ADDRESS)

	if err := app.Listen(LISTEN_ADDRESS); err != nil {
		log.Fatal().Err(err).Msg("Error starting server :/")
	}
}
