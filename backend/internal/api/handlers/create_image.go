package handlers

import (
	"encoding/json/v2"
	"fmt"
	"net/http"
	"strings"

	"github.com/P3rCh1/immersive-images/backend/internal/api/models"
	"github.com/P3rCh1/immersive-images/backend/internal/api/msgs"
	"github.com/P3rCh1/immersive-images/backend/internal/config"
	"github.com/P3rCh1/immersive-images/backend/internal/dal/postgres"
	"github.com/P3rCh1/immersive-images/backend/internal/generator"
)

const maxImageNameLen = 255

func (h *Handlers) CreateImage(w http.ResponseWriter, r *http.Request) {
	var req models.CreateImageRequest
	if err := json.UnmarshalRead(r.Body, &req, json.DefaultOptionsV2()); err != nil {
		h.Error(w, http.StatusBadRequest, msgs.InvalidBody)
		return
	}

	cfg, err := getGeneratorConfig(req)
	if err != nil {
		h.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	data, err := generator.GenerateImage(cfg)
	if err != nil {
		h.Error(w, http.StatusInternalServerError, msgs.Internal)
		return
	}

	img := postgres.Image{
		Name:       req.Name,
		Size:       int64(len(data)),
		Seed:       cfg.Seed,
		Style:      cfg.Style,
		Palette:    cfg.Palette,
		Additional: cfg.Additional,
		Width:      cfg.Width,
		Height:     cfg.Height,
		Scale:      cfg.Scale,
	}
	if err := h.db.CreateImage(r.Context(), &img); err != nil {
		h.Error(w, http.StatusInternalServerError, msgs.Internal)
		return
	}

	if err := h.s3.Upload(r.Context(), img.ID.String(), data); err != nil {
		if delErr := h.db.DeleteImage(r.Context(), img.ID); delErr != nil {
			h.log.Error(
				"failed to delete image after S3 upload failure",
				"need_cleanup", "DB",
				"image_id", img.ID,
				"error", delErr,
			)
		}

		h.Error(w, http.StatusServiceUnavailable, msgs.Unavailable)
		return
	}

	if err := h.db.SetUploaded(r.Context(), img.ID); err != nil {
		h.log.Error(
			"failed to store image uploaded flag",
			"need_cleanup", "S3",
			"image_id", img.ID,
			"error", err,
		)

		h.Error(w, http.StatusInternalServerError, msgs.Internal)
		return
	}

	h.JSON(w, http.StatusCreated, apiImageMetaFromDB(img))
}

func getGeneratorConfig(req models.CreateImageRequest) (generator.Config, error) {
	if len(req.Name) < 1 || len(req.Name) > maxImageNameLen {
		return generator.Config{}, fmt.Errorf(
			msgs.InvalidNameLenFmt,
			len(req.Name), 1, maxImageNameLen,
		)
	}

	images := config.Config.Images

	cfg := generator.Config{
		Seed:    req.Seed,
		Style:   strings.ToUpper(images.DefaultStyle),
		Palette: strings.ToUpper(images.DefaultPalette),
		Width:   images.DefaultWidth,
		Height:  images.DefaultHeight,
		Scale:   images.DefaultScale,
	}

	if req.Settings == nil {
		return cfg, nil
	}

	if req.Settings.Style != nil {
		style := strings.ToUpper(string(*req.Settings.Style))
		if !generator.ValidStyle(style) {
			return generator.Config{}, msgs.ErrInvalidStyle
		}

		cfg.Style = style
	}

	if req.Settings.Palette != nil {
		palette := strings.ToUpper(string(*req.Settings.Palette))
		if !generator.ValidPalette(palette) {
			return generator.Config{}, msgs.ErrInvalidPalette
		}

		cfg.Palette = palette
	}

	if req.Settings.Width != nil {
		width := *req.Settings.Width
		if width < images.MinWidth || width > images.MaxWidth {
			return generator.Config{}, fmt.Errorf(
				msgs.InvalidWidthFmt,
				width, images.MinWidth, images.MaxWidth,
			)
		}
		cfg.Width = width
	}

	if req.Settings.Height != nil {
		height := *req.Settings.Height
		if height < images.MinHeight || height > images.MaxHeight {
			return generator.Config{}, fmt.Errorf(
				msgs.InvalidHeightFmt,
				height, images.MinHeight, images.MaxHeight,
			)
		}
		cfg.Height = height
	}

	if req.Settings.Scale != nil {
		scale := *req.Settings.Scale
		if scale < images.MinScale || scale > images.MaxScale {
			return generator.Config{}, fmt.Errorf(
				msgs.InvalidScaleFmt,
				scale, images.MinScale, images.MaxScale,
			)
		}
		cfg.Scale = scale
	}

	if req.Settings.Additional != nil {
		additional := strings.ToUpper(string(*req.Settings.Additional))
		if !generator.ValidAdditional(additional) {
			return generator.Config{}, msgs.ErrInvalidAdditional
		}

		cfg.Additional = &additional
	}

	return cfg, nil
}
