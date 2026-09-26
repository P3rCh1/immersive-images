package generator

import (
	"math"

	"github.com/aquilax/go-perlin"
)

type scene struct {
	cfg        Config
	noiseGen   *perlin.Perlin
	ramp       ramp
	rotCos     float64
	rotSin     float64
	offX, offY float64
}

const noiseOffsetRange = 1000

func newScene(cfg Config) *scene {
	s := &scene{cfg: cfg}

	rng := newRng(cfg.Seed)
	angle := rng.Float() * 2 * math.Pi
	s.rotCos, s.rotSin = math.Cos(angle), math.Sin(angle)
	s.offX, s.offY = rng.Float()*noiseOffsetRange, rng.Float()*noiseOffsetRange

	s.noiseGen = perlin.NewPerlin(2.0, 0.5, 4, int64(cfg.Seed)) //nolint:gosec // integer overflow is not bad case
	s.ramp = hueShift(palettes[cfg.Palette], (rng.Float()-0.5)*0.12)

	return s
}

func hashUnit(seed int64, i, j int) float64 {
	h := uint64(seed) ^ 0x9E3779B97F4A7C15                           //nolint:gosec // integer overflow is not bad case
	h ^= uint64(i)*0xBF58476D1CE4E5B9 + uint64(j)*0x94D049BB133111EB //nolint:gosec // integer overflow is not bad case
	h ^= h >> 30
	h *= 0xBF58476D1CE4E5B9
	h ^= h >> 27
	h *= 0x94D049BB133111EB
	h ^= h >> 31

	return float64(h>>11) * (1.0 / float64(uint64(1)<<53))
}

func (s *scene) noise(x, y int, freq float64) float64 {
	cfg := s.cfg
	halfMin := math.Max(1, float64(min(cfg.Width, cfg.Height))/2)
	px := (float64(x) - float64(cfg.Width)/2) / halfMin
	py := (float64(y) - float64(cfg.Height)/2) / halfMin

	rx := px*s.rotCos - py*s.rotSin
	ry := px*s.rotSin + py*s.rotCos

	nx := rx + 0.5 + s.offX/noiseOffsetRange
	ny := ry + 0.5 + s.offY/noiseOffsetRange

	return s.noiseGen.Noise2D(nx*cfg.Scale*freq, ny*cfg.Scale*freq)
}
