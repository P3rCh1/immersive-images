package postgres

import (
	"context"

	"github.com/google/uuid"
)

func (s *Suite) TestCreateImage() {
	ctx := context.Background()

	s.log.Debug("Create image")
	img := GetTestImage()
	err := s.postgres.CreateImage(ctx, img)
	s.Require().NoError(err)
	s.Require().NotEqual(uuid.Nil, img.ID)

	s.log.Debug("Check created image is not uploaded yet")
	_, err = s.postgres.GetImage(ctx, img.ID)
	s.Require().ErrorIs(err, ErrNotFound)

	s.log.Debug("Set uploaded flag")
	err = s.postgres.SetUploaded(ctx, img.ID)
	s.Require().NoError(err)

	s.log.Debug("Check uploaded image")
	created, err := s.postgres.GetImage(ctx, img.ID)
	s.Require().NoError(err)
	s.Require().Equal(img.ID, created.ID)
	s.Require().Equal(img.Name, created.Name)
	s.Require().Equal(img.Size, created.Size)
	s.Require().Equal(img.Seed, created.Seed)
	s.Require().Equal(img.Style, created.Style)
	s.Require().Equal(img.Palette, created.Palette)
	s.Require().Equal(img.Additional, created.Additional)
	s.Require().Equal(img.Width, created.Width)
	s.Require().Equal(img.Height, created.Height)
	s.Require().Equal(img.Scale, created.Scale)
	s.Require().True(created.Uploaded)

	s.log.Debug("Set uploaded flag for unknown image")
	err = s.postgres.SetUploaded(ctx, uuid.New())
	s.Require().ErrorIs(err, ErrNotFound)
}

func (s *Suite) TestGetImage() {
	ctx := context.Background()

	img := GetTestImage()

	s.log.Debug("Create image")
	err := s.postgres.CreateImage(ctx, img)
	s.Require().NoError(err)

	s.log.Debug("Set uploaded flag")
	err = s.postgres.SetUploaded(ctx, img.ID)
	s.Require().NoError(err)

	s.log.Debug("Get uploaded image")
	got, err := s.postgres.GetImage(ctx, img.ID)
	s.Require().NoError(err)
	s.Require().Equal(img.Name, got.Name)
	s.Require().True(got.Uploaded)

	s.log.Debug("Get unknown image")
	_, err = s.postgres.GetImage(ctx, uuid.New())
	s.Require().ErrorIs(err, ErrNotFound)

	s.log.Debug("Get not uploaded image")
	pending := GetTestImage()

	err = s.postgres.CreateImage(ctx, pending)
	s.Require().NoError(err)

	_, err = s.postgres.GetImage(ctx, pending.ID)
	s.Require().ErrorIs(err, ErrNotFound)
}

func (s *Suite) TestUpdateImage() {
	ctx := context.Background()

	img := GetTestImage()

	s.log.Debug("Create image")
	err := s.postgres.CreateImage(ctx, img)
	s.Require().NoError(err)

	s.log.Debug("Set uploaded flag")
	err = s.postgres.SetUploaded(ctx, img.ID)
	s.Require().NoError(err)

	s.log.Debug("Update image name")
	updated := &Image{ID: img.ID, Name: "new name"}

	err = s.postgres.UpdateImage(ctx, updated)
	s.Require().NoError(err)
	s.Require().Equal("new name", updated.Name)
	s.Require().Equal(img.ID, updated.ID)
	s.Require().True(updated.Uploaded)

	s.log.Debug("Check update in DB")
	got, err := s.postgres.GetImage(ctx, img.ID)
	s.Require().NoError(err)
	s.Require().Equal("new name", got.Name)

	s.log.Debug("Update unknown image")
	err = s.postgres.UpdateImage(ctx, &Image{ID: uuid.New(), Name: "new name"})
	s.Require().ErrorIs(err, ErrNotFound)

	s.log.Debug("Create not uploaded image")
	pending := GetTestImage()

	err = s.postgres.CreateImage(ctx, pending)
	s.Require().NoError(err)

	s.log.Debug("Update not uploaded image")
	err = s.postgres.UpdateImage(ctx, &Image{ID: pending.ID, Name: "new name"})
	s.Require().ErrorIs(err, ErrNotFound)
}

