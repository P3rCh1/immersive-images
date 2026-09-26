package generator

import "math"

type rgb struct {
	r, g, b float64
}

type stop struct {
	pos   float64
	color rgb
}

type ramp []stop

var palettes = map[string]ramp{
	PaletteNeon: {
		{0.00, rgb{0.02, 0.01, 0.08}},
		{0.35, rgb{0.06, 0.02, 0.20}},
		{0.55, rgb{0.36, 0.12, 0.75}},
		{0.70, rgb{0.00, 0.78, 0.95}},
		{0.85, rgb{1.00, 0.35, 0.90}},
		{1.00, rgb{0.95, 1.00, 0.85}},
	},
	PaletteDark: {
		{0.00, rgb{0.004, 0.020, 0.060}},
		{0.30, rgb{0.010, 0.100, 0.200}},
		{0.55, rgb{0.040, 0.300, 0.400}},
		{0.75, rgb{0.150, 0.750, 0.650}},
		{1.00, rgb{0.550, 1.000, 0.900}},
	},
	PaletteToxic: {
		{0.00, rgb{0.03, 0.07, 0.03}},
		{0.25, rgb{0.10, 0.30, 0.06}},
		{0.50, rgb{0.35, 0.85, 0.12}},
		{0.70, rgb{0.72, 1.00, 0.22}},
		{0.85, rgb{0.95, 1.00, 0.42}},
		{1.00, rgb{0.97, 1.00, 0.86}},
	},
	PaletteRainbow: {
		{0.00, rgb{1.00, 0.42, 0.49}},
		{0.17, rgb{1.00, 0.63, 0.36}},
		{0.33, rgb{1.00, 0.85, 0.23}},
		{0.50, rgb{0.61, 0.88, 0.36}},
		{0.67, rgb{0.50, 0.85, 0.91}},
		{0.83, rgb{0.54, 0.56, 0.85}},
		{1.00, rgb{0.78, 0.61, 0.85}},
	},
	PaletteOcean: {
		{0.00, rgb{0.00, 0.08, 0.23}},
		{0.30, rgb{0.04, 0.30, 0.55}},
		{0.55, rgb{0.09, 0.65, 0.72}},
		{0.75, rgb{0.43, 0.91, 0.78}},
		{1.00, rgb{0.90, 1.00, 0.96}},
	},
	PaletteLava: {
		{0.00, rgb{0.04, 0.02, 0.02}},
		{0.40, rgb{0.48, 0.12, 0.06}},
		{0.60, rgb{0.88, 0.35, 0.17}},
		{0.80, rgb{0.96, 0.72, 0.24}},
		{1.00, rgb{1.00, 0.96, 0.85}},
	},
	PaletteMono: {
		{0.00, rgb{0.02, 0.02, 0.02}},
		{0.50, rgb{0.50, 0.50, 0.50}},
		{1.00, rgb{0.98, 0.98, 0.98}},
	},
}

func hueShift(source ramp, turn float64) ramp {
	if turn == 0 {
		return source
	}

	out := make(ramp, len(source))
	copy(out, source)

	for i := range out {
		h, s, v := rgb2hsv(out[i].color.r, out[i].color.g, out[i].color.b)
		h = math.Mod(h+turn, 1)
		if h < 0 {
			h++
		}
		out[i].color.r, out[i].color.g, out[i].color.b = hsv2rgb(h, s, v)
	}

	return out
}

func colorAt(source ramp, value float64) (float64, float64, float64) {
	value = clamp(value, 0.0, 1.0)

	if len(source) == 0 {
		return 0, 0, 0
	}

	if value <= source[0].pos {
		return source[0].color.r, source[0].color.g, source[0].color.b
	}

	for i := 0; i < len(source)-1; i++ {
		from, to := source[i], source[i+1]
		if value <= to.pos {
			k := (value - from.pos) / (to.pos - from.pos)

			return lerp(k, from.color, to.color)
		}
	}

	last := source[len(source)-1]

	return last.color.r, last.color.g, last.color.b
}

func lerp(k float64, from, to rgb) (float64, float64, float64) {
	return from.r + k*(to.r-from.r), from.g + k*(to.g-from.g), from.b + k*(to.b-from.b)
}
