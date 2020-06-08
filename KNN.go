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
	name   string
}

type distances struct {
	id       int
	distance float64
}

func calculateDistance(predict data, train data) float64 {
	result1 := math.Pow((predict.width - train.width), 2)
	result2 := math.Pow((predict.size - train.size), 2)
	result3 := math.Pow((predict.length - train.length), 2)
	result := math.Sqrt(result1 + result2 + result3)
	return result
}

func getAllDistances(predict data, array []data) []distances {
	distArray := make([]distances, 0, 20)
	for i := 0; i < len(array); i++ {
		dist := calculateDistance(predict, array[i])
		distStruc := distances{array[i].id, dist}
		distArray = append(distArray, distStruc)
	}
	return distArray
}

func bubbleSort(dists []distances) []distances {
	for i := 0; i < len(dists)-1; i++ {
		for j := 0; j < len(dists)-i-1; j++ {
			if dists[j].distance > dists[j+1].distance {
				aux := dists[j]
				dists[j] = dists[j+1]
				dists[j+1] = aux
			}
		}
	}
	return dists
}

func findClosestGroups(dists []distances, k int) []distances {
	sortedDists := bubbleSort(dists)
	minDists := make([]distances, 0, 20)
	for j := 0; j < k; j++ {
		minDists = append(minDists, sortedDists[j])
	}
	return minDists
}

func predict(predict data, array []data) string {
	distArray := getAllDistances(predict, array)
	minDists := findClosestGroups(distArray, 2)
	fmt.Println(minDists[0].distance)
	fmt.Println(minDists[1].distance)
	return "red"
}

func main() {
	array := make([]data, 0, 20)
	data1 := data{0, 100, 105, 100, "red"}
	data2 := data{1, 105, 101, 105, "red"}
	data3 := data{2, 102, 103, 104, "red"}
	data4 := data{3, 90, 91, 92, "red"}
	data5 := data{4, 79, 60, 80, "red"}

	data6 := data{5, 0, 0, 0, "blue"}
	data7 := data{6, 10, 16, 11, "blue"}
	data8 := data{8, 6, 13, 32, "blue"}

	wantPredict := data{3, 100, 100, 100, "None"}

	array = append(array, data1, data2, data3, data4, data5, data6, data7, data8)

	typeOf := predict(wantPredict, array)
	println(typeOf)
}
