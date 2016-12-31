package main

// // no C code needed.
import "C"
import "fmt"

var V int

func F() string { return fmt.Sprintf("Hello, number %d", V) }
