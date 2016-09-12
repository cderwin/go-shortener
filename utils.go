package main

import (
	"hash/crc32"
	"time"
)

var alphabet = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func base62Encode(src uint32) string {
	result := make([]rune, 0, 13)
	for src > 0 {
		value := src % 62
		src = (src - value) / 62
		result = append(result, alphabet[value])
	}

	size := len(result)
	buffer := make([]rune, size)
	for i, num := range result {
		buffer[size-1-i] = num
	}
	return string(buffer)
}

func hashUrl(long_url string) (short_url string) {
	checksum := crc32.ChecksumIEEE([]byte(long_url))
	short_url = base62Encode(checksum)
	return
}

// Clock interface for easy testing

type Clock interface {
	UTCNow() time.Time
}

type SystemClock struct{}

func NewSystemClock() SystemClock {
	return SystemClock{}
}

func (c SystemClock) UTCNow() time.Time {
	return time.Now().UTC()
}
