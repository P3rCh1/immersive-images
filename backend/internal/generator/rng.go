package generator

type rng struct {
	state uint64
}

func newRng(seed int) *rng {
	return &rng{state: uint64(seed)<<32 ^ uint64(seed)} //nolint:gosec // integer overflow is not bad case
}

func (r *rng) Next() uint64 {
	r.state += 0x9E3779B97F4A7C15
	x := r.state
	x = (x ^ (x >> 30)) * 0xBF58476D1CE4E5B9
	x = (x ^ (x >> 27)) * 0x94D049BB133111EB

	return x ^ (x >> 31)
}

func (r *rng) Float() float64 {
	return float64(r.Next()>>11) * (1.0 / float64(uint64(1)<<53))
}

func (r *rng) IntN(n int) int {
	if n <= 0 {
		return 0
	}

	return int(r.Next() % uint64(n)) //nolint:gosec // integer overflow is not bad case
}
