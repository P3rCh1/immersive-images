package generator

import (
	"testing"

	"github.com/P3rCh1/immersive-images/backend/internal/common"
	"github.com/stretchr/testify/assert"
)

func TestValidations(t *testing.T) {
	assert.True(t, ValidStyle(StylePatches))
	assert.False(t, ValidStyle("UNKNOWN"))

	assert.True(t, ValidPalette(PaletteNeon))
	assert.False(t, ValidPalette("UNKNOWN"))

	assert.False(t, ValidAdditional("UNKNOWN"))
	assert.True(t, ValidAdditional(AdditionalPixel))
	assert.True(t, ValidAdditional(AdditionalBlur))
}

func TestGenerateImageErrors(t *testing.T) {
	_, err := GenerateImage(Config{
		Seed:    1,
		Style:   "UNKNOWN",
		Palette: PaletteNeon,
		Width:   10,
		Height:  10,
		Scale:   1,
	})
	assert.ErrorIs(t, err, ErrUnknownStyle)

	_, err = GenerateImage(Config{
		Seed:    1,
		Style:   StylePatches,
		Palette: "UNKNOWN",
		Width:   10,
		Height:  10,
		Scale:   1,
	})
	assert.ErrorIs(t, err, ErrUnknownPalette)

	_, err = GenerateImage(Config{
		Seed:    1,
		Style:   StylePatches,
		Palette: PaletteNeon,
		Width:   0,
		Height:  10,
		Scale:   1,
	})
	assert.ErrorIs(t, err, ErrInvalidDimensions)

	_, err = GenerateImage(Config{
		Seed:    1,
		Style:   StylePatches,
		Palette: PaletteNeon,
		Width:   10,
		Height:  0,
		Scale:   1,
	})
	assert.ErrorIs(t, err, ErrInvalidDimensions)

	_, err = GenerateImage(Config{
		Seed:    1,
		Style:   StylePatches,
		Palette: PaletteNeon,
		Width:   10,
		Height:  10,
		Scale:   0,
	})
	assert.ErrorIs(t, err, ErrInvalidScale)

	_, err = GenerateImage(Config{
		Seed:    1,
		Style:   StylePatches,
		Palette: PaletteNeon,
		Width:   10,
		Height:  10,
		Scale:   -1,
	})
	assert.ErrorIs(t, err, ErrInvalidScale)

	_, err = GenerateImage(Config{
		Seed:       1,
		Style:      StylePatches,
		Palette:    PaletteNeon,
		Additional: new("UNKNOWN"),
		Width:      10,
		Height:     10,
		Scale:      1,
	})
	assert.ErrorIs(t, err, ErrUnknownAdditional)
}

func TestGenerateImageSmoke(t *testing.T) {
	data, err := GenerateImage(Config{
		Seed:    42,
		Style:   StylePatches,
		Palette: PaletteNeon,
		Width:   32,
		Height:  32,
		Scale:   1.0,
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, data)

	data, err = GenerateImage(Config{
		Seed:       123,
		Style:      StyleSpots,
		Palette:    PaletteDark,
		Additional: common.Ptr(AdditionalPixel),
		Width:      16,
		Height:     16,
		Scale:      0.5,
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, data)

	data, err = GenerateImage(Config{
		Seed:       7,
		Style:      StyleRings,
		Palette:    PaletteToxic,
		Additional: common.Ptr(AdditionalBlur),
		Width:      24,
		Height:     24,
		Scale:      0.8,
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, data)
}
