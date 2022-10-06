package unit

const (
	Ki = 1 * 1024
	Mi = Ki * 1024
	Gi = Mi * 1024
	Ti = Gi * 1024
	Pi = Ti * 1024
	Ei = Pi * 1024
	Zi = Ei * 1024
	Yi = Zi * 1024
)

func ToKi(n uint64) (ki uint64) {
	return n / Ki
}

func ToMi(n uint64) (mi uint64) {
	return n / Mi
}

func ToGi(n uint64) (gi uint64) {
	return n / Gi
}

func ToTi(n uint64) (ti uint64) {
	return n / Ti
}

func ToPi(n uint64) (pi uint64) {
	return n / Pi
}

func ToEi(n uint64) (ei uint64) {
	return n / Ei
}

func KiTo(ki uint64) (n uint64) {
	return ki * Ki
}

func MiTo(mi uint64) (n uint64) {
	return mi * Mi
}

func GiTo(gi uint64) (n uint64) {
	return gi * Gi
}

func TiTo(ti uint64) (n uint64) {
	return ti * Ti
}

func PiTo(pi uint64) (n uint64) {
	return pi * Pi
}

func EiTo(ei uint64) (n uint64) {
	return ei * Ei
}
