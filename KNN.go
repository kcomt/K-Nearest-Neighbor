package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

type data struct {
	id          int
	sepalLength float64
	sepalWidth  float64
	petalLength float64
	petalWidth  float64
	class       string
}

type distances struct {
	id       int
	distance float64
}

func calculateDistance(predict data, train data, distanceBuffer chan distances) {
	result1 := math.Pow((predict.sepalLength - train.sepalLength), 2)
	result2 := math.Pow((predict.sepalWidth - train.sepalWidth), 2)
	result3 := math.Pow((predict.petalLength - train.petalLength), 2)
	result4 := math.Pow((predict.petalWidth - train.petalWidth), 2)
	result := math.Sqrt(result1 + result2 + result3 + +result4)

	distStruct := distances{train.id, result}
	distanceBuffer <- distStruct
}

func getAllDistances(predict data, array []data) []distances {
	distArray := make([]distances, 0, 20)
	distanceBuffer := make(chan distances, 8)
	for i := 0; i < len(array); i++ {
		go calculateDistance(predict, array[i], distanceBuffer)
	}
	for i := 0; i < len(array); i++ {
		distArray = append(distArray, <-distanceBuffer)
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

func load(dataBuffer chan data, r *csv.Reader, i int) {

	record, err := r.Read()
	if err != nil {
		log.Fatal(err)
	}
	noCommasRow := strings.Split(record[0], ",")
	sepalLength, _ := strconv.ParseFloat(noCommasRow[0], 8)
	sepalWidth, _ := strconv.ParseFloat(noCommasRow[1], 8)
	petalLength, _ := strconv.ParseFloat(noCommasRow[2], 8)
	petalWidth, _ := strconv.ParseFloat(noCommasRow[3], 8)
	data := data{i, sepalLength, sepalWidth, petalLength, petalWidth, noCommasRow[4]}
	dataBuffer <- data
}

func main() {
	array := make([]data, 0, 2000)
	csvfile, err := os.Open("dataSetIris.csv")
	dataBuffer := make(chan data, 150)
	if err != nil {
		log.Fatalln("No se pudo abrir el archivo", err)
	}
	r := csv.NewReader(csvfile)
	for i := 0; i < 150; i++ {
		go load(dataBuffer, r, i)
	}
	for i := 0; i < 150; i++ {
		array = append(array, <-dataBuffer)
	}
	fmt.Println("yeeer")
}
