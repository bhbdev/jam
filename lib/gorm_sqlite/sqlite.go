package gorm_sqlite

import (
	"fmt"

	"bhbdev/jam/config"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Connect(config *config.Database) *gorm.DB {

	// FIXME: enable logger again
	//logger := logger.Logger.WithField("component", "gorm")

	dsn := fmt.Sprintf("file:%s?_foreign_keys=true", config.Name)

	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		// FIXME: enable logger again
		//Logger: gormlogruslogger.NewGormLogrusLogger(logger.WithField("component", "gorm"), 10*time.Millisecond),
	})

	if err != nil {
		panic(fmt.Sprintf(`😫: Connection failed, check your sqlite db file: %s.db, error: %s `, config.Name, err.Error()))
	}

	return db

}
