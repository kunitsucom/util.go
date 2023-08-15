package stringz

func MaskPrefix(s, mask string, unmaskSuffix int) (masked string) {
	for i, r := range s {
		if unmaskSuffix <= i {
			masked += string(r)
			continue
		}
		masked += mask
	}
	return masked
}

func MaskSuffix(s, mask string, unmaskPrefix int) (masked string) {
	for i, r := range s {
		if unmaskPrefix > i {
			masked += string(r)
			continue
		}
		masked += mask
	}
	return masked
}
