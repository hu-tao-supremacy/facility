package main

import (
	"context"
	"fmt"
	"log"
	"net"

	empty "github.com/golang/protobuf/ptypes/empty"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"onepass.app/facility/hts/common"
	facility "onepass.app/facility/hts/facility"
	database "onepass.app/facility/internal/database"
	typing "onepass.app/facility/internal/typing"

	_ "github.com/lib/pq"
)

// FacilityServer is for handling facility endpoint
type FacilityServer struct {
	facility.UnimplementedFacilityServiceServer
	dbs *database.DataService
}

// GetFacilityList is a function to list all facilities owned by organization
func (fs *FacilityServer) GetFacilityList(ctx context.Context, in *facility.GetFacilityListRequest) (*facility.GetFacilityListResponse, error) {
	list, err := fs.dbs.GetFacilityList(in.OrganizationId)

	if err != nil {
		return nil, status.Error(err.StatusCode, err.Error())
	}

	return &facility.GetFacilityListResponse{
		Facilities: list,
	}, nil
}

// GetAvailableFacilityList is a function to list all available facilities
func (fs *FacilityServer) GetAvailableFacilityList(ctx context.Context, in *empty.Empty) (*facility.GetAvailableFacilityListResponse, error) {
	list, err := fs.dbs.GetAvailableFacilityList()

	if err != nil {
		return nil, status.Error(err.StatusCode, err.Error())
	}

	return &facility.GetAvailableFacilityListResponse{
		Facilities: list,
	}, nil
}

// GetFacilityInfo is a function to get facility’s information
func (fs *FacilityServer) GetFacilityInfo(ctx context.Context, in *facility.GetFacilityInfoRequest) (*common.Facility, error) {
	result, err := fs.dbs.GetFacilityInfo(in.FacilityId)

	if err != nil {
		return nil, status.Error(err.StatusCode, err.Error())
	}

	return result, nil
}

// ApproveFacilityRequest is a function to reject facility’s request by id
func (fs *FacilityServer) ApproveFacilityRequest(ctx context.Context, in *facility.ApproveFacilityRequestRequest) (*common.Result, error) {
	permission := common.Permission_UPDATE_FACILITY
	isAbleToApproveRequest := hasPermission(in.UserId, in.OrganizationId, permission)
	if !isAbleToApproveRequest {
		return nil, status.Error(codes.PermissionDenied, (&typing.PermissionError{Type: permission}).Error())
	}

	facilityRequest, err := fs.dbs.GetFacilityRequest(in.RequestId)
	if err != nil {
		return nil, status.Error(err.StatusCode, err.Error())
	}

	isTimeOverlap, err := fs.dbs.IsOverLapTime(facilityRequest.FacilityId, facilityRequest.Start, facilityRequest.Finish)
	if err != nil {
		return nil, status.Error(err.StatusCode, err.Error())
	}
	if isTimeOverlap {
		return nil, status.Error(codes.AlreadyExists, (&typing.AlreadyExistError{Name: "Facility is booked at that time"}).Error())
	}

	err = fs.dbs.ApproveFacilityRequest(in.RequestId)
	if err != nil {
		return nil, status.Error(err.StatusCode, err.Error())
	}

	description := fmt.Sprintf("Request ID: %d has been aproved", in.RequestId)
	return &common.Result{
		IsOk:        true,
		Description: description,
	}, nil
}

// RejectFacilityRequest is a function to reject facility’s request by id
func (fs *FacilityServer) RejectFacilityRequest(ctx context.Context, in *facility.RejectFacilityRequestRequest) (*common.Result, error) {
	permission := common.Permission_UPDATE_FACILITY
	isAbleToRejectRequest := hasPermission(in.UserId, in.OrganizationId, permission)
	if !isAbleToRejectRequest {
		return nil, status.Error(codes.PermissionDenied, (&typing.PermissionError{Type: permission}).Error())
	}

	err := fs.dbs.RejectFacilityRequest(in.RequestId, in.Reason)

	if err != nil {
		return nil, status.Error(err.StatusCode, err.Error())
	}

	description := fmt.Sprintf("Request ID: %d has been rejected", in.RequestId)
	return &common.Result{
		IsOk:        true,
		Description: description,
	}, nil
}

// CreateFacilityRequest is a function to create facility’s request by id
func (fs *FacilityServer) CreateFacilityRequest(ctx context.Context, in *facility.CreateFacilityRequestRequest) (*common.FacilityRequest, error) {
	permssion := common.Permission_UPDATE_EVENT
	havingPermissionChannel := make(chan bool)
	eventOwnerChannel := make(chan bool)
	overlapTimeChannel := make(chan bool)
	errorChannel := make(chan *typing.DatabaseError)
	go func() { havingPermissionChannel <- hasPermission(in.UserId, in.OrganizationId, permssion) }()
	go func() { eventOwnerChannel <- hasEvent(in.UserId, in.OrganizationId, in.EventId) }()
	go func() {
		isTimeOverlap, err := fs.dbs.IsOverLapTime(in.FacilityId, in.Start, in.End)
		overlapTimeChannel <- isTimeOverlap
		errorChannel <- err
	}()

	isPermission := <-havingPermissionChannel
	isEventOwner := <-eventOwnerChannel
	isTimeOverlap := <-overlapTimeChannel
	overalpErr := <-errorChannel

	close(havingPermissionChannel)
	close(eventOwnerChannel)
	close(overlapTimeChannel)
	close(errorChannel)

	if !(isPermission && isEventOwner) {
		return nil, status.Error(codes.PermissionDenied, (&typing.PermissionError{Type: permssion}).Error())
	}

	if overalpErr != nil {
		return nil, status.Error(codes.Internal, overalpErr.Error())
	}
	if isTimeOverlap {
		return nil, status.Error(codes.AlreadyExists, (&typing.AlreadyExistError{Name: "Facility is booked at that time"}).Error())
	}

	result, err := fs.dbs.CreateFacilityRequest(in.EventId, in.FacilityId, in.Start, in.End)

	if err != nil {
		return nil, status.Error(err.StatusCode, err.Error())
	}

	return result, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()

	facilityServer := &FacilityServer{}
	db := &database.DataService{}
	db.ConnectToDB()
	facilityServer.dbs = db

	facility.RegisterFacilityServiceServer(s, facilityServer)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
