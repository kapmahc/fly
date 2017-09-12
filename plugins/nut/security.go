package nut

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha512"
)

var (
	_aes  *Aes
	_hmac *Hmac
)

// AES aes
func AES() *Aes {
	return _aes
}

// HMAC hmac
func HMAC() *Hmac {
	return _hmac
}

// Aes aes
type Aes struct {
	cip cipher.Block
}

// Encrypt aes encrypt
func (p *Aes) Encrypt(buf []byte) ([]byte, error) {
	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return nil, err
	}
	cfb := cipher.NewCFBEncrypter(p.cip, iv)
	val := make([]byte, len(buf))
	cfb.XORKeyStream(val, buf)

	return append(val, iv...), nil
}

// Decrypt aes decrypt
func (p *Aes) Decrypt(buf []byte) ([]byte, error) {
	bln := len(buf)
	cln := bln - aes.BlockSize
	ct := buf[0:cln]
	iv := buf[cln:bln]

	cfb := cipher.NewCFBDecrypter(p.cip, iv)
	val := make([]byte, cln)
	cfb.XORKeyStream(val, ct)
	return val, nil
}

// Hmac hmac
type Hmac struct {
	key []byte
}

// Sum sum hmac
func (p *Hmac) Sum(plain []byte) []byte {
	mac := hmac.New(sha512.New, p.key)
	return mac.Sum(plain)
}

// Chk chk hmac
func (p *Hmac) Chk(plain, code []byte) bool {
	return hmac.Equal(p.Sum(plain), code)
}
