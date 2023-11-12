package util

import (
	"crypto/rand"
	"encoding/hex"
	"math"
	"strconv"
)

// IsNumeric judges whether given string is numeric or not.
func IsNumeric(number string) bool {
	_, err := strconv.Atoi(number)
	return err == nil
}

// ConvertToInt converts given string to int.
func ConvertToInt(number string) int {
	value, err := strconv.Atoi(number)
	if err != nil {
		return 0
	}
	return value
}

// ConvertToUint converts given string to uint.
func ConvertToUint(number string) uint {
	return uint(ConvertToInt(number))
}

func RandomBase16String(l int) string {
	buff := make([]byte, int(math.Ceil(float64(l)/2)))
	rand.Read(buff)
	str := hex.EncodeToString(buff)
	return str[:l]
}
