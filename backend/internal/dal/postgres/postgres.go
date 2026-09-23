package postgres

import (
	"fmt"
	"log/slog"

	"github.com/P3rCh1/immersive-images/backend/internal/config"
	"github.com/jmoiron/sqlx"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Postgres struct {
	log *slog.Logger
	db  *sqlx.DB
}

func New(log *slog.Logger) (*Postgres, error) {
	db, err := sqlx.Connect(
		"pgx",
		fmt.Sprintf(
			"host=%s port=%s user=postgres password=%s dbname=%s",
			config.Config.DB.Host,
			config.Config.DB.Port,
			config.Config.DB.Password,
			config.Config.DB.Name,
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to init postgres: %w", err)
	}

	return &Postgres{
		log: log,
		db:  db,
	}, nil
}
