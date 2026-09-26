package generator

import "github.com/fogleman/gg"

const (
	pixelBlockDiv = 32
	color8Divisor = 257
	color8Max     = 255
	blurRadius    = 4
)

func pixelate(dc *gg.Context) {
	img := dc.Image()
	width, height := img.Bounds().Dx(), img.Bounds().Dy()
	block := max(1, min(width, height)/pixelBlockDiv)

	for blockY := 0; blockY < height; blockY += block {
		for blockX := 0; blockX < width; blockX += block {
			var sumR, sumG, sumB, count float64
			for y := blockY; y < blockY+block && y < height; y++ {
				for x := blockX; x < blockX+block && x < width; x++ {
					r, g, b, _ := img.At(x, y).RGBA()
					sumR += float64(r) / color8Divisor
					sumG += float64(g) / color8Divisor
					sumB += float64(b) / color8Divisor
					count++
				}
			}

			dc.SetRGBA(sumR/count/color8Max, sumG/count/color8Max, sumB/count/color8Max, 1)
			dc.DrawRectangle(float64(blockX), float64(blockY), float64(block), float64(block))
			dc.Fill()
		}
	}
}

func blur(dc *gg.Context) {
	img := dc.Image()
	width, height := img.Bounds().Dx(), img.Bounds().Dy()
	radius := blurRadius

	src := make([]float64, width*height*3)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			idx := (y*width + x) * 3
			src[idx] = float64(r) / color8Divisor
			src[idx+1] = float64(g) / color8Divisor
			src[idx+2] = float64(b) / color8Divisor
		}
	}

	hBuf := make([]float64, width*height*3)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			var sum [3]float64
			for offset := -radius; offset <= radius; offset++ {
				sampleX := clamp(x+offset, 0, width-1)
				i := (y*width + sampleX) * 3
				sum[0] += src[i]
				sum[1] += src[i+1]
				sum[2] += src[i+2]
			}
			idx := (y*width + x) * 3
			for channel := 0; channel < 3; channel++ {
				hBuf[idx+channel] = sum[channel] / float64(2*radius+1)
			}
		}
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			var sum [3]float64
			for offset := -radius; offset <= radius; offset++ {
				sampleY := clamp(y+offset, 0, height-1)
				idx := (sampleY*width + x) * 3
				sum[0] += hBuf[idx]
				sum[1] += hBuf[idx+1]
				sum[2] += hBuf[idx+2]
			}
			weight := float64(2*radius + 1)
			setPix(dc, x, y, sum[0]/weight/color8Max, sum[1]/weight/color8Max, sum[2]/weight/color8Max, 1)
		}
	}
}

func clamp[T ~int | ~float64](value, low, high T) T {
	if value < low {
		return low
	}
	if value > high {
		return high
	}

	return value
}
