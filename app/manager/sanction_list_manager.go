package manager

import (
	"crypto/sha256"
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"path"
	"time"
)

// SanctionListManager offers functionalities for managing updates to local public sanction list from Open Sanctions.
// Open Sanction: https://www.opensanctions.org/
type SanctionListManager struct{}

// CreateSanctionListManager function creates new SanctionListManager that can be used
// for managing local public sanction list.
func CreateSanctionListManager() *SanctionListManager {
	var sanctionListManager = &SanctionListManager{}

	// Initialize sanction list folder
	if _, err := os.Stat("./sanction-lists"); err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir("sanction-lists", 0777)
			if err != nil {
				panic(fmt.Sprint("error while creating sanction-lists folder %w", err.Error()))
			}
		} else {
			panic(fmt.Sprint("error while searching for sanction-lists folder %w", err.Error()))
		}
	}

	// Initialize sanction list
	if _, err := os.Stat("./sanction-lists/UN_List.csv"); err != nil {
		if os.IsNotExist(err) {
			if _, err := sanctionListManager.GetNewestSanctionsList(); err != nil {
				panic(fmt.Sprint("error while downloading latest sanction list %w", err.Error()))
			}
		} else {
			panic(fmt.Sprint("error while searching for sanction list %w", err.Error()))
		}
	}

	return sanctionListManager
}

func (*SanctionListManager) GetNewestSanctionsList() (string, error) {
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

// LoadSanctionListForNoninteractiveCheck reads a sanctions list csv file where each row contains only one string element and returns a slice of strings.
// additionally, it hashes the elements in the result using SHA-256.
func (s *SanctionListManager) LoadSanctionListForNoninteractiveCheck() ([][]int, error) {
	// Read all the records from the CSV file
	records, err := s.LoadSanctionList()
	if err != nil {
		return nil, err
	}

	// Create a slice to store the hashed strings
	var result [][]int = make([][]int, len(records))

	// Iterate through the records
	for i, record := range records {
		if len(record) > 0 {
			// Hash each string using SHA-256
			hash := sha256.New()
			hash.Write([]byte(record[0]))
			hashedBytes := hash.Sum(nil)

			var intArray []int = make([]int, len(hashedBytes))
			for j, hashedByte := range hashedBytes {
				intArray[j] = int(hashedByte)
			}

			result[i] = intArray
		}
	}

	return result, nil
}

func (*SanctionListManager) LoadSanctionList() ([][]string, error) {
	dir := "./sanction-lists"
	filePath := path.Join(dir, "UN_List.csv")

	// Open the CSV file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Create a new CSV reader
	reader := csv.NewReader(file)

	// Read all the records from the CSV file
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return records, nil
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
