package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/google/uuid"

	"ats-verify/internal/models"
)

// Seed creates initial users if the users table is empty.
func Seed(ctx context.Context, db *sql.DB, hashFn func(string) (string, error)) error {
	var count int
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM users").Scan(&count); err != nil {
		return fmt.Errorf("checking users count: %w", err)
	}

	if count > 0 {
		log.Println("seed: users table not empty, skipping")
		return nil
	}

	log.Println("seed: creating initial users...")

	wbPrefix := "wb"
	users := []struct {
		username string
		password string
		role     models.UserRole
		prefix   *string
	}{
		{"admin", "admin", models.RoleAdmin, nil},
		{"marketplace_wb", "marketplace_wb", models.RoleMarketplace, &wbPrefix},
		{"ats_staff", "ats_staff", models.RoleATSStaff, nil},
	}

	for _, u := range users {
		hash, err := hashFn(u.password)
		if err != nil {
			return fmt.Errorf("hashing password for %s: %w", u.username, err)
		}

		_, err = db.ExecContext(ctx,
			`INSERT INTO users (id, username, password_hash, role, marketplace_prefix, is_approved, created_at, updated_at)
			 VALUES ($1, $2, $3, $4, $5, true, NOW(), NOW())`,
			uuid.New(), u.username, hash, u.role, u.prefix,
		)
		if err != nil {
			return fmt.Errorf("inserting seed user %s: %w", u.username, err)
		}
		log.Printf("seed: created user %s (role: %s)", u.username, u.role)
	}

	return nil
}
