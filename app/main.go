package app

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/bhbdev/jam/config"
	"github.com/bhbdev/jam/lib/database"
	"github.com/bhbdev/jam/lib/logger"
	"github.com/bhbdev/jam/lib/middleware"
	"github.com/bhbdev/jam/lib/redis"
	"github.com/bhbdev/jam/lib/server"
	"github.com/bhbdev/jam/lib/session"
	"github.com/bhbdev/jam/router"
	"github.com/bhbdev/jam/services"
	//"bhbdev/jam/migrations"
)

func ServerApp(cfg *config.Config) {

	ctx := context.Background()

	db, err := database.Connect(&cfg.Database)
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}

	err = redis.Setup(cfg.Redis.Address())
	if err != nil {
		logger.Error("failed to connect to redis", "error", err)
		os.Exit(1)
	}
	session := session.NewSessionHandler(redis.Client)

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
	app.AddMiddleware(middleware.NewUserSessionHandler(session, userService))
	app.AddMiddleware(middleware.NewContentType())
	app.AddMiddleware(middleware.NewDatabase(db))
	app.AddMiddleware(middleware.NewPageServices(services))
	app.AddMiddleware(middleware.NewLogger(logger.Logger())) // keep last

	slog.Debug("starting 🍇 JAM server", "address", "http://"+cfg.Server.Address())

	log.Fatal(app.ListenAndServe())

}
