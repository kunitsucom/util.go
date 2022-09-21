package unit

const (
	Ki = 1024
	Mi = Ki * 1024
	Gi = Mi * 1024
	Ti = Gi * 1024
	Pi = Ti * 1024
	Ei = Pi * 1024
	Zi = Ei * 1024
	Yi = Zi * 1024
)

func ToKi(n uint64) uint64 {
	return n / Ki
}

func ToMi(n uint64) uint64 {
	return n / Mi
}

func ToGi(n uint64) uint64 {
	return n / Gi
}

func ToTi(n uint64) uint64 {
	return n / Ti
}

func ToPi(n uint64) uint64 {
	return n / Pi
}

func ToEi(n uint64) uint64 {
	return n / Ei
}
