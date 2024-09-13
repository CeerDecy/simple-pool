package tools

import (
	"fmt"
	"testing"
)

func TestRandString(t *testing.T) {
	fmt.Println(RandString(1))
	fmt.Println(RandString(2))
	fmt.Println(RandString(4))
	fmt.Println(RandString(8))
	fmt.Println(RandString(16))
	fmt.Println(RandString(32))
	fmt.Println(RandString(64))
	fmt.Println(RandString(128))
	fmt.Println(RandString(256))
}
