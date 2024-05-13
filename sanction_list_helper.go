package main

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"path"
	"time"
)

func getNewestSanctionsList() (string, error) {
	// Open sanctions url contains date
	// Start from today's date
	currentDate := time.Now()

	var (
		fileName string
	)

	for {
		// Construct URL with current date
		dateString := currentDate.Format("20060102")
		url := constructURL(dateString)

		// Try downloading CSV with current date
		csvData, err := downloadCSVWithRetry(url)
		if err == nil {
			// Filter columns
			filteredData := filterColumns(csvData)

			// Remove the first row (column names) from the filtered data
			if len(filteredData) > 0 {
				filteredData = filteredData[1:]
			}

			// Save filtered CSV content
			fileName, err = saveCSV(filteredData, dateString)
			if err != nil {
				return "", err
			}

			// Exit the loop if CSV is successfully downloaded
			break
		}

		// Move to the previous day
		currentDate = currentDate.AddDate(0, 0, -1)
	}

	return fileName, nil
}

func saveCSV(records [][]string, date string) (string, error) {
	dir := "./sanction-lists"
	fileName := fmt.Sprintf("UN_list_%s.csv", date)

	// Save to UN_List.csv so we dont have a lot of unnecessary files
	filePath := path.Join(dir, "UN_List.csv")

	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, row := range records {
		if err := writer.Write(row); err != nil {
			return "", err
		}
	}

	return fileName, nil
}

func filterColumns(records [][]string) [][]string {
	// We only need "name" column for the sanctions check
	var filteredRecords [][]string

	// Find the index of the "name" column
	var nameIndex int
	for i, columnName := range records[0] {
		if columnName == "name" {
			nameIndex = i
			break
		}
	}

	// Copy only the "name" column to filteredRecords
	for _, row := range records {
		filteredRow := []string{row[nameIndex]}
		filteredRecords = append(filteredRecords, filteredRow)
	}

	return filteredRecords
}

func constructURL(date string) string {
	return fmt.Sprintf("https://data.opensanctions.org/datasets/%s/un_sc_sanctions/targets.simple.csv", date)
}

func downloadCSVWithRetry(url string) ([][]string, error) {
	// Download CSV file
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Read CSV content
	reader := csv.NewReader(response.Body)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return records, nil
}
