package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/golang/protobuf/ptypes"
	"github.com/jmoiron/sqlx/reflectx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	common "onepass.app/facility/hts/common"
	model "onepass.app/facility/internal/model"
	typing "onepass.app/facility/internal/typing"

	"github.com/iancoleman/strcase"
	"github.com/jmoiron/sqlx"
)

// DataService is for handling data layer
type DataService struct {
	SQL *sqlx.DB
}

// GetFacilityList is a function to get facility list owned by the organization from database
func (dbs *DataService) GetFacilityList(organizationID int64) ([]*common.Facility, typing.CustomError) {
	var facilities []*model.Facility
	query := fmt.Sprintf(`
	SELECT * 
	FROM facility 
	WHERE facility.organization_id = %d;`,
		organizationID)
	err := dbs.SQL.Select(&facilities, query)

	if err != nil {
		return nil, &typing.DatabaseError{
			Err:        err,
			StatusCode: codes.Internal,
		}
	}

	result := make([]*common.Facility, len(facilities))
	for i, item := range facilities {
		result[i] = convertFacilityModelToProto(item)
	}

	return result, nil
}

// GetAvailableFacilityList is a function to list all available facilities
func (dbs *DataService) GetAvailableFacilityList() ([]*common.Facility, typing.CustomError) {
	var facilities []*model.Facility
	query := `
	SELECT * 
	FROM facility`
	err := dbs.SQL.Select(&facilities, query)

	if err != nil {
		return nil, &typing.DatabaseError{
			Err:        err,
			StatusCode: codes.Internal,
		}
	}

	result := make([]*common.Facility, len(facilities))
	for i, item := range facilities {
		result[i] = convertFacilityModelToProto(item)
	}

	return result, nil
}

// GetFacilityInfo is a function to get facility’s information by id
func (dbs *DataService) GetFacilityInfo(facilityID int64) (*common.Facility, typing.CustomError) {
	var _facility model.Facility
	query := fmt.Sprintf(`
	SELECT * 
	FROM facility 
	WHERE facility.id = %d`,
		facilityID)
	err := dbs.SQL.Get(&_facility, query)

	switch {
	case err == sql.ErrNoRows:
		return nil, &typing.DatabaseError{
			Err:        &typing.NotFoundError{Name: "facility"},
			StatusCode: codes.NotFound,
		}
	case err != nil:
		return nil, &typing.DatabaseError{
			Err:        err,
			StatusCode: codes.Internal,
		}
	default:
		return convertFacilityModelToProto(&_facility), nil
	}
}

func (dbs *DataService) updateFacilityRequest(requestID int64, status common.Status, reason *wrapperspb.StringValue) typing.CustomError {
	var queryReason string
	if reason != nil {
		queryReason = ", reject_reason=:reason "
	}

	query := fmt.Sprintf(`
	UPDATE facility_request 
	SET status=:status%s 
	WHERE facility_request.id = :id`,
		queryReason)
	result, err := dbs.SQL.NamedExec(query, map[string]interface{}{
		"id":     requestID,
		"status": status.String(),
		"reason": reason.GetValue(),
	})
	if err != nil {
		return &typing.DatabaseError{
			Err:        err,
			StatusCode: codes.Internal,
		}
	}

	count, err := result.RowsAffected()
	switch {
	case err != nil:
		return &typing.DatabaseError{
			Err:        err,
			StatusCode: codes.Internal,
		}
	case count != 1:
		return &typing.DatabaseError{
			Err:        &typing.NotFoundError{Name: "FacilityRequest"},
			StatusCode: codes.NotFound,
		}
	default:
		return nil
	}
}

// RejectFacilityRequest is a function to reject facility’s request by id
func (dbs *DataService) RejectFacilityRequest(requestID int64, reason *wrapperspb.StringValue) typing.CustomError {
	return dbs.updateFacilityRequest(requestID, common.Status_REJECTED, reason)
}

// ApproveFacilityRequest is a function to approve facility request
func (dbs *DataService) ApproveFacilityRequest(requestID int64) typing.CustomError {
	return dbs.updateFacilityRequest(requestID, common.Status_APPROVED, nil)
}

