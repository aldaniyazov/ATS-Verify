package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"

	"ats-verify/internal/models"
)

// RiskRepository handles IIN/BIN risk profile operations.
type RiskRepository struct {
	db *sql.DB
}

// NewRiskRepository creates a new RiskRepository.
func NewRiskRepository(db *sql.DB) *RiskRepository {
	return &RiskRepository{db: db}
}

// Upsert creates or updates a risk profile for an IIN/BIN.
func (r *RiskRepository) Upsert(ctx context.Context, profile *models.RiskProfile) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO risk_profiles (id, iin_bin, risk_level, flagged_by, reason, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		 ON CONFLICT (iin_bin)
		 DO UPDATE SET risk_level = EXCLUDED.risk_level, flagged_by = EXCLUDED.flagged_by, reason = EXCLUDED.reason, updated_at = NOW()`,
		uuid.New(), profile.IINBIN, profile.RiskLevel, profile.FlaggedBy, profile.Reason,
	)
	if err != nil {
		return fmt.Errorf("upserting risk profile: %w", err)
	}
	return nil
}

// GetByIINBIN retrieves a risk profile by IIN/BIN.
func (r *RiskRepository) GetByIINBIN(ctx context.Context, iinBin string) (*models.RiskProfile, error) {
	var p models.RiskProfile
	err := r.db.QueryRowContext(ctx,
		`SELECT id, iin_bin, risk_level, flagged_by, reason, created_at, updated_at
		 FROM risk_profiles WHERE iin_bin = $1`,
		iinBin,
	).Scan(&p.ID, &p.IINBIN, &p.RiskLevel, &p.FlaggedBy, &p.Reason, &p.CreatedAt, &p.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("querying risk profile: %w", err)
	}
	return &p, nil
}

// ListAll returns all risk profiles.
func (r *RiskRepository) ListAll(ctx context.Context) ([]models.RiskProfile, error) {
	rows, err := r.db.QueryContext(ctx,
		"SELECT id, iin_bin, risk_level, flagged_by, reason, created_at, updated_at FROM risk_profiles ORDER BY updated_at DESC",
	)
	if err != nil {
		return nil, fmt.Errorf("listing risk profiles: %w", err)
	}
	defer rows.Close()

	var profiles []models.RiskProfile
	for rows.Next() {
		var p models.RiskProfile
		if err := rows.Scan(&p.ID, &p.IINBIN, &p.RiskLevel, &p.FlaggedBy, &p.Reason, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scanning risk profile: %w", err)
		}
		profiles = append(profiles, p)
	}
	return profiles, nil
}

// Delete removes a risk profile by ID.
func (r *RiskRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM risk_profiles WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("deleting risk profile: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("risk profile not found")
	}
	return nil
}
