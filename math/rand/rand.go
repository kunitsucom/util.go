package randz

import (
	"math/rand"
	"time"
)

func Range(minValue, maxValue int) int {
	//nolint:gosec
	return rand.Intn(maxValue-minValue+1) + minValue
}

func RangeRand(r *rand.Rand, minValue, maxValue int) int {
	return r.Intn(maxValue-minValue+1) + minValue
}

func Range31(minValue, maxValue int32) int32 {
	//nolint:gosec
	return rand.Int31n(maxValue-minValue+1) + minValue
}

func Range31Rand(r *rand.Rand, minValue, maxValue int32) int32 {
	return r.Int31n(maxValue-minValue+1) + minValue
}

func Range63(minValue, maxValue int64) int64 {
	//nolint:gosec
	return rand.Int63n(maxValue-minValue+1) + minValue
}

func Range63Rand(r *rand.Rand, minValue, maxValue int64) int64 {
	return r.Int63n(maxValue-minValue+1) + minValue
}

func RangeDuration(minValue, maxValue time.Duration) time.Duration {
	//nolint:gosec
	return time.Duration(rand.Int63n(int64(maxValue)-int64(minValue)+1) + int64(minValue))
}

func RangeDurationRand(r *rand.Rand, minValue, maxValue time.Duration) time.Duration {
	//nolint:gosec
	return time.Duration(r.Int63n(int64(maxValue)-int64(minValue)+1) + int64(minValue))
}