func (s *Suite) TestListImages() {
	ctx := context.Background()

	s.log.Debug("Create not uploaded image")
	pending := GetTestImage()

	err := s.postgres.CreateImage(ctx, pending)
	s.Require().NoError(err)

	s.log.Debug("Create and upload first image")
	first := GetTestImage()

	err = s.postgres.CreateImage(ctx, first)
	s.Require().NoError(err)

	err = s.postgres.SetUploaded(ctx, first.ID)
	s.Require().NoError(err)

	s.log.Debug("Create and upload second image")
	second := GetTestImage()

	err = s.postgres.CreateImage(ctx, second)
	s.Require().NoError(err)

	err = s.postgres.SetUploaded(ctx, second.ID)
	s.Require().NoError(err)

	s.log.Debug("Create and upload third image")
	third := GetTestImage()

	err = s.postgres.CreateImage(ctx, third)
	s.Require().NoError(err)

	err = s.postgres.SetUploaded(ctx, third.ID)
	s.Require().NoError(err)

	s.log.Debug("Total images counts uploaded images only")
	total, err := s.postgres.TotalImages(ctx)
	s.Require().NoError(err)
	s.Require().Equal(3, total)

	s.log.Debug("List all uploaded images")
	list, err := s.postgres.ListImages(ctx, 100, nil)
	s.Require().NoError(err)

	s.Require().Len(list, 3)
	s.Require().Equal(list[2].ID, first.ID)
	s.Require().Equal(list[1].ID, second.ID)
	s.Require().Equal(list[0].ID, third.ID)

	s.log.Debug("List uploaded images with limit")
	list, err = s.postgres.ListImages(ctx, 2, nil)
	s.Require().NoError(err)
	s.Require().Len(list, 2)

	s.log.Debug("List with start")
	list, err = s.postgres.ListImages(ctx, 2, &second.ID)
	s.Require().NoError(err)
	s.Require().Len(list, 2)
	s.Require().Equal(second.ID, list[0].ID)
}

func (s *Suite) TestDeleteImage() {
	ctx := context.Background()

	s.log.Debug("Create and upload image to delete")
	deleted := GetTestImage()

	err := s.postgres.CreateImage(ctx, deleted)
	s.Require().NoError(err)

	err = s.postgres.SetUploaded(ctx, deleted.ID)
	s.Require().NoError(err)

	s.log.Debug("Create and upload image to keep")
	kept := GetTestImage()

	err = s.postgres.CreateImage(ctx, kept)
	s.Require().NoError(err)

	err = s.postgres.SetUploaded(ctx, kept.ID)
	s.Require().NoError(err)

	s.log.Debug("Delete image")
	err = s.postgres.DeleteImage(ctx, deleted.ID)
	s.Require().NoError(err)

	s.log.Debug("Check image is deleted")
	_, err = s.postgres.GetImage(ctx, deleted.ID)
	s.Require().ErrorIs(err, ErrNotFound)

	s.log.Debug("Delete image twice")
	err = s.postgres.DeleteImage(ctx, deleted.ID)
	s.Require().ErrorIs(err, ErrNotFound)

	s.log.Debug("Create not uploaded image")
	pending := GetTestImage()

	err = s.postgres.CreateImage(ctx, pending)
	s.Require().NoError(err)

	s.log.Debug("Delete not uploaded image")
	err = s.postgres.DeleteImage(ctx, pending.ID)
	s.Require().ErrorIs(err, ErrNotFound)
}

func (s *Suite) TestDeleteImageForce() {
	ctx := context.Background()

	s.log.Debug("Create and upload image to delete")
	deleted := GetTestImage()

	err := s.postgres.CreateImage(ctx, deleted)
	s.Require().NoError(err)

	err = s.postgres.SetUploaded(ctx, deleted.ID)
	s.Require().NoError(err)

	s.log.Debug("Create and upload image to keep")
	kept := GetTestImage()

	err = s.postgres.CreateImage(ctx, kept)
	s.Require().NoError(err)

	err = s.postgres.SetUploaded(ctx, kept.ID)
	s.Require().NoError(err)

	s.log.Debug("Delete image")
	err = s.postgres.DeleteImageForce(ctx, deleted.ID)
	s.Require().NoError(err)

	s.log.Debug("Check image is deleted")
	_, err = s.postgres.GetImage(ctx, deleted.ID)
	s.Require().ErrorIs(err, ErrNotFound)

	s.log.Debug("Delete image twice")
	err = s.postgres.DeleteImageForce(ctx, deleted.ID)
	s.Require().NoError(err)

	s.log.Debug("Create not uploaded image")
	pending := GetTestImage()

	err = s.postgres.CreateImage(ctx, pending)
	s.Require().NoError(err)

	s.log.Debug("Delete not uploaded image")
	err = s.postgres.DeleteImageForce(ctx, pending.ID)
	s.Require().NoError(err)
}
