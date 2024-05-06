package jws

import "strings"

func Parse(jwt string) (headerEncoded, payloadEncoded, signatureEncoded string, err error) {
	parts := strings.Split(jwt, ".")
	const expectedPartsLen = 3
	if len(parts) != expectedPartsLen {
		return "", "", "", ErrInvalidTokenReceived
	}

	return parts[0], parts[1], parts[2], nil
}
