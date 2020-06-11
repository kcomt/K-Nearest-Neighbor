package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
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
	quality            int
}

type dataJson struct {
	Id                 string `json:"id"`
	FixedAcidity       string `json:"fixedAcidity"`
	VolatileAcidity    string `json:"volatileAcidity"`
	CitricAcid         string `json:"citricAcid"`
	ResidualSugar      string `json:"residualSugar"`
	Chlorides          string `json:"chlorides"`
	TotalSulfurDioxide string `json:"totalSulfurDioxide"`
	Density            string `json:"density"`
	Ph                 string `json:"ph"`
	Alcohol            string `json:"alcohol"`
	Quality            string `json:"quality"`
}

type distances struct {
	id       int
	distance float64
	quality  int
}

var arrayOfData = make([]data, 0, 2000)

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

	distStruct := distances{train.id, result, train.quality}
	distanceBuffer <- distStruct
}

func getAllDistances(predict data) []distances {
	distArray := make([]distances, 0, len(arrayOfData))
	distanceBuffer := make(chan distances)
	for i := 0; i < len(arrayOfData); i++ {
		go calculateDistance(predict, arrayOfData[i], distanceBuffer)
	}
	for i := 0; i < len(arrayOfData); i++ {
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

func predict(predict data, k int) int {
	distArray := getAllDistances(predict)
	minDists := findClosestGroups(distArray, k)
	score := make([]int, 11, 11)
	for i := 0; i < len(score); i++ {
		score[i] = 0
	}
	for i := 0; i < len(minDists); i++ {
		score[minDists[i].quality] = score[minDists[i].quality] + 1
	}
	quality := 1
	maxi := score[1]
	for i := 1; i < len(score); i++ {
		if score[i] >= maxi {
			maxi = score[i]
			quality = i
		}
	}
	return quality
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
	quality, _ := strconv.Atoi(noCommasRow[11])
	data := data{i, fixedAcidity, volatileAcidity, citricAcid, residualSugar, chlorides, totalSulfurDioxide, density, ph, alcohol, quality}
	dataBuffer <- data
}

func train() {
	dataBuffer := make(chan data, 1600)
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
		arrayOfData = append(arrayOfData, <-dataBuffer)
	}
}

func homeLink(w http.ResponseWriter, r *http.Request) {
	s := "Welcome"
	fmt.Fprintf(w, s)
}

func createEvent(w http.ResponseWriter, r *http.Request) {
	var predictPlease dataJson
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the event title and description only in order to update")
	}
	json.Unmarshal(reqBody, &predictPlease)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(predictPlease)

	fixedAcidity, _ := strconv.ParseFloat(predictPlease.FixedAcidity, 8)
	volatileAcidity, _ := strconv.ParseFloat(predictPlease.VolatileAcidity, 8)
	citricAcid, _ := strconv.ParseFloat(predictPlease.CitricAcid, 8)
	residualSugar, _ := strconv.ParseFloat(predictPlease.ResidualSugar, 8)
	chlorides, _ := strconv.ParseFloat(predictPlease.Chlorides, 8)
	totalSulfurDioxide, _ := strconv.ParseFloat(predictPlease.TotalSulfurDioxide, 8)
	density, _ := strconv.ParseFloat(predictPlease.Density, 8)
	ph, _ := strconv.ParseFloat(predictPlease.Ph, 8)
	alcohol, _ := strconv.ParseFloat(predictPlease.Alcohol, 8)
	quality, _ := strconv.Atoi(predictPlease.Quality)

	data := data{2000, fixedAcidity, volatileAcidity, citricAcid, residualSugar, chlorides, totalSulfurDioxide, density, ph, alcohol, quality}
	qualityPrediction := predict(data, 3)
	fmt.Println(qualityPrediction)
}

func main() {
	train()
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/data", createEvent).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", router))
}
