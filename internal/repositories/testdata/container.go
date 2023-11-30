package testdata

import (
	"context"
	"fmt"
	"github.com/bonus2k/go-musthave-diploma-tpl/internal/migrations"
	"github.com/jmoiron/sqlx"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"time"
)

const (
	migrationsInitDB   = "file://testdata/migration/init_db"
	migrationsInitData = "file://testdata/migration/init_data"
)

type PostgresContainer struct {
	testcontainers.Container
	MappedPort string
	Host       string
}

func (container *PostgresContainer) InitData() (*sqlx.DB, error) {
	db, err := sqlx.Open("pgx", container.getDSN())
	if err != nil {
		log.Printf("err open connection to db %v", err)
		return nil, err
	}
	db.MustExec(`TRUNCATE public.schema_migrations RESTART IDENTITY CASCADE`)
	log.Println("init data", container.getDSN())
	err = migrations.Start(container.getDSN(), migrationsInitData)
	if err != nil {
		log.Printf("err migration db %v", err)
		return nil, err
	}
	return db, nil
}

func (container *PostgresContainer) InitDB(ctx context.Context) error {
	log.Println("init db", container.getDSN())
	if err := migrations.Start(container.getDSN(), migrationsInitDB); err != nil {
		log.Printf("err migration db: %v", err)
		return err
	}
	return nil
}

func NewPostgresContainer(ctx context.Context) (*PostgresContainer, error) {
	req := testcontainers.ContainerRequest{
		Env: map[string]string{
			"POSTGRES_USER":     "user",
			"POSTGRES_PASSWORD": "password",
			"POSTGRES_DB":       "postgres_test",
		},
		ExposedPorts: []string{"5432/tcp"},
		Image:        "postgres:alpine",
		WaitingFor: wait.ForExec([]string{"pg_isready"}).
			WithPollInterval(2 * time.Second).
			WithExitCodeMatcher(func(exitCode int) bool {
				return exitCode == 0
			}),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}
	host, err := container.Host(ctx)
	if err != nil {
		return nil, err
	}

	mappedPort, err := container.MappedPort(ctx, "5432")
	if err != nil {
		return nil, err
	}

	return &PostgresContainer{
		Container:  container,
		MappedPort: mappedPort.Port(),
		Host:       host,
	}, nil
}

func (container *PostgresContainer) getDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", "user", "password", container.Host, container.MappedPort, "postgres")
}
