package postgres

import (
	"time"

	"github.com/google/uuid"
)

type Image struct {
	ID         uuid.UUID `db:"id"`
	Name       string    `db:"name"`
	Size       int64     `db:"size"`
	Seed       int       `db:"seed"`
	Style      string    `db:"style"`
	Palette    string    `db:"palette"`
	Additional *string   `db:"additional"`
	Width      int       `db:"width"`
	Height     int       `db:"height"`
	Scale      float64   `db:"scale"`
	CreatedAt  time.Time `db:"created_at"`
	Uploaded   bool      `db:"uploaded"`
}
