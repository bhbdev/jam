package gorm_postgres

import (
	"fmt"
	"strconv"

	"bhbdev/jam/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(cfg *config.Database) *gorm.DB {

	// FIXME: enable logger again
	//logger := logger.Logger()
	//gormlogger := logger.With("component", "gorm")

	addressStr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Host,
		cfg.User,
		cfg.Pass,
		cfg.Name,
		strconv.Itoa(cfg.Port),
		cfg.SslMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		// FIXME: enable logger again
		//Logger: gormlogruslogger.NewGormLogrusLogger(logger.WithField("component", "gorm"), 10*time.Millisecond),
	})

	if err != nil {
		panic(fmt.Sprintf(`😫: Connection failed, check Postgres with address %s, error: %s `, addressStr, err.Error()))
	}

	return db

}
