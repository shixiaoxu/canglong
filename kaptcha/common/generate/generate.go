package generate

import (
"fmt"
"math/rand"
"time"
)

func GenDecimalCode4() (string) {
	return fmt.Sprintf("%0.4d",
		rand.New(rand.NewSource(time.Now().UnixNano())).Int63n(10000))
}

func GenDecimalCode6() (string) {
	return fmt.Sprintf("%0.6d",
		rand.New(rand.NewSource(time.Now().UnixNano())).Int63n(1000000))
}