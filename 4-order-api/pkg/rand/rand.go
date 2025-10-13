package rand

import (
	"fmt"
	"math/rand"
)

func RandSession() string {
	sessionId := "aAbBcCdDeEfFgGhHiIjJkKlLmMnNoOpPqQrRsStTuUvVwWxXyYzZ0123456789"
	readySesId := make([]byte, 12)
	for i := 0; i < 12; i++ {
		readySesId[i] = sessionId[rand.Intn(len(sessionId))]

	}
	fmt.Println(string(readySesId))
	return string(readySesId)
}
