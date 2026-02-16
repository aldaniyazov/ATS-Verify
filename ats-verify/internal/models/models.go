package models

import (
	"time"

	"github.com/google/uuid"
)

// UserRole defines the type for user roles.
type UserRole string

const (
	RoleAdmin       UserRole = "admin"
	RolePaidUser    UserRole = "paid_user"
	RoleATSStaff    UserRole = "ats_staff"
	RoleCustoms     UserRole = "customs_staff"
	RoleMarketplace UserRole = "marketplace_staff"
)

// RiskLevel defines the type for IIN/BIN risk assessment.
type RiskLevel string

const (
	RiskGreen  RiskLevel = "green"
	RiskYellow RiskLevel = "yellow"
	RiskRed    RiskLevel = "red"
)

// User represents a system user.
type User struct {
	ID                uuid.UUID `json:"id" db:"id"`
	Username          string    `json:"username" db:"username"`
	PasswordHash      string    `json:"-" db:"password_hash"`
	Role              UserRole  `json:"role" db:"role"`
	MarketplacePrefix *string   `json:"marketplace_prefix,omitempty" db:"marketplace_prefix"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
}

// Parcel represents a tracked parcel in the system.
type Parcel struct {
	ID          uuid.UUID `json:"id" db:"id"`
	TrackNumber string    `json:"track_number" db:"track_number"`
	Marketplace string    `json:"marketplace" db:"marketplace"`
	Country     string    `json:"country,omitempty" db:"country"`
	Brand       string    `json:"brand,omitempty" db:"brand"`
	ProductName string    `json:"product_name,omitempty" db:"product_name"`
	SNT         string    `json:"snt,omitempty" db:"snt"`
	IsUsed      bool      `json:"is_used" db:"is_used"`
	UploadDate  time.Time `json:"upload_date" db:"upload_date"`
	UploadedBy  uuid.UUID `json:"uploaded_by" db:"uploaded_by"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// RiskProfile represents a risk assessment for an IIN/BIN.
type RiskProfile struct {
	ID        uuid.UUID `json:"id" db:"id"`
	IINBIN    string    `json:"iin_bin" db:"iin_bin"`
	RiskLevel RiskLevel `json:"risk_level" db:"risk_level"`
	FlaggedBy uuid.UUID `json:"flagged_by" db:"flagged_by"`
	Reason    string    `json:"reason,omitempty" db:"reason"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// TrackingEvent represents a single tracking event from Kazpost/CDEK.
type TrackingEvent struct {
	ID          uuid.UUID `json:"id" db:"id"`
	ParcelID    uuid.UUID `json:"parcel_id" db:"parcel_id"`
	StatusCode  string    `json:"status_code" db:"status_code"`
	Description string    `json:"description" db:"description"`
	Location    string    `json:"location" db:"location"`
	EventTime   time.Time `json:"event_time" db:"event_time"`
	Source      string    `json:"source" db:"source"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// AnalysisReport stores results of IMEI verification or risk analysis.
type AnalysisReport struct {
	ID            uuid.UUID              `json:"id" db:"id"`
	UserID        uuid.UUID              `json:"user_id" db:"user_id"`
	ReportType    string                 `json:"report_type" db:"report_type"`
	InputFileName string                 `json:"input_file_name" db:"input_file_name"`
	ResultSummary map[string]interface{} `json:"result_summary" db:"result_summary"`
	RawDataURL    string                 `json:"raw_data_url,omitempty" db:"raw_data_url"`
	CreatedAt     time.Time              `json:"created_at" db:"created_at"`
}

// MarketplacePrefixMap maps user role suffixes to marketplace names.
var MarketplacePrefixMap = map[string]string{
	"wb":    "Wildberries",
	"ozon":  "Ozon",
	"kaspi": "Kaspi",
	"ali":   "AliExpress",
	"temu":  "Temu",
}
