package postgres

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/P3rCh1/immersive-images/backend/internal/common"
	"github.com/P3rCh1/immersive-images/backend/internal/testutils"
	sqlmigrate "github.com/rubenv/sql-migrate"
	"github.com/stretchr/testify/suite"
)

const (
	migrationsDir = "../../../migrations"
)

type Suite struct {
	suite.Suite
	log      *slog.Logger
	postgres *Postgres
}

func TestPostgres(t *testing.T) {
	suite.Run(t, &Suite{})
}

func (s *Suite) SetupSuite() {
	if !testutils.AcceptanceEnabed() {
		s.T().Skip("acceptance tests are disabled")
	}

	s.log = slog.New(slog.NewJSONHandler(
		os.Stdout,
		&slog.HandlerOptions{Level: slog.LevelDebug},
	))

	testutils.InitTestConfig()

	p, err := New(s.log)
	s.Require().NoError(err)

	s.postgres = p
}

func (s *Suite) SetupTest() {
	s.MigrateUp()
}

func (s *Suite) TearDownTest() {
	s.MigrateDown()
}

func (s *Suite) TearDownSuite() {
	if s.postgres == nil {
		return
	}

	s.postgres.Close()
}

func (s *Suite) MigrateUp() {
	s.T().Helper()
	s.migrate(sqlmigrate.Up)
}

func (s *Suite) MigrateDown() {
	s.T().Helper()
	s.migrate(sqlmigrate.Down)
}

func (s *Suite) migrate(direction sqlmigrate.MigrationDirection) {
	s.T().Helper()

	_, err := sqlmigrate.ExecContext(
		context.Background(),
		s.postgres.db.DB,
		"postgres",
		sqlmigrate.FileMigrationSource{Dir: migrationsDir},
		direction,
	)

	s.Require().NoError(err)
}

func GetTestImage() *Image {
	return &Image{
		Name:       "name",
		Size:       100,
		Seed:       101,
		Style:      "style",
		Palette:    "pallete",
		Additional: common.Ptr("additional"),
		Width:      102,
		Height:     103,
		Scale:      1.0,
	}
}
