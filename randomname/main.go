package main

import (
	"fmt"
	"randomname/services"
)

func main() {
	nickname := services.RandomNickName()
	fmt.Println(nickname)
}
