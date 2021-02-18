package main

import (
	"context"
	"fmt"
	"log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net"
	"os"
	"google.golang.org/grpc"
	"github.com/jmoiron/sqlx"
	facility "onepass.app/facility/hts/facility"
	
	_ "github.com/lib/pq"
)

// DataService is very gppd
type DataService struct {
	sql *sqlx.DB
}

// FacilityServer is very gppd
type FacilityServer struct {
	facility.UnimplementedFacilityServiceServer
	dataService *DataService
}


// GetFacilityList is a function that is VERY GOOD
func (fs *FacilityServer) GetFacilityList(ctx context.Context, in *facility.GetFacilityListRequest) (*facility.GetFacilityListResponse, error) {
	list := make([]*facility.Facility, 1)
	fmt.Println("emp:", list)

	return &facility.GetFacilityListResponse{
		Facilities: list,
	}, nil
}

// ConnectToDB is very good
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
	db.ConnectToDB()
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
