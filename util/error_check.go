package util

import "fmt"

func Check(f func() error) {
	if err := f(); err != nil {
		fmt.Println("Received error:", err)
	}
}
