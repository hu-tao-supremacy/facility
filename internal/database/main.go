package database

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jmoiron/sqlx/reflectx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	facility "onepass.app/facility/hts/facility"

	"github.com/jmoiron/sqlx"
)

// DataService is for handling data layer
type DataService struct {
	SQL *sqlx.DB
}

// GetFacilityList is a function to get facility list owned by the organization from database
func (dbs *DataService) GetFacilityList(organizationID int64) ([]*facility.Facility, error) {
	var facilities []*facility.Facility
	query := fmt.Sprintf("SELECT * FROM facility WHERE facility.organization_id = %d;", organizationID)
	err := dbs.SQL.Select(&facilities, query)

	if err != nil {
		return nil, err
	}

	return facilities, nil
}

// GetAvailableFacilityList is a function to list all available facilities
func (dbs *DataService) GetAvailableFacilityList() ([]*facility.Facility, error) {
	var facilities []*facility.Facility
	query := "SELECT * FROM facility"
	err := dbs.SQL.Select(&facilities, query)

	if err != nil {
		return nil, err
	}

	return facilities, nil
}

// GetFacilityInfo is a function to get facility’s information by id
func (dbs *DataService) GetFacilityInfo(facilityID int64) (*facility.Facility, error) {
	var _facility facility.Facility
	query := fmt.Sprintf("SELECT * FROM facility WHERE facility.id = %d", facilityID)
	err := dbs.SQL.Get(&_facility, query)

	if err != nil {
		return nil, err
	}

	return &_facility, nil
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

	db.Mapper = reflectx.NewMapperFunc("json", strings.ToLower)
	dbs.SQL = db
	version, err := dbs.ping()
	if err == nil {
		log.Println("SQL version:", version)
	}
}
