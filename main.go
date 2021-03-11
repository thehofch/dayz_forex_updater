package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"time"
)

var config map[string]interface{}
var initialSellPrice float64
var initialBuyPrice float64
var newSellValue int
var newBuyValue int

func main() {
	loadConfig()
	readExistingData()
	calculateNewData()
	writeNewData()
}

func loadConfig() {
	jsonFile, err := os.Open("config.json")

	if err != nil {
		fmt.Println(err)
	}

	defer func() {
		if err = jsonFile.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal([]byte(byteValue), &config)
}

func readExistingData() {
	filePath := config["forex_trader_file_path"]

	file, err := ioutil.ReadFile(fmt.Sprintf("%v", filePath))
	if err != nil {
		log.Fatal(err)
	}

	initialText := string(file)

	r, err := regexp.Compile(`(\d+),\s+(\d+)`)

	if err != nil {
		fmt.Printf("There is a problem with your regexp.\n")
		return
	}

	res := r.FindStringSubmatch(initialText)

	if len(res) != 3 {
		os.Exit(3)
	}

	sellPrice, err := strconv.ParseFloat(res[1], 0)

	if err != nil {
		log.Fatal(err)
	} else {
		initialSellPrice = sellPrice
	}

	buyPrice, err := strconv.ParseFloat(res[2], 0)

	if err != nil {
		log.Fatal(err)
	} else {
		initialBuyPrice = buyPrice
	}
}

func calculateNewData() {
	rand.Seed(time.Now().UnixNano())
	rand := rand.Float64() * 100
	randInt := int(rand)

	userChanceOfIncreasing, err := strconv.ParseFloat(fmt.Sprintf("%v", config["chanse_of_increasing"]), 0)

	if err != nil {
		log.Fatal(err)
	}

	userChanceOfIncreasingInt := int(userChanceOfIncreasing)

	increasing := true

	if randInt > userChanceOfIncreasingInt {
		increasing = false
	}

	newSellValue = changeValue(increasing)
	newBuyValue = int(float64(newSellValue) * 0.99)
}

func changeValue(increase bool) int {
	minThreshold, err := strconv.ParseFloat(fmt.Sprintf("%v", config["min_threshold"]), 0)

	if err != nil {
		log.Fatal(err)
	}

	minThresholdInt := int(minThreshold)

	maxThreshold, err := strconv.ParseFloat(fmt.Sprintf("%v", config["max_threshold"]), 0)

	if err != nil {
		log.Fatal(err)
	}

	maxThresholdInt := int(maxThreshold)

	number := rand.Intn(maxThresholdInt-minThresholdInt) + minThresholdInt

	var newPrice float64

	if increase == true {
		newPrice = initialSellPrice * (float64(number)/100 + 1.0)
	} else {
		newPrice = initialSellPrice * (1 - float64(number)/100)
	}

	return int(newPrice)
}

func writeNewData() {
	text := fmt.Sprintf("<Trader> %s\n<Category> Currency\n%s, *, *, %d, %d\n<FileEnd>",
		config["trader_name"],
		config["currency_name"],
		int(newSellValue),
		int(newBuyValue),
	)
	dataBytes := []byte(text)
	ioutil.WriteFile(fmt.Sprintf("%v", config["forex_trader_file_path"]), dataBytes, 0)
}
