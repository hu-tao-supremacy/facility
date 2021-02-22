package model

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx/types"
)

// OperatingHour is struct for facility operationHours
type OperatingHour struct {
	StartHour  int64  `json:"startHour"`
	FinishHour int64  `json:"finishHour"`
	Day        string `json:"day"`
}

// Facility is model for database
type Facility struct {
	ID             int64
	OrganizationID int64
	Name           string
	Latitude       float64
	Longitude      float64
	OperatingHours types.JSONText
	Description    string
}

// FacilityRequest is model for database
type FacilityRequest struct {
	ID           int64
	EventID      int64
	FacilityID   int64
	Status       string
	RejectReason sql.NullString
	Start        time.Time
	Finish       time.Time
}

// FacilityRequestWithInfo is joint model between Facility and FacilityRequest for database
type FacilityRequestWithInfo struct {
	ID             int64
	EventID        int64
	FacilityID     int64
	Status         string
	RejectReason   sql.NullString
	Start          time.Time
	Finish         time.Time
	FaciltiyID     int64
	OrganizationID int64
	FacilityName   string
	Latitude       float64
	Longitude      float64
	OperatingHours types.JSONText
	Description    string
}
