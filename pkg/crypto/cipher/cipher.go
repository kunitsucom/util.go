package cipherz

type Cipher interface {
	Encrypt(plainBytes []byte) (cipherBytes []byte, err error)
	Decrypt(cipherBytes []byte) (plainBytes []byte, err error)
}
