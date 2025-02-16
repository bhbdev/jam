package app

import (
	"context"
	"log"
	"log/slog"
	"os"

	"bhbdev/jam/config"
	"bhbdev/jam/lib/database"
	"bhbdev/jam/lib/logger"
	"bhbdev/jam/lib/middleware"
	"bhbdev/jam/lib/server"
	"bhbdev/jam/router"
	"bhbdev/jam/services"
	//"bhbdev/jam/migrations"
)

func ServerApp(cfg *config.Config) {

	ctx := context.Background()

	db, err := database.Connect(&cfg.Database)
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}

	userService := services.NewUserService(db.Instance())
	jobAppService := services.NewJobAppService(db.Instance())
	services := services.NewServices(userService, jobAppService)

	// create new http server app
	app := server.NewApp(&cfg.Server)

	// add http handlers
	router.Apply(ctx, app)
	router.Assets(app, "assets")

	// add middleware
	app.AddMiddleware(middleware.NewTemplate("templates"))
	app.AddMiddleware(middleware.NewContentType())
	app.AddMiddleware(middleware.NewDatabase(db))
	app.AddMiddleware(middleware.NewPageServices(services))
	app.AddMiddleware(middleware.NewLogger(logger.Logger())) // keep last

	slog.Debug("starting 🍇 JAM server", "address", "http://"+cfg.Server.Address())

	log.Fatal(app.ListenAndServe())

}
