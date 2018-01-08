// Package md5 возвращает hash строки
package md5

import (
	"crypto/md5"
	"encoding/hex"
)

// GetHash возвращает hash строки
func GetHash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
