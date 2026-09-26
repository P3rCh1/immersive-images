package generator

import (
	"bytes"
	"errors"
	"fmt"
	"image/png"
	"math"

	"github.com/fogleman/gg"
)

var (
	ErrUnknownStyle      = errors.New("generator: unknown style")
	ErrUnknownPalette    = errors.New("generator: unknown palette")
	ErrUnknownAdditional = errors.New("generator: unknown additional")
	ErrInvalidDimensions = errors.New("generator: width and height must be positive")
	ErrInvalidScale      = errors.New("generator: scale must be positive")
)

const (
	StylePatches  = "PATCHES"
	StyleSpots    = "SPOTS"
	StyleRings    = "RINGS"
	StyleTerraces = "TERRACES"
	StyleWaves    = "WAVES"
	StyleMottled  = "MOTTLED"
)

const (
	PaletteNeon    = "NEON"
	PaletteDark    = "DARK"
	PaletteToxic   = "TOXIC"
	PaletteRainbow = "RAINBOW"
	PaletteOcean   = "OCEAN"
	PaletteLava    = "LAVA"
	PaletteMono    = "MONO"
)

const (
	AdditionalPixel = "PIXEL"
	AdditionalBlur  = "BLUR"
)

type Config struct {
	Seed       int
	Style      string
	Palette    string
	Additional *string
	Width      int
	Height     int
	Scale      float64
}

var renderers = map[string]func(*scene) *gg.Context{
	StylePatches:  renderPatches,
	StyleSpots:    renderSpots,
	StyleRings:    renderRings,
	StyleTerraces: renderTerraces,
	StyleWaves:    renderWaves,
	StyleMottled:  renderMottled,
}

var additions = map[string]func(*gg.Context){
	AdditionalPixel: pixelate,
	AdditionalBlur:  blur,
}

func ValidStyle(style string) bool {
	_, ok := renderers[style]

	return ok
}

func ValidPalette(palette string) bool {
	_, ok := palettes[palette]

	return ok
}

func ValidAdditional(additional string) bool {
	_, ok := additions[additional]
	return ok
}

func normalizeSeed(seed int) int {
	const seedModulo = 1 << 32

	normalized := int64(seed) % seedModulo
	if normalized < 0 {
		normalized += seedModulo
	}

	return int(normalized)
}

func GenerateImage(cfg Config) ([]byte, error) {
	cfg.Seed = normalizeSeed(cfg.Seed)

	renderer, ok := renderers[cfg.Style]
	if !ok {
		return nil, ErrUnknownStyle
	}
	if _, ok := palettes[cfg.Palette]; !ok {
		return nil, ErrUnknownPalette
	}
	if cfg.Width <= 0 || cfg.Height <= 0 {
		return nil, ErrInvalidDimensions
	}
	if cfg.Scale <= 0 {
		return nil, ErrInvalidScale
	}
	if cfg.Additional != nil {
		if _, ok := additions[*cfg.Additional]; !ok {
			return nil, ErrUnknownAdditional
		}
	}

	s := newScene(cfg)
	dc := renderer(s)

	if cfg.Additional != nil {
		if apply, ok := additions[*cfg.Additional]; ok {
			apply(dc)
		}
	}

	return encodePNG(dc)
}

func encodePNG(dc *gg.Context) ([]byte, error) {
	var buf bytes.Buffer
	if err := png.Encode(&buf, dc.Image()); err != nil {
		return nil, fmt.Errorf("failed to encode generated image: %w", err)
	}

	return buf.Bytes(), nil
}

func setPix(dc *gg.Context, x, y int, r, g, b, a float64) {
	dc.SetRGBA(r, g, b, a)
	dc.SetPixel(x, y)
}

func smoothstep(e0, e1, x float64) float64 {
	t := clamp((x-e0)/(e1-e0), 0.0, 1.0)
	return t * t * (3 - 2*t)
}

func hsv2rgb(h, s, v float64) (float64, float64, float64) {
	i := int(h * 6)
	f := h*6 - float64(i)
	p := v * (1 - s)
	q := v * (1 - f*s)
	t := v * (1 - (1-f)*s)

	switch i % 6 {
	case 0:
		return v, t, p
	case 1:
		return q, v, p
	case 2:
		return p, v, t
	case 3:
		return p, q, v
	case 4:
		return t, p, v
	default:
		return v, p, q
	}
}

func rgb2hsv(r, g, b float64) (float64, float64, float64) {
	mx := math.Max(r, math.Max(g, b))
	mn := math.Min(r, math.Min(g, b))
	d := mx - mn

	h := 0.0
	switch {
	case d == 0:
		h = 0
	case mx == r:
		h = math.Mod((g-b)/d, 6)
	case mx == g:
		h = (b-r)/d + 2
	default:
		h = (r-g)/d + 4
	}
	h /= 6
	if h < 0 {
		h++
	}

	s := 0.0
	if mx != 0 {
		s = d / mx
	}

	return h, s, mx
}
