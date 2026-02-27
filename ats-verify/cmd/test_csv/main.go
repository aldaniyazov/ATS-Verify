package main

import (
	"fmt"
	"log"
	"os"

	"ats-verify/internal/service"
)

func main() {
	f, err := os.Open("/Users/aldaniyazov/Documents/projects/example/marketplaces - dump.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	csvReader, err := service.NewRobustCSVReader(f)
	if err != nil {
		log.Fatal(err)
	}
	headerRow, _ := csvReader.Read()
	fmt.Printf("Headers: %v\n", headerRow)
	record, _ := csvReader.Read()
	fmt.Printf("Row 1: %v\n", record)
}
