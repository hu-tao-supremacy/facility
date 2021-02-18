package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	facility "onepass.app/facility/hts/facility"

	_ "github.com/go-sql-driver/mysql"
)

type FacilityServer struct {
	facility.UnimplementedFacilityServiceServer
}




func (fs *FacilityServer)GetFacilityList(ctx context.Context, in *facility.GetFacilityListRequest) (*facility.GetFacilityListResponse, error) {
	list := make([]*facility.Facility, 1)
	fmt.Println("emp:", list)
	return &facility.GetFacilityListResponse{
		Facilities: list,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	facility.RegisterFacilityServiceServer(s, &FacilityServer{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
