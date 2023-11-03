package random

import (
	"math/rand"
	"time"

	"smatyx.com/internal/cast"
)

func randRand() *rand.Rand {
	source := rand.NewSource(time.Now().UnixNano())
	result := rand.New(source)

	return result
}

var letterRunes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func String(n int) string {
    b := make([]byte, n)
    for i := range b {
		ranInt := rand.Intn(len(letterRunes))
        b[i] = letterRunes[ranInt]
    }
    return cast.BytesToString(b)
}
