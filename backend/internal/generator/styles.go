package generator

import (
	"math"

	"github.com/fogleman/gg"
)

func newCanvas(s *scene) (*gg.Context, int, int) {
	width, height := s.cfg.Width, s.cfg.Height

	return gg.NewContext(width, height), width, height
}

func paint(s *scene, dc *gg.Context, x, y int, value float64) {
	r, g, b := colorAt(s.ramp, value)
	setPix(dc, x, y, r, g, b, 1)
}

func renderPatches(s *scene) *gg.Context {
	dc, width, height := newCanvas(s)
	cellSize := max(32, min(width, height)/8)
	seed := s.cfg.Seed

	jitterX := func(i, j int) float64 { return (hashUnit(int64(seed), i, j) - 0.5) * float64(cellSize) }
	jitterY := func(i, j int) float64 { return (hashUnit(int64(seed), i+1013, j+5197) - 0.5) * float64(cellSize) }
	cellColor := func(i, j int) float64 { return hashUnit(int64(seed)*31+7, i, j) }

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			cellX, cellY := x/cellSize, y/cellSize
			bestDist := 1.0e30
			bestI, bestJ := cellX, cellY

			for di := -1; di <= 1; di++ {
				for dj := -1; dj <= 1; dj++ {
					ni, nj := cellX+di, cellY+dj
					centerX := float64(ni*cellSize+cellSize/2) + jitterX(ni, nj)
					centerY := float64(nj*cellSize+cellSize/2) + jitterY(ni, nj)
					offX := float64(x) - centerX
					offY := float64(y) - centerY
					if d := offX*offX + offY*offY; d < bestDist {
						bestDist = d
						bestI, bestJ = ni, nj
					}
				}
			}

			paint(s, dc, x, y, cellColor(bestI, bestJ))
		}
	}

	return dc
}

func renderSpots(s *scene) *gg.Context {
	dc, width, height := newCanvas(s)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			field := (s.noise(x, y, 1) + 1) / 2
			fieldFine := (s.noise(x, y, 1.8) + 1) / 2
			value := smoothstep(0.52, 0.76, field+(fieldFine-0.5)*0.45)

			paint(s, dc, x, y, value)
		}
	}

	return dc
}

func renderRings(s *scene) *gg.Context {
	dc, width, height := newCanvas(s)
	minDim := float64(min(width, height))
	random := newRng(s.cfg.Seed + 7919)
	count := 2 + random.IntN(5)
	spread := 0.9 + random.Float()*0.95
	spacing := spread / float64(count)
	bandWidth := spacing * (0.26 + random.Float()*0.2)

	centerX := float64(width)/2 + (random.Float()-0.5)*minDim*0.18
	centerY := float64(height)/2 + (random.Float()-0.5)*minDim*0.18

	shimmer := make([]float64, count)
	for i := 0; i < count; i++ {
		shimmer[i] = 0.7 + random.Float()*0.3
	}

	rotation := random.Float() * 2 * math.Pi
	cosA, sinA := math.Cos(rotation), math.Sin(rotation)
	axisX := 0.85 + 0.35*random.Float()
	axisY := 0.85 + 0.35*random.Float()
	rippleAmp := random.Float() * 0.05
	rippleFreq := float64(2 + random.IntN(4))
	ripplePhase := random.Float() * 2 * math.Pi

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			offX := float64(x) - centerX
			offY := float64(y) - centerY
			rx := offX*cosA - offY*sinA
			ry := offX*sinA + offY*cosA
			angle := math.Atan2(ry, rx)

			radius := math.Hypot(rx/axisX, ry/axisY)
			ringDist := radius * (1 + rippleAmp*math.Sin(rippleFreq*angle+ripplePhase)) * 2 / minDim
			wobble := s.noise(x, y, 1)*0.02 + 0.012*math.Sin(float64(y)*0.004)

			value := 0.0
			for i := 0; i < count; i++ {
				position := spacing/2 + spacing*float64(i) + wobble*0.5
				delta := math.Abs(ringDist - position)
				brightness := math.Pow(math.Max(0, 1-delta/bandWidth), 1.6) * shimmer[i]
				if brightness > value {
					value = brightness
				}
			}

			value = clamp(value+(1-clamp((ringDist-0.92)/0.2, 0.0, 1.0))*0.08, 0.0, 1.0)
			paint(s, dc, x, y, value)
		}
	}

	return dc
}

func renderTerraces(s *scene) *gg.Context {
	dc, width, height := newCanvas(s)
	random := newRng(s.cfg.Seed + 127)
	levels := float64(4 + random.IntN(5))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			value := float64(int((s.noise(x, y, 1)+1)/2*levels)) / levels
			paint(s, dc, x, y, value)
		}
	}

	return dc
}

func renderWaves(s *scene) *gg.Context {
	dc, width, height := newCanvas(s)
	random := newRng(s.cfg.Seed + 7)
	angle := random.Float()*0.5 - 0.25
	cosA, sinA := math.Cos(angle), math.Sin(angle)
	bandCount := float64(3 + random.IntN(4))
	span := float64(height)*cosA + float64(width)*math.Abs(sinA)
	period := span / bandCount

	bendAmp := random.Float()*0.5 - 0.25
	bendFreq := float64(1 + random.IntN(3))
	bendPhase := random.Float() * 2 * math.Pi

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			along := (float64(x)*cosA - float64(y)*sinA) / span
			across := (float64(x)*sinA + float64(y)*cosA) / period
			across += bendAmp * math.Sin(2*math.Pi*bendFreq*along+bendPhase)

			layer := math.Floor(across)
			fraction := across - layer
			base := 0.12
			if math.Mod(layer, 2) != 0 {
				base = 0.6
			}
			value := clamp(base+smoothstep(0.12, 0.4, fraction)*0.25, 0.0, 1.0)

			paint(s, dc, x, y, value)
		}
	}

	return dc
}

func renderMottled(s *scene) *gg.Context {
	dc, width, height := newCanvas(s)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			background := (s.noise(x, y, 1) + 1) / 3
			coarse := smoothstep(0.66, 0.78, (s.noise(x, y, 2.6)+1)/2)
			fine := smoothstep(0.6, 0.7, (s.noise(x, y, 4)+1)/2)
			value := clamp(math.Max(background*0.85, coarse+fine*0.7), 0.0, 1.0)

			paint(s, dc, x, y, value)
		}
	}

	return dc
}
