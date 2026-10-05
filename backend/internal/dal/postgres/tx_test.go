package postgres

import (
	"context"
	"database/sql"
)

func (s *Suite) TestBeginX() {
	ctx := context.Background()

	s.log.Debug("Begin transaction")
	ctxWithTx, err := s.postgres.Beginx(ctx)
	s.Require().NoError(err)

	img := GetTestImage()

	s.log.Debug("Insert image")
	err = s.postgres.CreateImage(ctxWithTx, img)
	s.Require().NoError(err)

	s.log.Debug("Set uploaded flag")
	err = s.postgres.SetUploaded(ctxWithTx, img.ID)
	s.Require().NoError(err)

	s.log.Debug("Commit transaction")
	err = s.postgres.Commit(ctxWithTx)
	s.Require().NoError(err)

	s.log.Debug("Check if image exists")
	inserted, err := s.postgres.GetImage(ctx, img.ID)
	s.Require().NoError(err)
	s.Require().Equal(img.Name, inserted.Name)
}

func (s *Suite) TestBeginTxx() {
	ctx := context.Background()

	s.log.Debug("Begin transaction")
	ctxWithTx, err := s.postgres.BeginTxx(ctx, &sql.TxOptions{})
	s.Require().NoError(err)

	img := GetTestImage()

	s.log.Debug("Insert image")
	err = s.postgres.CreateImage(ctxWithTx, img)
	s.Require().NoError(err)

	s.log.Debug("Set uploaded flag")
	err = s.postgres.SetUploaded(ctxWithTx, img.ID)
	s.Require().NoError(err)

	s.log.Debug("Commit transaction")
	err = s.postgres.Commit(ctxWithTx)
	s.Require().NoError(err)

	s.log.Debug("Check if image exists")
	inserted, err := s.postgres.GetImage(ctx, img.ID)
	s.Require().NoError(err)
	s.Require().Equal(img.Name, inserted.Name)
}

func (s *Suite) TestCommit() {
	ctx := context.Background()

	s.log.Debug("Begin transaction")
	ctxWithTx, err := s.postgres.BeginTxx(ctx, &sql.TxOptions{})
	s.Require().NoError(err)

	img := GetTestImage()

	s.log.Debug("Insert image")
	err = s.postgres.CreateImage(ctxWithTx, img)
	s.Require().NoError(err)

	s.log.Debug("Set uploaded flag")
	err = s.postgres.SetUploaded(ctxWithTx, img.ID)
	s.Require().NoError(err)

	s.log.Debug("Commit transaction")
	err = s.postgres.Commit(ctxWithTx)
	s.Require().NoError(err)

	s.log.Debug("Check if image exists")
	inserted, err := s.postgres.GetImage(ctx, img.ID)
	s.Require().NoError(err)
	s.Require().Equal(img.Name, inserted.Name)

	s.log.Debug("Try commit again")
	err = s.postgres.Commit(ctxWithTx)
	s.Require().Error(err)

	s.log.Debug("Try commit without begin")
	err = s.postgres.Commit(ctx)
	s.Require().ErrorIs(err, ErrNoTx)
}

func (s *Suite) TestRollback() {
	ctx := context.Background()

	s.log.Debug("Begin transaction")
	ctxWithTx, err := s.postgres.BeginTxx(ctx, &sql.TxOptions{})
	s.Require().NoError(err)

	img := GetTestImage()

	s.log.Debug("Insert image")
	err = s.postgres.CreateImage(ctxWithTx, img)
	s.Require().NoError(err)

	s.log.Debug("Set uploaded flag")
	err = s.postgres.SetUploaded(ctxWithTx, img.ID)
	s.Require().NoError(err)

	s.log.Debug("Rollback transaction")
	s.postgres.Rollback(ctxWithTx)

	s.log.Debug("Check if image not exists")
	_, err = s.postgres.GetImage(ctx, img.ID)
	s.Require().ErrorIs(err, ErrNotFound)
}
