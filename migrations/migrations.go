package migrations

// import (
// 	"bhbdev/jam/lib/database"
// 	"fmt"
// 	"log/slog"

// 	_ "github.com/go-sql-driver/mysql"
// 	"github.com/golang-migrate/migrate/v4"
// 	migrate_db "github.com/golang-migrate/migrate/v4/database"
// 	"github.com/golang-migrate/migrate/v4/database/mysql"
// 	"github.com/golang-migrate/migrate/v4/database/postgres"
// 	"github.com/golang-migrate/migrate/v4/database/sqlite3"
// 	_ "github.com/golang-migrate/migrate/v4/source/file"
// 	_ "github.com/mattn/go-sqlite3"
// )

// func  Migrate(db database.Database, direction string, steps int) (err error) {

// 	var driver migrate_db.Driver

// 	switch db.Engine() {
// 	case database.EngineMySQL:
// 		driver, err = mysql.WithInstance(db.Instance(), &mysql.Config{})
// 		if err != nil {
// 			slog.Error("failed to create mysql driver", "error", err)
// 			return err
// 		}
// 		break
// 	case database.EnginePostgres:
// 		driver, err = postgres.WithInstance(db.Instance(), &postgres.Config{})
// 		if err != nil {
// 			slog.Error("failed to create postgres driver", "error", err)
// 			return err
// 		}
// 		break
// 	case database.EngineSQLite:
// 		driver, err = sqlite3.WithInstance(db.Instance(), &sqlite3.Config{})
// 		if err != nil {
// 			slog.Error("failed to create sqlite3 driver", "error", err)
// 			return err
// 		}
// 	default:
// 		return fmt.Errorf("invalid engine: %s", db.Engine())
// 	}

// 	m, err := migrate.NewWithDatabaseInstance(
// 		"file://migrations",
// 		"jam", //db name
// 		driver,
// 	)
// 	if err != nil {
// 		slog.Error("failed to create migration instance", "error", err)
// 		return err
// 	}

// 	slog.Info("running migrations...")

// 	if direction == "down" {
// 		if steps == 0 {
// 			err = m.Down()
// 		} else {
// 			err = m.Steps(-steps)
// 		}
// 	} else { // up
// 		if steps == 0 {
// 			err = m.Up()
// 		} else {
// 			err = m.Steps(steps)
// 		}
// 	}

// 	switch err {
// 	case migrate.ErrNoChange:
// 		slog.Info("no change")
// 		break
// 	case migrate.ErrNilVersion:
// 		slog.Error("no migration")
// 		break
// 	case migrate.ErrInvalidVersion:
// 		slog.Error("version must be >= -1")
// 		break
// 	case migrate.ErrLocked:
// 		slog.Error("database locked")
// 		break
// 	case migrate.ErrLockTimeout:
// 		slog.Error("timeout: can't acquire database lock")
// 		break
// 	default:
// 		slog.Error("failed to run migrations", "error", err)
// 	}
// 	slog.Info("migrations complete")

// 	return nil
// }
//
