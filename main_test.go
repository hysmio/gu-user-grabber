package main

import (
	"fmt"
	"testing"
)

func TestGetID(t *testing.T) {
	id, err := GetID(`C:\Users\Hayden\AppData\LocalLow\Immutable\gods\debug.log`)
	if err != nil {
		panic(err)
	}

	fmt.Println(id)
}
