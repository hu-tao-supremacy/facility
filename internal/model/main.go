package model

import (
	"database/sql"

	"github.com/golang/protobuf/ptypes/timestamp"
	"onepass.app/facility/hts/facility"
)

// Facility is model for database
type Facility struct {
	ID             int64
	OrganizationID int64
	Name           string
	Latitude       int64
	Longitude      int64
	OperatingHours sql.NullString
	Description    sql.NullString
}

// FacilityRequest is model for database
type FacilityRequest struct {
	ID           int64
	EventID      int64
	FacilityID   int64
	Status       facility.RequestStatus
	RejectReason sql.NullString
	Start        *timestamp.Timestamp
	Finish       *timestamp.Timestamp
}
