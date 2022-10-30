package randz

import (
	"math/rand"
	"time"
)

func Range(min, max int) int {
	// nolint: gosec
	return rand.Intn(max-min+1) + min
}

func RangeRand(r *rand.Rand, min, max int) int {
	return r.Intn(max-min+1) + min
}

func Range31(min, max int32) int32 {
	// nolint: gosec
	return rand.Int31n(max-min+1) + min
}

func Range31Rand(r *rand.Rand, min, max int32) int32 {
	return r.Int31n(max-min+1) + min
}

func Range63(min, max int64) int64 {
	// nolint: gosec
	return rand.Int63n(max-min+1) + min
}

func Range63Rand(r *rand.Rand, min, max int64) int64 {
	return r.Int63n(max-min+1) + min
}

func RangeDuration(min, max time.Duration) time.Duration {
	// nolint: gosec
	return time.Duration(rand.Int63n(int64(max)-int64(min)+1) + int64(min))
}

func RangeDurationRand(r *rand.Rand, min, max time.Duration) time.Duration {
	// nolint: gosec
	return time.Duration(r.Int63n(int64(max)-int64(min)+1) + int64(min))
}
