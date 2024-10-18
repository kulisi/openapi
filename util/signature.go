package util

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

func asd(keyId string, secret string, timeRub string, urlPath string) string {
	random := RandStr(5)
	value := keyId + "&" + timeRub + "&" + urlPath
	key := secret + random
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(value))

	sha := hex.EncodeToString(mac.Sum(nil))

	return random + base64.StdEncoding.EncodeToString([]byte(sha))
}
