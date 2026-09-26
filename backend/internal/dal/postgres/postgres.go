package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/P3rCh1/immersive-images/backend/internal/config"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var ErrNotFound = errors.New("not found")

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

func (p *Postgres) Close() {
	if err := p.db.Close(); err != nil {
		p.log.Error(
			"failed to close db",
			"error", err,
		)
	}
}

func (p *Postgres) ListImages(ctx context.Context, limit int, start *uuid.UUID) ([]Image, error) {
	const query = `
		SELECT
			id,
			name,
			size,
			seed,
			style,
			palette,
			additional,
			width,
			height,
			scale,
			uploaded,
			created_at
		FROM images
		WHERE
			uploaded AND
			(
				$1::uuid IS NULL
				OR
				(created_at, id) <= (SELECT created_at, id FROM images WHERE id = $1::uuid)
			)
		ORDER BY created_at DESC, id DESC	
		LIMIT $2
	`

	handler := p.db.SelectContext

	tx, ok := ctx.Value(txCtxValueKey).(*sqlx.Tx)
	if ok {
		handler = tx.SelectContext
	}

	var images []Image
	if err := handler(ctx, &images, query, start, limit); err != nil {
		p.log.Error(
			"failed to get image list from DB",
			"error", err,
		)

		return nil, fmt.Errorf("failed to get image list from DB: %w", err)
	}

	return images, nil
}

func (p *Postgres) CreateImage(ctx context.Context, image *Image) error {
	const query = `
		INSERT INTO images (
			name,
			size,
			seed,
			style,
			palette,
			additional,
			width,
			height,
			scale
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at
	`

	handler := p.db.QueryRowx

	tx, ok := ctx.Value(txCtxValueKey).(*sqlx.Tx)
	if ok {
		handler = tx.QueryRowx
	}

	if err := handler(
		query,
		image.Name,
		image.Size,
		image.Seed,
		image.Style,
		image.Palette,
		image.Additional,
		image.Width,
		image.Height,
		image.Scale,
	).StructScan(image); err != nil {
		p.log.Error(
			"failed to create image in DB",
			"error", err,
		)

		return fmt.Errorf("failed to create image in DB: %w", err)
	}

	return nil
}

func (p *Postgres) GetImage(ctx context.Context, id uuid.UUID) (*Image, error) {
	const query = `
		SELECT
			id,
			name,
			size,
			seed,
			style,
			palette,
			additional,
			width,
			height,
			scale,
			uploaded,
			created_at
		FROM images
		WHERE id = $1 AND uploaded
	`

	handler := p.db.GetContext

	tx, ok := ctx.Value(txCtxValueKey).(*sqlx.Tx)
	if ok {
		handler = tx.GetContext
	}

	var image Image
	if err := handler(ctx, &image, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}

		p.log.Error(
			"failed to get image from DB",
			"error", err,
		)

		return nil, fmt.Errorf("failed to get image from DB: %w", err)
	}

	return &image, nil
}

func (p *Postgres) UpdateImage(ctx context.Context, image *Image) error {
	const query = `
		UPDATE images
		SET name = :name
		WHERE id = :id AND uploaded
	`

	handler := p.db.NamedExecContext

	tx, ok := ctx.Value(txCtxValueKey).(*sqlx.Tx)
	if ok {
		handler = tx.NamedExecContext
	}

	result, err := handler(ctx, query, image)
	if err != nil {
		p.log.Error(
			"failed to update image name in DB",
			"error", err,
		)

		return fmt.Errorf("failed to update image name in DB: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		p.log.Error(
			"failed to get affected rows",
			"error", err,
		)

		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if affected == 0 {
		return ErrNotFound
	}

	return nil
}

func (p *Postgres) SetUploaded(ctx context.Context, id uuid.UUID) error {
	const query = `
		UPDATE images
		SET uploaded = true
		WHERE id = $1
	`

	handler := p.db.ExecContext

	tx, ok := ctx.Value(txCtxValueKey).(*sqlx.Tx)
	if ok {
		handler = tx.ExecContext
	}

	result, err := handler(ctx, query, id)
	if err != nil {
		p.log.Error(
			"failed to set image uploaded flag",
			"error", err,
		)

		return fmt.Errorf("failed to set image uploaded flag: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		p.log.Error(
			"failed to get affected rows",
			"error", err,
		)

		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if affected == 0 {
		return ErrNotFound
	}

	return nil
}

func (p *Postgres) DeleteImage(ctx context.Context, id uuid.UUID) error {
	const query = `
		DELETE FROM images
		WHERE id = $1
	`

	handler := p.db.ExecContext

	tx, ok := ctx.Value(txCtxValueKey).(*sqlx.Tx)
	if ok {
		handler = tx.ExecContext
	}

	if _, err := handler(ctx, query, id); err != nil {
		p.log.Error(
			"failed to delete image from DB",
			"error", err,
		)

		return fmt.Errorf("failed to delete image from DB: %w", err)
	}

	return nil
}

func (p *Postgres) TotalImages(ctx context.Context) (int, error) {
	const query = `
		SELECT count(*) FROM images WHERE uploaded
	`

	handler := p.db.QueryRowx

	tx, ok := ctx.Value(txCtxValueKey).(*sqlx.Tx)
	if ok {
		handler = tx.QueryRowx
	}

	var count int
	if err := handler(query).Scan(&count); err != nil {
		p.log.Error(
			"failed to get image count from DB",
			"error", err,
		)

		return 0, fmt.Errorf("failed to get image count from DB: %w", err)
	}

	return count, nil
}
