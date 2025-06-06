package main

import (
	"fmt"
	"tic-tac-toe/randomname/services"
)

func main() {
	nickname := services.RandomNickName()
	fmt.Println(nickname)
}
