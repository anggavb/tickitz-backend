package pkg

import (
	"encoding/base64"
	"math/rand"
)

func GenerateRandomToken(n int) string {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return base64.URLEncoding.EncodeToString(b)
}
