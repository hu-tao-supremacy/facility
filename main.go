package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	facility "onepass.app/facility/hts/facility"

	_ "github.com/lib/pq"
)

// DataService is for handling data layer
type DataService struct {
	sql *sqlx.DB
}

// FacilityServer is for handling facility endpoint 
type FacilityServer struct {
	facility.UnimplementedFacilityServiceServer
	dataService *DataService
}


// GetFacilityList is a function to list all facilities owned by organization
func (fs *FacilityServer) GetFacilityList(ctx context.Context, in *facility.GetFacilityListRequest) (*facility.GetFacilityListResponse, error) {
	list := make([]*facility.Facility, 1)
	fmt.Println("emp:", list)

	return &facility.GetFacilityListResponse{
		Facilities: list,
	}, nil
}

func (dbs *DataService) connectToDB() {
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

	dbs.sql = db
}

func (fs *FacilityServer) ping() (string, error) {
	var version string
	err := fs.dataService.sql.Get(&version, "SELECT VERSION();")

	if err != nil {
		return version, status.Error(codes.Internal, err.Error())
	}

	return version, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	
	facilityServer := &FacilityServer{}
	db := &DataService{}
	db.connectToDB()
	facilityServer.dataService = db
	
	version, err := facilityServer.ping()
	if err == nil {
		log.Println("SQL version:", version)
	}

	facility.RegisterFacilityServiceServer(s, facilityServer)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}
