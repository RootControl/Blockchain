package utils

import "strconv"

func IntToHex(num int64) []byte {
	return []byte(strconv.FormatInt(num, 16))
}