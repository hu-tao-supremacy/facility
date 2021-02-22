package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/golang/protobuf/ptypes"
	empty "github.com/golang/protobuf/ptypes/empty"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"onepass.app/facility/hts/common"
	facility "onepass.app/facility/hts/facility"
	database "onepass.app/facility/internal/database"
	"onepass.app/facility/internal/helper"
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
		return nil, status.Error(err.Code(), err.Error())
	}

	return &facility.GetFacilityListResponse{
		Facilities: list,
	}, nil
}

// GetAvailableFacilityList is a function to list all available facilities
func (fs *FacilityServer) GetAvailableFacilityList(ctx context.Context, in *empty.Empty) (*facility.GetAvailableFacilityListResponse, error) {
	list, err := fs.dbs.GetAvailableFacilityList()

	if err != nil {
		return nil, status.Error(err.Code(), err.Error())
	}

	return &facility.GetAvailableFacilityListResponse{
		Facilities: list,
	}, nil
}

// GetFacilityInfo is a function to get facility’s information
func (fs *FacilityServer) GetFacilityInfo(ctx context.Context, in *facility.GetFacilityInfoRequest) (*common.Facility, error) {
	result, err := fs.dbs.GetFacilityInfo(in.FacilityId)

	if err != nil {
		return nil, status.Error(err.Code(), err.Error())
	}

	return result, nil
}

// ApproveFacilityRequest is a function to reject facility’s request by id
func (fs *FacilityServer) ApproveFacilityRequest(ctx context.Context, in *facility.ApproveFacilityRequestRequest) (*common.Result, error) {
	isConditionPassed, err := isAbleToApproveFacilityRequest(fs, in)

	if !isConditionPassed || err != nil {
		return nil, status.Error(err.Code(), err.Error())
	}

	err = fs.dbs.ApproveFacilityRequest(in.RequestId)
	if err != nil {
		return nil, status.Error(err.Code(), err.Error())
	}

	description := fmt.Sprintf("Request ID: %d has been aproved", in.RequestId)
	return &common.Result{
		IsOk:        true,
		Description: description,
	}, nil
}

// RejectFacilityRequest is a function to reject facility’s request by id
func (fs *FacilityServer) RejectFacilityRequest(ctx context.Context, in *facility.RejectFacilityRequestRequest) (*common.Result, error) {
	isConditionPassed, err := isAbleToRejectFacilityRequest(fs, in)
	if !isConditionPassed || err != nil {
		return nil, status.Error(err.Code(), err.Error())
	}

	err = fs.dbs.RejectFacilityRequest(in.RequestId, in.Reason)

	if err != nil {
		return nil, status.Error(err.Code(), err.Error())
	}

	description := fmt.Sprintf("Request ID: %d has been rejected", in.RequestId)
	return &common.Result{
		IsOk:        true,
		Description: description,
	}, nil
}

// CreateFacilityRequest is a function to create facility’s request by id
func (fs *FacilityServer) CreateFacilityRequest(ctx context.Context, in *facility.CreateFacilityRequestRequest) (*common.FacilityRequest, error) {
	isConditionPassed, err := isAbleToCreateFacilityRequest(fs, in)

	if !isConditionPassed || err != nil {
		return nil, status.Error(err.Code(), err.Error())
	}

	result, err := fs.dbs.CreateFacilityRequest(in.EventId, in.FacilityId, in.Start, in.End)

	if err != nil {
		return nil, status.Error(err.Code(), err.Error())
	}

	return result, nil
}

// GetFacilityRequestList is a function to get facility request’s of the organization
func (fs *FacilityServer) GetFacilityRequestList(ctx context.Context, in *facility.GetFacilityRequestListRequest) (*facility.GetFacilityRequestListResponse, error) {
	permission := common.Permission_UPDATE_FACILITY
	isPermission := hasPermission(in.UserId, in.OrganizationId, permission)

	if !isPermission {
		return nil, status.Error(codes.PermissionDenied, (&typing.PermissionError{Type: permission}).Error())
	}

	result, err := fs.dbs.GetFacilityRequestList(in.OrganizationId)

	if err != nil {
		return nil, status.Error(err.Code(), err.Error())
	}

	return &facility.GetFacilityRequestListResponse{
		Requests: result,
	}, nil
}

// GetFacilityRequestsListStatus is a function to get facility’s of the event
func (fs *FacilityServer) GetFacilityRequestsListStatus(ctx context.Context, in *facility.GetFacilityRequestsListStatusRequest) (*facility.GetFacilityRequestsListStatusResponse, error) {
	permission := common.Permission_UPDATE_FACILITY
	event := getEvent(in.EventId)
	isPermission := hasPermission(in.UserId, event.OrganizationId, permission)

	if !isPermission {
		return nil, status.Error(codes.PermissionDenied, (&typing.PermissionError{Type: permission}).Error())
	}

	result, err := fs.dbs.GetFacilityRequestsListStatus(in.EventId)

	if err != nil {
		return nil, status.Error(err.Code(), err.Error())
	}

	return &facility.GetFacilityRequestsListStatusResponse{
		Requests: result,
	}, nil
}

