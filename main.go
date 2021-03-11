package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
)

var config map[string]interface{}
var initialSellPrice float64
var initialBuyPrice float64

func main() {
	loadConfig()
	readExistingData()
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
	fmt.Println(initialText)

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

}

func writeNewData() {
	text := fmt.Sprintf("<Trader> %s\n<Category> Currency\n%s, *, *, %d, %d\n<FileEnd>",
		config["trader_name"],
		config["currency_name"],
		int(50),
		int(60),
	)
	dataBytes := []byte(text)
	ioutil.WriteFile(fmt.Sprintf("%v", config["forex_trader_file_path"]), dataBytes, 0)
}
