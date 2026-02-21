package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"ats-verify/internal/models"
)

// ParcelRepository handles parcel database operations.
type ParcelRepository struct {
	db *sql.DB
}

// NewParcelRepository creates a new ParcelRepository.
func NewParcelRepository(db *sql.DB) *ParcelRepository {
	return &ParcelRepository{db: db}
}

// ParcelUpsertResult describes what happened during an upsert attempt.
type ParcelUpsertResult struct {
	TrackNumber string
	Action      string // "inserted", "updated", "skipped_used", "error"
	Message     string
}

// UpsertParcel implements the deduplication logic from GOALS.md:
// 1. New Track Number -> INSERT (is_used = false)
// 2. Existing Track (is_used=false) -> UPDATE (overwrite data)
// 3. Existing Track (is_used=true) -> ERROR ("Track already used")
func (r *ParcelRepository) UpsertParcel(ctx context.Context, p *models.Parcel) (*ParcelUpsertResult, error) {
	// Step 1: Check if track number exists
	var existingID uuid.UUID
	var isUsed bool

	err := r.db.QueryRowContext(ctx,
		"SELECT id, is_used FROM parcels WHERE track_number = $1",
		p.TrackNumber,
	).Scan(&existingID, &isUsed)

	if err == sql.ErrNoRows {
		// Case 1: New track number -> INSERT
		_, err := r.db.ExecContext(ctx,
			`INSERT INTO parcels (id, track_number, marketplace, country, brand, product_name, snt, is_used, upload_date, uploaded_by, created_at, updated_at)
			 VALUES ($1, $2, $3, $4, $5, $6, $7, false, $8, $9, NOW(), NOW())`,
			uuid.New(), p.TrackNumber, p.Marketplace, p.Country, p.Brand, p.ProductName, p.SNT, p.UploadDate, p.UploadedBy,
		)
		if err != nil {
			return nil, fmt.Errorf("inserting parcel: %w", err)
		}
		return &ParcelUpsertResult{
			TrackNumber: p.TrackNumber,
			Action:      "inserted",
			Message:     "New parcel created",
		}, nil
	}

	if err != nil {
		return nil, fmt.Errorf("checking existing parcel: %w", err)
	}

	// Case 3: Track exists and is_used=true -> REJECT
	if isUsed {
		return &ParcelUpsertResult{
			TrackNumber: p.TrackNumber,
			Action:      "skipped_used",
			Message:     "Track already used (is_used=true). Cannot overwrite.",
		}, nil
	}

	// Case 2: Track exists and is_used=false -> UPDATE
	_, err = r.db.ExecContext(ctx,
		`UPDATE parcels SET marketplace=$1, country=$2, brand=$3, product_name=$4, snt=$5, upload_date=$6, uploaded_by=$7, updated_at=NOW()
		 WHERE id=$8`,
		p.Marketplace, p.Country, p.Brand, p.ProductName, p.SNT, p.UploadDate, p.UploadedBy, existingID,
	)
	if err != nil {
		return nil, fmt.Errorf("updating parcel: %w", err)
	}

	return &ParcelUpsertResult{
		TrackNumber: p.TrackNumber,
		Action:      "updated",
		Message:     "Existing parcel updated (was not used)",
	}, nil
}

// GetByTrackNumber retrieves a parcel by its track number.
func (r *ParcelRepository) GetByTrackNumber(ctx context.Context, trackNumber string) (*models.Parcel, error) {
	var p models.Parcel
	err := r.db.QueryRowContext(ctx,
		`SELECT id, track_number, marketplace, country, brand, product_name, snt, is_used, upload_date, uploaded_by, created_at, updated_at
		 FROM parcels WHERE track_number = $1`,
		trackNumber,
	).Scan(&p.ID, &p.TrackNumber, &p.Marketplace, &p.Country, &p.Brand, &p.ProductName, &p.SNT, &p.IsUsed, &p.UploadDate, &p.UploadedBy, &p.CreatedAt, &p.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("querying parcel by track: %w", err)
	}
	return &p, nil
}

