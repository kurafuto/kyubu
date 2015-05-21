package main

import "fmt"

var tmpCount = 0

func tmpVar() string {
	tmpCount++
	return fmt.Sprintf("tmp%d", tmpCount-1)
}
