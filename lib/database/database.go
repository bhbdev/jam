package database

import (
	"fmt"

	"github.com/bhbdev/jam/config"

	gorm_postgres "github.com/bhbdev/jam/lib/gorm_postgres"
	gorm_sqlite "github.com/bhbdev/jam/lib/gorm_sqlite"

	"gorm.io/gorm"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

const (
	EngineMySQL    = "mysql"
	EngineSQLite   = "sqlite3"
	EnginePostgres = "postgres"
)

type Database interface {
	Engine() string
	Instance() *gorm.DB
}

type DB struct {
	*gorm.DB
	engine string
}

func Connect(cfg *config.Database) (db *DB, err error) {
	db = &DB{
		engine: cfg.Engine,
	}
	switch db.engine {
	case EngineMySQL:
		// todo: implement mysql
	case EnginePostgres:
		db.DB = gorm_postgres.Connect(cfg)
	case EngineSQLite:
		db.DB = gorm_sqlite.Connect(cfg)
	default:
		return nil, fmt.Errorf("invalid engine: %s", cfg.Engine)
	}
	return
}

func (db DB) Instance() *gorm.DB {
	return db.DB
}

func (db DB) Engine() string {
	return db.engine
}
