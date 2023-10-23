package refreshtoken

import (
	"math/rand"
	"time"
)

func Create() (string, error) {
	refreshTokenBytes := make([]byte, 32)

	source := rand.NewSource(time.Now().Unix())
	randomSource := rand.New(source)

	_, err := randomSource.Read(refreshTokenBytes)
	if err != nil {
		return "", err
	}

	return string(refreshTokenBytes), nil
}
