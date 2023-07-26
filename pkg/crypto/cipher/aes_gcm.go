package cipherz

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

const (
	recommendedGCMNonceSize = 12 // AES-GCM の初期化ベクトル IV は 12 バイト (96ビット) を推奨 ref. https://www.cryptrec.go.jp/list/cryptrec-ls-0001-2012r5.pdf ref. https://go.dev/src/crypto/cipher/gcm.go#L157
)

// NewAESGCM returns a Cipher using AES-GCM.
// To select AES-128-GCM, AES-192-GCM, or AES-256-GCM, specify either 16, 24, or 32 bytes for the key length.
func NewAESGCM(key []byte) (Cipher, error) { //nolint:ireturn
	c, err := newAESGCM(key, recommendedGCMNonceSize)
	if err != nil {
		return nil, fmt.Errorf("newAESGCM: %w", err)
	}

	return c, nil
}

type _GCM struct {
	key  []byte
	aead cipher.AEAD
	rand io.Reader
}

func newAESGCM(key []byte, nonceSize int) (*_GCM, error) { //nolint:ireturn
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("aes.NewCipher: %w", err)
	}

	gcm, err := cipher.NewGCMWithNonceSize(block, nonceSize)
	if err != nil {
		return nil, fmt.Errorf("cipher.NewGCMWithNonceSize: %w", err)
	}

	return &_GCM{
		key:  key,
		aead: gcm,
		rand: rand.Reader,
	}, nil
}

func (gcm *_GCM) Encrypt(plainBytes []byte) (cipherBytes []byte, err error) {
	nonce := make([]byte, gcm.aead.NonceSize())

	if _, err := gcm.rand.Read(nonce); err != nil {
		return nil, fmt.Errorf("gcm.rand.Read: %w", err)
	}

	cipherBytes = gcm.aead.Seal(nil, nonce, plainBytes, nil)

	//nolint: makezero
	cipherBytes = append(nonce, cipherBytes...)

	return cipherBytes, nil
}

func (gcm *_GCM) Decrypt(cipherBytes []byte) (plainBytes []byte, err error) {
	if len(cipherBytes) < gcm.aead.NonceSize() {
		return nil, fmt.Errorf("len(cipherBytes)=%d: %w", len(cipherBytes), ErrCipherBytesIsTooShort)
	}

	nonce := cipherBytes[:gcm.aead.NonceSize()]
	cipherBytes = cipherBytes[gcm.aead.NonceSize():]

	plainBytes, err = gcm.aead.Open(nil, nonce, cipherBytes, nil)
	if err != nil {
		return nil, fmt.Errorf("gcm.aead.Open: %w", err)
	}

	return plainBytes, nil
}
