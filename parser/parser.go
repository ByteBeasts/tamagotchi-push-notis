package parser

import (
	"encoding/csv"
	"io"
	"log"
)

// CSVData represents the parsed CSV data with headers and rows
type CSVData struct {
	Headers []string
	Rows    [][]string
}

// Parser handles CSV parsing operations
type Parser struct {
	csvReader *csv.Reader
}

// NewParser creates a new CSV parser from an io.Reader
func NewParser(r io.Reader) *Parser {
	reader := csv.NewReader(r)
	reader.LazyQuotes = true
	reader.TrimLeadingSpace = true

	return &Parser{
		csvReader: reader,
	}
}

// Parse reads and parses the CSV data, returning structured data
func (p *Parser) Parse() (*CSVData, error) {
	records, err := p.csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return &CSVData{Headers: []string{}, Rows: [][]string{}}, nil
	}

	// First row is headers, rest are data rows
	headers := records[0]
	rows := records[1:]

	return &CSVData{
		Headers: headers,
		Rows:    rows,
	}, nil
}

// ParseToMaps converts CSV data to a slice of maps for easier manipulation
func (p *Parser) ParseToMaps() ([]map[string]string, error) {
	csvData, err := p.Parse()
	if err != nil {
		return nil, err
	}

	var result []map[string]string

	for _, row := range csvData.Rows {
		rowMap := make(map[string]string)
		for i, value := range row {
			if i < len(csvData.Headers) {
				rowMap[csvData.Headers[i]] = value
			}
		}
		result = append(result, rowMap)
	}

	return result, nil
}

// GetColumn returns all values for a specific column
func (csvData *CSVData) GetColumn(columnName string) []string {
	var columnIndex int = -1

	// Find the column index
	for i, header := range csvData.Headers {
		if header == columnName {
			columnIndex = i
			break
		}
	}

	if columnIndex == -1 {
		log.Printf("Column '%s' not found", columnName)
		return []string{}
	}

	var values []string
	for _, row := range csvData.Rows {
		if columnIndex < len(row) {
			values = append(values, row[columnIndex])
		}
	}

	return values
}

// FilterRows filters rows based on a column value
func (csvData *CSVData) FilterRows(columnName, value string) *CSVData {
	var columnIndex int = -1

	// Find the column index
	for i, header := range csvData.Headers {
		if header == columnName {
			columnIndex = i
			break
		}
	}

	if columnIndex == -1 {
		log.Printf("Column '%s' not found for filtering", columnName)
		return csvData
	}

	var filteredRows [][]string
	for _, row := range csvData.Rows {
		if columnIndex < len(row) && row[columnIndex] == value {
			filteredRows = append(filteredRows, row)
		}
	}

	return &CSVData{
		Headers: csvData.Headers,
		Rows:    filteredRows,
	}
}

// GetRowCount returns the number of data rows
func (csvData *CSVData) GetRowCount() int {
	return len(csvData.Rows)
}

// GetColumnCount returns the number of columns
func (csvData *CSVData) GetColumnCount() int {
	return len(csvData.Headers)
}
