package migrations

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/bonus2k/go-musthave-diploma-tpl/internal"
	"github.com/bonus2k/go-musthave-diploma-tpl/internal/utils"
	"github.com/golang-migrate/migrate/v4"
	mpgx "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"time"
)

const migrationsPath = "file://internal/migrations/sql"

func Start(connect string) error {
	dataBase, err := sql.Open("pgx", connect)
	if err != nil {
		return fmt.Errorf("can't open connection to db %w", err)
	}

	defer func() {
		if dataBase != nil {
			_ = dataBase.Close()
		}
	}()

	f := func() error { return establishConnection(dataBase) }
	if err = utils.RetryAfterError(f); err != nil {
		return fmt.Errorf("can't connected to db %w", err)
	}

	if err = migrateSQL(dataBase); err != nil {
		return err
	}
	internal.Log.Info("migration successfully finished")
	return nil
}

func establishConnection(db *sql.DB) error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancelFunc()
	return db.PingContext(ctx)
}

func migrateSQL(db *sql.DB) error {
	driver, err := mpgx.WithInstance(db, &mpgx.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		migrationsPath,
		"pgx", driver)
	if err != nil {
		return err
	}
	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	return nil
}