// GetFacilityRequestStatus is a function to get facility request’s of the event
func (fs *FacilityServer) GetFacilityRequestStatus(ctx context.Context, in *facility.GetFacilityRequestStatusRequest) (*common.FacilityRequest, error) {
	result, err := fs.dbs.GetFacilityRequest(in.RequestId)
	if err != nil {
		return nil, status.Error(err.Code(), err.Error())
	}

	isAbleToviewRequest, permission, err := isAbleToViewFacilityRequest(fs, in.UserId, result)
	if err != nil {
		return nil, status.Error(err.Code(), err.Error())
	}

	if !isAbleToviewRequest {
		return nil, status.Error(codes.PermissionDenied, (&typing.PermissionError{Type: permission}).Error())
	}

	return result, nil
}

// GetFacilityRequestStatusFull is a function to get facility request’s of the event
func (fs *FacilityServer) GetFacilityRequestStatusFull(ctx context.Context, in *facility.GetFacilityRequestStatusFullRequest) (*facility.FacilityRequestWithFacilityInfo, error) {
	result, err := fs.dbs.GetFacilityRequestStatusFull(in.RequestId)
	if err != nil {
		return nil, status.Error(err.Code(), err.Error())
	}

	isAbleToviewRequest, permission, err := isAbleToViewFacilityRequestFull(fs, in.UserId, result)
	if err != nil {
		return nil, status.Error(err.Code(), err.Error())
	}

	if !isAbleToviewRequest {
		return nil, status.Error(codes.PermissionDenied, (&typing.PermissionError{Type: permission}).Error())
	}

	return result, nil
}

// GetAvailableTimeOfFacility is a function to get available of facility
func (fs *FacilityServer) GetAvailableTimeOfFacility(ctx context.Context, in *facility.GetAvailableTimeOfFacilityRequest) (*facility.GetAvailableTimeOfFacilityResponse, error) {
	facilityRequests, err := fs.dbs.GetApprovedFacilityRequestList(in.FacilityId, in.Start, in.End)
	if err != nil {
		return nil, status.Error(err.Code(), err.Error())
	}

	facilityInfo, err := fs.dbs.GetFacilityInfo(in.FacilityId)
	if err != nil {
		return nil, status.Error(err.Code(), err.Error())
	}

	startTime, _ := ptypes.Timestamp(in.Start)
	finishTime, _ := ptypes.Timestamp(in.End)

	operatingHours := map[int32]*common.OperatingHour{}
	for _, operatingHour := range facilityInfo.OperatingHours {
		operatingHours[common.DayOfWeek_value[operatingHour.Day.String()]] = operatingHour
	}

	// should be change in other helper too
	day := helper.DayDifference(&startTime, &finishTime) + 1
	if day <= 0 {
		return nil, &typing.InputError{Name: "Start must be earlier than Finish"}
	}
	result := make([]*facility.GetAvailableTimeOfFacilityResponse_Day, day)
	var currentDay time.Time
	for i := range result {
		currentDay = startTime.AddDate(0, 0, i)
		operationHour := operatingHours[int32(currentDay.Weekday())]
		if operationHour == nil {
			result[i] = &facility.GetAvailableTimeOfFacilityResponse_Day{Items: nil}
			continue
		}
		startHour := operationHour.StartHour
		finishHour := operationHour.FinishHour
		hour := finishHour - startHour
		avaialbleTime := make([]bool, hour)
		for j := range avaialbleTime {
			avaialbleTime[j] = true
		}
		result[i] = &facility.GetAvailableTimeOfFacilityResponse_Day{Items: avaialbleTime}
	}

	for _, request := range facilityRequests {
		if request.Finish.Seconds >= in.Start.Seconds && request.Start.Seconds <= in.End.Seconds {
			requestStartTime, _ := ptypes.Timestamp(request.Start)
			requestFinishTime, _ := ptypes.Timestamp(request.Finish)
			index := requestStartTime.Day() - startTime.Day()
			operatiingHour := operatingHours[int32(requestStartTime.Weekday())]
			if operatiingHour == nil {
				continue
			}
			startHour := operatiingHour.StartHour
			requestStartHour := requestStartTime.Hour()
			requestFinishHour := requestFinishTime.Hour()
			for i, item := range result[index].Items {
				currentHour := int(startHour) + i
				if item && currentHour <= requestStartHour || currentHour >= requestFinishHour {
					result[index].Items[i] = false
				}
			}
		}
	}

	return &facility.GetAvailableTimeOfFacilityResponse{Day: result}, nil
}

func main() {
	port := os.Getenv("GRPC_PORT")
	lis, err := net.Listen("tcp", ":"+port)
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
