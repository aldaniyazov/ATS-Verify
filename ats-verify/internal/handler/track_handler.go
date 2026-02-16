package handler

import (
	"net/http"
	"strings"

	"ats-verify/internal/middleware"
	"ats-verify/internal/models"
	"ats-verify/internal/service"
)

// TrackHandler handles track search and tracking endpoints.
type TrackHandler struct {
	parcelService *service.ParcelService
}

// NewTrackHandler creates a new TrackHandler.
func NewTrackHandler(parcelService *service.ParcelService) *TrackHandler {
	return &TrackHandler{parcelService: parcelService}
}

// RegisterRoutes registers track routes.
func (h *TrackHandler) RegisterRoutes(mux *http.ServeMux, authMw func(http.Handler) http.Handler) {
	mux.Handle("POST /api/v1/track/bulk", authMw(
		middleware.RequireRole(models.RoleATSStaff, models.RoleAdmin, models.RoleCustoms)(http.HandlerFunc(h.BulkSearch)),
	))
	mux.Handle("GET /api/v1/tracking/{track}", authMw(http.HandlerFunc(h.GetTracking)))
}

// bulkSearchRequest is the payload for bulk track search.
type bulkSearchRequest struct {
	Tracks []string `json:"tracks"`
}

// BulkSearch handles POST /api/v1/track/bulk
func (h *TrackHandler) BulkSearch(w http.ResponseWriter, r *http.Request) {
	var req bulkSearchRequest
	if err := Decode(r, &req); err != nil {
		Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if len(req.Tracks) == 0 {
		Error(w, http.StatusBadRequest, "tracks array is required")
		return
	}

	// Cap at 500 tracks per request
	if len(req.Tracks) > 500 {
		req.Tracks = req.Tracks[:500]
	}

	results, err := h.parcelService.BulkTrackLookup(r.Context(), req.Tracks)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	found := 0
	for _, r := range results {
		if r.Found {
			found++
		}
	}

	JSON(w, http.StatusOK, map[string]interface{}{
		"total":   len(results),
		"found":   found,
		"missing": len(results) - found,
		"results": results,
	})
}

// trackingEvent is a mock tracking event.
type trackingEvent struct {
	Status      string `json:"status"`
	Description string `json:"description"`
	Location    string `json:"location"`
	Date        string `json:"date"`
	Source      string `json:"source"`
}

// GetTracking handles GET /api/v1/tracking/{track}
// Returns mock tracking data (real Kazpost/CDEK integration is a future phase).
func (h *TrackHandler) GetTracking(w http.ResponseWriter, r *http.Request) {
	track := r.PathValue("track")
	if track == "" {
		Error(w, http.StatusBadRequest, "track number is required")
		return
	}
	track = strings.TrimSpace(track)

	// Check if parcel exists
	results, err := h.parcelService.BulkTrackLookup(r.Context(), []string{track})
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	if len(results) == 0 || !results[0].Found {
		Error(w, http.StatusNotFound, "parcel not found in database")
		return
	}

	// Return mock tracking events
	events := []trackingEvent{
		{Status: "ACCEPTED", Description: "Посылка принята на склад отправителя", Location: "Шэньчжэнь, Китай", Date: "2026-02-10T10:00:00Z", Source: "mock"},
		{Status: "IN_TRANSIT", Description: "Отправлена международным рейсом", Location: "Шэньчжэнь → Алматы", Date: "2026-02-11T14:30:00Z", Source: "mock"},
		{Status: "CUSTOMS", Description: "Прибыла на таможню", Location: "Алматы, Казахстан", Date: "2026-02-13T09:00:00Z", Source: "mock"},
		{Status: "CLEARED", Description: "Таможенное оформление завершено", Location: "Алматы, Казахстан", Date: "2026-02-14T16:00:00Z", Source: "mock"},
		{Status: "DELIVERED", Description: "Доставлена получателю", Location: "Астана, Казахстан", Date: "2026-02-15T12:00:00Z", Source: "mock"},
	}

	JSON(w, http.StatusOK, map[string]interface{}{
		"track_number": track,
		"parcel":       results[0].Parcel,
		"events":       events,
	})
}
