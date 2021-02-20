package model

import (
	"database/sql"

	"github.com/golang/protobuf/ptypes/timestamp"
	"onepass.app/facility/hts/common"
)

// Facility is model for database
type Facility struct {
	ID             int64
	OrganizationID int64
	Name           string
	Latitude       float64
	Longitude      float64
	OperatingHours string
	Description    string
}

// FacilityRequest is model for database
type FacilityRequest struct {
	ID           int64
	EventID      int64
	FacilityID   int64
	Status       common.Status
	RejectReason sql.NullString
	Start        *timestamp.Timestamp
	Finish       *timestamp.Timestamp
}
