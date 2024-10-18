package util

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

// Md5 同过md5加密文本
func Md5(text string) string {
	hash := md5.New()
	hash.Write([]byte(text))
	return hex.EncodeToString(hash.Sum(nil))
}

// Sha256Hmac 通过 Sha256Hmac 加密文本
func Sha256Hmac(key string, value string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(value))

	sha := hex.EncodeToString(mac.Sum(nil))

	return base64.StdEncoding.EncodeToString([]byte(sha))
}
