package utils

import (
	"fmt"
	"testing"
)

func TestRandInt(t *testing.T) {
	for i := 0; i < 5; i++ {
		n := RandInt(0, 2)
		fmt.Println(n)
	}
}
