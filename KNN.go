package main

import (
	"fmt"
	"math"
)

type data struct {
	width  float64
	size   float64
	length float64
	name   string
}

func distance(x float64, y float64, x2 float64, y2 float64) {
	result1 := math.Pow((x - x2), 2)
	result2 := math.Pow((y - y2), 2)
	result := math.Sqrt(result1 + result2)
	fmt.Println(result)
}

func main() {
	distance(3, 4, 0, 0)
}
