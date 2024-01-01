package errorz

import "errors"

func PanicOrIgnore(err error, ignores ...error) {
	if err == nil {
		return
	}

	for _, ignore := range ignores {
		if errors.Is(err, ignore) {
			return
		}
	}

	panic(err)
}
