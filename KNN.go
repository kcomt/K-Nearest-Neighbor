package main

import (
	"fmt"
	"math"
)

type data struct {
	width  float64
	size   float64
	length float64
	value  string
}

func distance(x float64, y float64, x2 float64, y2 float64) float64 {
	result1 := math.Pow((x - x2), 2)
	result2 := math.Pow((y - y2), 2)
	result := math.Sqrt(result1 + result2)
	return result
}

func predict(predict data, array []data) data {
	for v := range array {
		fmt.Println(array[v].length)
	}
	return predict
}

func main() {
	array := make([]data, 0, 20)
	data1 := data{100, 100, 100, "red"}
	data2 := data{89, 88, 91, "red"}
	data3 := data{102, 103, 104, "red"}
	data4 := data{90, 91, 92, "red"}
	data5 := data{79, 60, 80, "red"}

	data6 := data{0, 0, 0, "blue"}
	data7 := data{10, 16, 11, "blue"}
	data8 := data{6, 13, 32, "blue"}

	predict := data{80, 70, 104, "None"}

	array = append(array, data1, data2, data3, data4, data5, data6, data7, data8)

	fmt.Println(predict.value)
}
