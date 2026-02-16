package service

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"
)

// IMEIService handles IMEI verification logic.
type IMEIService struct{}

// NewIMEIService creates a new IMEIService.
func NewIMEIService() *IMEIService {
	return &IMEIService{}
}

// IMEIResult is the result for a single IMEI check.
type IMEIResult struct {
	IMEI       string `json:"imei"`
	FoundInPDF bool   `json:"found_in_pdf"`
}

// IMEIReport is the full analysis report.
type IMEIReport struct {
	TotalIMEIs int          `json:"total_imeis"`
	Found      int          `json:"found"`
	Missing    int          `json:"missing"`
	Results    []IMEIResult `json:"results"`
}

// Analyze compares IMEIs from CSV against text content of a PDF.
// Note: This is a simplified text-search approach. Full PDF parsing is a separate phase.
func (s *IMEIService) Analyze(csvReader io.Reader, pdfTextContent string) (*IMEIReport, error) {
	reader := csv.NewReader(csvReader)
	reader.TrimLeadingSpace = true

	// Read header
	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("reading CSV header: %w", err)
	}

	// Find IMEI column
	imeiIdx := -1
	for i, col := range header {
		lower := strings.ToLower(strings.TrimSpace(col))
		if lower == "imei" || lower == "imei1" || lower == "imei_number" {
			imeiIdx = i
			break
		}
	}
	if imeiIdx == -1 {
		return nil, fmt.Errorf("CSV must contain an 'imei' column")
	}

	report := &IMEIReport{}
	pdfLower := strings.ToLower(pdfTextContent)

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}

		if imeiIdx >= len(record) {
			continue
		}

		imei := strings.TrimSpace(record[imeiIdx])
		if imei == "" {
			continue
		}

		report.TotalIMEIs++
		found := strings.Contains(pdfLower, strings.ToLower(imei))

		if found {
			report.Found++
		} else {
			report.Missing++
		}

		report.Results = append(report.Results, IMEIResult{
			IMEI:       imei,
			FoundInPDF: found,
		})
	}

	return report, nil
}
