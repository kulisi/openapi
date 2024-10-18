package util

import "math/rand"

// RandStr 获取随机字符串
func RandStr(length int) string {
	char := []rune("1234567890qwertyuiopasdfghjklzxcvbnm")
	rChar := make([]rune, length)
	for i := range rChar {
		rChar[i] = char[rand.Intn(len(char))]
	}
	return string(rChar)
}
