package security

import (
	"crypto/hmac"
	"crypto/sha512"
)

var hmacKey []byte

// HmacSum sum hmac
func HmacSum(plain []byte) []byte {
	mac := hmac.New(sha512.New, hmacKey)
	return mac.Sum(plain)
}

// HmacChk chk hmac
func HmacChk(plain, code []byte) bool {
	return hmac.Equal(HmacSum(plain), code)
}
