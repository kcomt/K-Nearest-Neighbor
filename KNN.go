package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

type data struct {
	id                 int
	fixedAcidity       float64
	volatileAcidity    float64
	citricAcid         float64
	residualSugar      float64
	chlorides          float64
	totalSulfurDioxide float64
	density            float64
	ph                 float64
	alcohol            float64
	quality            float64
}

type distances struct {
	id       int
	distance float64
}

func calculateDistance(predict data, train data, distanceBuffer chan distances) {
	result1 := math.Pow((predict.fixedAcidity - train.fixedAcidity), 2)
	result2 := math.Pow((predict.volatileAcidity - train.volatileAcidity), 2)
	result3 := math.Pow((predict.citricAcid - train.citricAcid), 2)
	result4 := math.Pow((predict.residualSugar - train.residualSugar), 2)
	result5 := math.Pow((predict.chlorides - train.chlorides), 2)
	result6 := math.Pow((predict.totalSulfurDioxide - train.totalSulfurDioxide), 2)
	result7 := math.Pow((predict.density - train.density), 2)
	result8 := math.Pow((predict.ph - train.ph), 2)
	result9 := math.Pow((predict.alcohol - train.alcohol), 2)
	result := math.Sqrt(result1 + result2 + result3 + +result4 + result5 + result6 + result7 + result8 + result9)

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

func load(i int, noCommasRow []string, dataBuffer chan data) {
	fixedAcidity, _ := strconv.ParseFloat(noCommasRow[0], 8)
	volatileAcidity, _ := strconv.ParseFloat(noCommasRow[1], 8)
	citricAcid, _ := strconv.ParseFloat(noCommasRow[2], 8)
	residualSugar, _ := strconv.ParseFloat(noCommasRow[3], 8)
	chlorides, _ := strconv.ParseFloat(noCommasRow[4], 8)
	totalSulfurDioxide, _ := strconv.ParseFloat(noCommasRow[6], 8)
	density, _ := strconv.ParseFloat(noCommasRow[7], 8)
	ph, _ := strconv.ParseFloat(noCommasRow[8], 8)
	alcohol, _ := strconv.ParseFloat(noCommasRow[10], 8)
	quality, _ := strconv.ParseFloat(noCommasRow[11], 8)
	data := data{i, fixedAcidity, volatileAcidity, citricAcid, residualSugar, chlorides, totalSulfurDioxide, density, ph, alcohol, quality}
	fmt.Println(i)
	dataBuffer <- data
}

func main() {
	dataBuffer := make(chan data, 1600)
	array := make([]data, 0, 2000)
	csvfile, err := os.Open("wine.csv")
	if err != nil {
		log.Fatalln("No se pudo abrir el archivo", err)
	}
	r := csv.NewReader(csvfile)
	for i := 0; i < 1600; i++ {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		noCommasRow := strings.Split(record[0], ";")
		go load(i, noCommasRow, dataBuffer)
	}
	for i := 0; i < 1600; i++ {
		array = append(array, <-dataBuffer)
	}
}
