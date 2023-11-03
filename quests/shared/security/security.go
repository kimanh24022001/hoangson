package security

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"math/big"

	"golang.org/x/crypto/pbkdf2"
	"smatyx.com/shared/cast"
)

func HashPassword(password, salt string) string {
	hash := pbkdf2.Key(
		cast.StringToBytes(password),
		cast.StringToBytes(salt),
		100, 256/8, sha256.New)
	result := base64.StdEncoding.EncodeToString(hash)

	return result
}

// Might be really slow
func NewSalt() string {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	result := make([]byte, 32)
	for i := 0; i < 32; i++ {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		result[i] = letters[num.Int64()]
	}

	return cast.BytesToString(result)
}
