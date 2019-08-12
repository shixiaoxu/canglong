package main

import (
	"fmt"
	kaptcha "canglong/kaptcha/common/generate"
)

func main() {
	fmt.Println("hello kaptcha")
	fmt.Println(kaptcha.GenDecimalCode4())
	fmt.Println(kaptcha.GenDecimalCode6())
}
