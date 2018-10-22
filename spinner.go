package moni

import (
	"fmt"
	"time"
)

func FibinucciSpinner() {
	fmt.Println("Starting spinner... ")
	go spinner(300 * time.Millisecond)
	const n = 45
	fibN := fib(n) // slow
	fmt.Printf("  fib(%d) = %d\n", n, fibN)
}

func spinner(delay time.Duration) {
	for {
		for _, r := range `-\|/` {
			fmt.Printf("\r%c", r)
			time.Sleep(delay)
		}
	}
}

// fib
func fib(x int) int {
	if x < 2 {
		return x
	}
	return fib(x-1) + fib(x-2)
}
