package service

import (
	"strings"
	"testing"
)

func TestIMEIServiceAnalyze_Success(t *testing.T) {
	svc := NewIMEIService()

	csvContent := `id,imei1,other_col
1,12345678901234,foo
2,99999999999999,bar
3,11111111111111,baz`

	pdfText := `Document text
Here is an IMEI: 123456789012345
And another one not in CSV: 888888888888888
Also one that matches row 3: 111111111111119`

	report, err := svc.Analyze(strings.NewReader(csvContent), pdfText)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if report == nil {
		t.Fatal("expected report to not be nil")
	}

	if report.TotalIMEIs != 3 {
		t.Errorf("expected 3 IMEIs, got %d", report.TotalIMEIs)
	}

	if report.TotalFound != 2 {
		t.Errorf("expected 2 found, got %d", report.TotalFound)
	}

	if report.TotalMissing != 1 {
		t.Errorf("expected 1 missing, got %d", report.TotalMissing)
	}

	// Verify the missing one is 99999999999999
	missingFound := false
	for _, res := range report.Results {
		if res.IMEI14 == "99999999999999" && !res.Found {
			missingFound = true
		}
	}
	if !missingFound {
		t.Error("expected 99999999999999 to be missing")
	}
}

func TestIMEIServiceAnalyze_NoIMEIColumn(t *testing.T) {
	svc := NewIMEIService()

	csvContent := `id,name,value
1,test,100`

	pdfText := `123456789012345`

	_, err := svc.Analyze(strings.NewReader(csvContent), pdfText)
	if err == nil {
		t.Fatal("expected error for missing IMEI column")
	}
}