// CreateFacilityRequest is a function to create facilityRequest
func (dbs *DataService) CreateFacilityRequest(eventID int64, facilityID int64, start *timestamppb.Timestamp, finish *timestamppb.Timestamp) (*common.FacilityRequest, typing.CustomError) {
	var id int64
	query := `
	INSERT INTO facility_request (event_id, facility_id, status, start, finish) 
	VALUES (:event_id, :facility_id, :status, :start, :finish) 
	RETURNING id`
	startTime, _ := ptypes.Timestamp(start)
	finishTime, _ := ptypes.Timestamp(finish)
	rows, err := dbs.SQL.NamedQuery(query, map[string]interface{}{
		"event_id":    eventID,
		"facility_id": facilityID,
		"status":      "PENDING",
		"start":       startTime,
		"finish":      finishTime,
	})
	if rows.Next() {
		rows.Scan(&id)
	}

	if err != nil {
		return nil, &typing.DatabaseError{
			Err:        err,
			StatusCode: codes.Internal,
		}
	}

	result := common.FacilityRequest{
		Id:         id,
		EventId:    eventID,
		FacilityId: facilityID,
		Status:     common.Status_PENDING,
		Start:      start,
		Finish:     finish,
	}
	return &result, nil
}

// IsOverlapTime is function to check whether time is overlap with already booked facility
func (dbs *DataService) IsOverlapTime(facilityID int64, start *timestamppb.Timestamp, finish *timestamppb.Timestamp) (bool, typing.CustomError) {
	_, facilityNotFoundError := dbs.GetFacilityInfo(facilityID)
	if facilityNotFoundError != nil {
		return false, facilityNotFoundError
	}

	var count int64
	startTime, _ := ptypes.Timestamp(start)
	finishTime, _ := ptypes.Timestamp(finish)

	layoutTime := "2006-01-02 15:04:05"
	startTimeText := startTime.Format(layoutTime)
	finishTimeText := finishTime.Format(layoutTime)

	query := fmt.Sprintf(`
	SELECT COUNT(*) 
	FROM facility_request 
	WHERE (('%s' >= start AND '%s' < finish) OR ('%s' > start AND '%s' <= finish)) 
	AND facility_id = %d 
	AND status='APPROVED' 
	LIMIT 1;`,
		startTimeText, startTimeText, finishTimeText, finishTimeText, facilityID)
	err := dbs.SQL.Get(&count, query)

	if err != nil {
		return false, &typing.DatabaseError{
			Err:        err,
			StatusCode: codes.Internal,
		}
	}

	return count != 0, nil
}

// GetFacilityRequest is function to get facility request by id
func (dbs *DataService) GetFacilityRequest(RequestID int64) (*common.FacilityRequest, typing.CustomError) {
	var facilityRequest model.FacilityRequest

	query := fmt.Sprintf(`
	SELECT * 
	FROM facility_request 
	WHERE id=%d 
	LIMIT 1
	`, RequestID)
	err := dbs.SQL.Get(&facilityRequest, query)

	switch {
	case err == sql.ErrNoRows:
		return nil, &typing.DatabaseError{
			Err:        &typing.NotFoundError{Name: "facility"},
			StatusCode: codes.NotFound,
		}
	case err != nil:
		return nil, &typing.DatabaseError{
			Err:        err,
			StatusCode: codes.Internal,
		}
	default:
		return convertFacilityRequestModelToProto(&facilityRequest), nil
	}
}

func (dbs *DataService) ping() (string, error) {
	var version string
	err := dbs.SQL.Get(&version, "SELECT VERSION();")

	if err != nil {
		return version, status.Error(codes.Internal, err.Error())
	}

	return version, nil
}

// ConnectToDB is a function to connect to DB and setup sqlx config
func (dbs *DataService) ConnectToDB() {
	host := os.Getenv("POSTGRES_HOST")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	database := os.Getenv("POSTGRES_DB")
	port := os.Getenv("POSTGRES_PORT")

	dsn := fmt.Sprintf("user=%s password=%s host=%s database=%s port=%s sslmode=disable", user, password, host, database, port)
	db, err := sqlx.Connect("postgres", dsn)

	if err != nil {
		log.Fatalln(err)
	}

	strcase.ConfigureAcronym("ID", "id")
	db.Mapper = reflectx.NewMapperFunc("json", strcase.ToSnake)
	dbs.SQL = db
	version, err := dbs.ping()
	if err == nil {
		log.Println("SQL version:", version)
	}
}
