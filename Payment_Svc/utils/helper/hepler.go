package helper

import "crypto/rand"

func RandString() string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	bytes := make([]byte, 8)
	_, err := rand.Read(bytes)
	if err != nil {
		return ""
	}

	for i := 0; i < 8; i++ {
		bytes[i] = letters[bytes[i]%byte(len(letters))]
	}
	return string(bytes)
}
