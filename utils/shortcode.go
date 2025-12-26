// internal/utils/shortcode.go
package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
)

const (
	base62Chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	DefaultLength = 7
)

func GenerateShortCode(length int) (string, error) {
	if length <= 0 {
		length = DefaultLength
	}

	result := make([]byte, length)
	charsLen := big.NewInt(int64(len(base62Chars)))

	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, charsLen)
		if err != nil {
			return "", err
		}
		result[i] = base62Chars[num.Int64()]
	}

	return string(result), nil
}

func EncodeBase62(num uint64) string {
	if num == 0 {
		return string(base62Chars[0])
	}

	result := make([]byte, 0)
	base := uint64(len(base62Chars))

	for num > 0 {
		remainder := num % base
		result = append([]byte{base62Chars[remainder]}, result...)
		num = num / base
	}

	return string(result)
}

func DecodeBase62(str string) (uint64, error) {
	var num uint64
	base := uint64(len(base62Chars))

	charIndex := make(map[byte]uint64)
	for i := 0; i < len(base62Chars); i++ {
		charIndex[base62Chars[i]] = uint64(i)
	}

	for i := 0; i < len(str); i++ {
		char := str[i]
		value, exists := charIndex[char]
		if !exists {
			return 0, fmt.Errorf("invalid character: %c", char)
		}
		num = num*base + value
	}

	return num, nil
}

func ValidateShortCode(code string) bool {
	if len(code) < 1 || len(code) > 10 {
		return false
	}

	for _, char := range code {
		if !strings.ContainsRune(base62Chars, char) {
			return false
		}
	}

	return true
}
