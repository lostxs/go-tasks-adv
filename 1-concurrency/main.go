package main

import (
	"fmt"
	"math/rand"
)

func main() {
	numbers := make([]int, 10)
	squares := make([]int, 0, 10)

	numCh := make(chan int)
	sqCh := make(chan int)

	go func() {
		for i := range numbers {
			num := rand.Intn(101)
			numbers[i] = num
			numCh <- num
		}
		close(numCh)
	}()

	go func() {
		for num := range numCh {
			sqCh <- num * num
		}
		close(sqCh)
	}()

	for sq := range sqCh {
		squares = append(squares, sq)
	}

	fmt.Println("Исходные числа:", numbers)
	fmt.Println("Квадраты чисел:", squares)
}
