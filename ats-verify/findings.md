# Technical Findings & Decisions

## 1. PDF Parsing Strategy
- Need a library that can extract text accurately from "Graph 31" of the declaration.
- **Go Libs:** `unidoc/unipdf` (commercial), `dslipak/pdf`, or wrapper around `pdftotext` (poppler).
- **Python Libs:** `PyMuPDF` (fitz), `pdfplumber`.
- *Decision:* Start with Go's `ledongthuc/pdf` or similar open-source lib. If text extraction is poor, invoke a Python script sidecar.

## 2. IMEI Logic
- IMEI is 15 digits.
- The 15th digit is a Luhn check digit.
- Logic: "Find 14-digit substring in PDF text".
- This implies Regex: `\b\d{14}\d?\b`.

## 3. Consignment Note Logic (SNT)
- Marketplace uploads include `SNT`.
- Need to validate if SNT format is consistent.

## 4. Duplicate Logic (SQL)
- Need `UPSERT` logic (ON CONFLICT).
- Constraint: `track_number` must be unique? No, duplicates allowed if `used=false`.
- Actually, logic says: "If track exists and used=false -> UPDATE".
- "If track exists and used=true -> ERROR".
- So `track_number` is the key.
