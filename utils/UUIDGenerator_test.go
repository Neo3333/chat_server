package utils_test

import (
	"../utils"
	"fmt"
	"testing"
)

func TestNewUUIDGenerator(t *testing.T) {
	UUIDFactory := utils.NewUUIDGenerator("Neo")
	for i := 0; i < 50; i++{
		fmt.Println(UUIDFactory.Get())
	}
}
