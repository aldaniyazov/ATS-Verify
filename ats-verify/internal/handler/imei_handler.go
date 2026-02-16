package handler

import (
	"io"
	"net/http"

	"ats-verify/internal/middleware"
	"ats-verify/internal/models"
	"ats-verify/internal/service"
)

// IMEIHandler handles IMEI verification endpoints.
type IMEIHandler struct {
	imeiService *service.IMEIService
}

// NewIMEIHandler creates a new IMEIHandler.
func NewIMEIHandler(imeiService *service.IMEIService) *IMEIHandler {
	return &IMEIHandler{imeiService: imeiService}
}

// RegisterRoutes registers IMEI routes.
func (h *IMEIHandler) RegisterRoutes(mux *http.ServeMux, authMw func(http.Handler) http.Handler) {
	roleMw := middleware.RequireRole(models.RoleCustoms, models.RolePaidUser, models.RoleAdmin)
	mux.Handle("POST /api/v1/imei/analyze", authMw(roleMw(http.HandlerFunc(h.Analyze))))
}

// Analyze handles POST /api/v1/imei/analyze (multipart: csv_file + pdf_file)
func (h *IMEIHandler) Analyze(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form (max 50MB total)
	if err := r.ParseMultipartForm(50 << 20); err != nil {
		Error(w, http.StatusBadRequest, "failed to parse form: "+err.Error())
		return
	}

	// Get CSV file
	csvFile, _, err := r.FormFile("csv_file")
	if err != nil {
		Error(w, http.StatusBadRequest, "csv_file is required")
		return
	}
	defer csvFile.Close()

	// Get PDF file
	pdfFile, _, err := r.FormFile("pdf_file")
	if err != nil {
		Error(w, http.StatusBadRequest, "pdf_file is required")
		return
	}
	defer pdfFile.Close()

	// Read PDF content as text (simplified approach)
	pdfBytes, err := io.ReadAll(pdfFile)
	if err != nil {
		Error(w, http.StatusInternalServerError, "failed to read PDF file")
		return
	}

	// Note: This is a simplified text extraction.
	// Real PDF parsing would use a library like pdfcpu or unidoc.
	pdfText := string(pdfBytes)

	report, err := h.imeiService.Analyze(csvFile, pdfText)
	if err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	JSON(w, http.StatusOK, report)
}
