package main

import (
	"fmt"
	"math"
)

type data struct {
	id     int
	width  float64
	size   float64
	length float64
	value  string
}

func distance(predict data, train data) float64 {
	result1 := math.Pow((predict.width - train.width), 2)
	result2 := math.Pow((predict.size - train.size), 2)
	result3 := math.Pow((predict.length - train.length), 2)
	result := math.Sqrt(result1 + result2 + result3)
	return result
}

func predict(predict data, array []data) float64 {
	dist := 0.5
	for i := 0; i < len(array); i++ {
		dist = distance(predict, array[i])
	}
	return dist
}

func main() {
	array := make([]data, 0, 20)
	data1 := data{0, 100, 100, 100, "red"}
	data2 := data{1, 89, 88, 91, "red"}
	data3 := data{2, 102, 103, 104, "red"}
	data4 := data{3, 90, 91, 92, "red"}
	data5 := data{4, 79, 60, 80, "red"}

	data6 := data{5, 0, 0, 0, "blue"}
	data7 := data{6, 10, 16, 11, "blue"}
	data8 := data{8, 6, 13, 32, "blue"}

	predict := data{3, 80, 70, 104, "None"}

	array = append(array, data1, data2, data3, data4, data5, data6, data7, data8)

	fmt.Println(predict.value)
}