// MarkUsed sets is_used=true for a given parcel.
func (r *ParcelRepository) MarkUsed(ctx context.Context, trackNumber string) error {
	trackNumber = strings.TrimSpace(trackNumber)
	result, err := r.db.ExecContext(ctx,
		"UPDATE parcels SET is_used = true, updated_at = NOW() WHERE track_number = $1",
		trackNumber,
	)
	if err != nil {
		return fmt.Errorf("marking parcel as used: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("parcel with track_number %s not found (trimmed)", trackNumber)
	}
	return nil
}

// ListAll returns all parcels, optionally filtered by period.
func (r *ParcelRepository) ListAll(ctx context.Context, from, to *time.Time) ([]models.Parcel, error) {
	query := "SELECT id, track_number, marketplace, country, brand, product_name, snt, is_used, upload_date, uploaded_by, created_at, updated_at FROM parcels"
	args := []interface{}{}

	if from != nil && to != nil {
		query += " WHERE created_at >= $1 AND created_at <= $2"
		args = append(args, from, to)
	}
	query += " ORDER BY created_at DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("listing parcels: %w", err)
	}
	defer rows.Close()

	var parcels []models.Parcel
	for rows.Next() {
		var p models.Parcel
		if err := rows.Scan(&p.ID, &p.TrackNumber, &p.Marketplace, &p.Country, &p.Brand, &p.ProductName, &p.SNT, &p.IsUsed, &p.UploadDate, &p.UploadedBy, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scanning parcel row: %w", err)
		}
		parcels = append(parcels, p)
	}
	return parcels, nil
}

// ListWithFilters returns parcels with filtering, search, and pagination.
func (r *ParcelRepository) ListWithFilters(ctx context.Context, status, search string, page, limit int) ([]models.Parcel, int, error) {
	where := []string{}
	args := []interface{}{}
	argIdx := 1

	if status == "used" {
		where = append(where, fmt.Sprintf("is_used = $%d", argIdx))
		args = append(args, true)
		argIdx++
	} else if status == "unused" {
		where = append(where, fmt.Sprintf("is_used = $%d", argIdx))
		args = append(args, false)
		argIdx++
	}

	if search != "" {
		where = append(where, fmt.Sprintf("(track_number ILIKE $%d OR product_name ILIKE $%d OR brand ILIKE $%d)", argIdx, argIdx, argIdx))
		args = append(args, "%"+search+"%")
		argIdx++
	}

	whereClause := ""
	if len(where) > 0 {
		whereClause = " WHERE " + strings.Join(where, " AND ")
	}

	// Count total
	var total int
	countQuery := "SELECT COUNT(*) FROM parcels" + whereClause
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("counting parcels: %w", err)
	}

	// Fetch page
	offset := (page - 1) * limit
	dataQuery := "SELECT id, track_number, marketplace, country, brand, product_name, snt, is_used, upload_date, uploaded_by, created_at, updated_at FROM parcels" +
		whereClause + fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("listing parcels with filters: %w", err)
	}
	defer rows.Close()

	var parcels []models.Parcel
	for rows.Next() {
		var p models.Parcel
		if err := rows.Scan(&p.ID, &p.TrackNumber, &p.Marketplace, &p.Country, &p.Brand, &p.ProductName, &p.SNT, &p.IsUsed, &p.UploadDate, &p.UploadedBy, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("scanning parcel row: %w", err)
		}
		parcels = append(parcels, p)
	}
	return parcels, total, nil
}

// BulkLookup retrieves parcels for a slice of track numbers.
func (r *ParcelRepository) BulkLookup(ctx context.Context, tracks []string) ([]models.Parcel, error) {
	if len(tracks) == 0 {
		return nil, nil
	}

	// Build IN clause
	placeholders := make([]string, len(tracks))
	args := make([]interface{}, len(tracks))
	for i, t := range tracks {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = strings.TrimSpace(t)
	}

	query := "SELECT id, track_number, marketplace, country, brand, product_name, snt, is_used, upload_date, uploaded_by, created_at, updated_at FROM parcels WHERE track_number IN (" +
		strings.Join(placeholders, ",") + ")"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("bulk lookup parcels: %w", err)
	}
	defer rows.Close()

	var parcels []models.Parcel
	for rows.Next() {
		var p models.Parcel
		if err := rows.Scan(&p.ID, &p.TrackNumber, &p.Marketplace, &p.Country, &p.Brand, &p.ProductName, &p.SNT, &p.IsUsed, &p.UploadDate, &p.UploadedBy, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scanning bulk parcel: %w", err)
		}
		parcels = append(parcels, p)
	}
	return parcels, nil
}
