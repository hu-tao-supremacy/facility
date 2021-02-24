package main

import (
	"context"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	_ "github.com/lib/pq"
	account "onepass.app/facility/hts/account"
	common "onepass.app/facility/hts/common"
	facility "onepass.app/facility/hts/facility"
	organizer "onepass.app/facility/hts/organizer"
	participant "onepass.app/facility/hts/participant"
	"onepass.app/facility/internal/helper"
	typing "onepass.app/facility/internal/typing"
)

// hasPermission is mock function for account.hasPermission
func hasPermission(accountClient account.AccountServiceClient, userID int64, organizationID int64, permissionName common.Permission) (bool, typing.CustomError) {
	in := account.HasPermissionRequest{
		OrganizationId: organizationID,
		UserId:         userID,
		PermissionName: permissionName,
	}
	result, err := accountClient.HasPermission(context.Background(), &in)
	if err != nil {
		return false, &typing.PermissionError{Type: permissionName}
	}
	return result.IsOk, nil
}

// hasEvent is mock function for organization.hasEvent
func hasEvent(oragnizationClient organizer.OrganizationServiceClient, organizationID int64, userID int64, eventID int64) (bool, typing.CustomError) {
	in := organizer.HasEventReq{
		OrganizationId: organizationID,
		UserId:         userID,
		EventId:        eventID,
	}
	result, err := oragnizationClient.HasEvent(context.Background(), &in)
	if err != nil {
		return false, &typing.NotFoundError{Name: "The organization doesn't own event"}
	}
	return result.IsOk, nil
}

// getEvent is mock function for Participant.getEvent
func getEvent(participantClient participant.ParticipantServiceClient, eventID int64) (*common.Event, typing.CustomError) {
	in := participant.GetEventRequest{
		EventId: eventID,
	}
	result, err := participantClient.GetEvent(context.Background(), &in)
	if err != nil {
		return nil, &typing.NotFoundError{Name: "The organization doesn't own event"}
	}
	return result, nil
}

// isAbleToCreateFacilityRequest is function to check if a facility is able to book according to user psermission
func isAbleToCreateFacilityRequest(fs *FacilityServer, in *facility.CreateFacilityRequestRequest) (bool, typing.CustomError) {
	havingPermissionChannel := make(chan bool)
	eventOwnerChannel := make(chan bool)
	overlapTimeChannel := make(chan bool)
	errorChannel := make(chan typing.CustomError)

	go func() {
		isTimeOverlap, err := fs.dbs.IsOverlapTime(in.FacilityId, in.Start, in.End, true)
		errorChannel <- err
		overlapTimeChannel <- isTimeOverlap
	}()

	event, err := getEvent(fs.participant, in.EventId)
	if err != nil {
		return false, err
	}
	go func() {
		result, err := hasPermission(fs.account, in.UserId, event.OrganizationId, common.Permission_UPDATE_EVENT)
		if err != nil {
			errorChannel <- err
			havingPermissionChannel <- false
			return
		}
		havingPermissionChannel <- result
	}()
	go func() {
		result, err := hasEvent(fs.organizer, in.UserId, event.OrganizationId, in.EventId)
		if err != nil {
			errorChannel <- err
			eventOwnerChannel <- false
			return
		}
		eventOwnerChannel <- result
	}()

	isPermission := <-havingPermissionChannel
	isTimeOverlap := <-overlapTimeChannel

	close(errorChannel)
	for err := range errorChannel {
		return false, err
	}

	isEventOwner := <-eventOwnerChannel

	close(havingPermissionChannel)
	close(eventOwnerChannel)
	close(overlapTimeChannel)
	close(errorChannel)

	if !(isPermission && isEventOwner) {
		return false, &typing.PermissionError{Type: common.Permission_UPDATE_EVENT}
	}

	if isTimeOverlap {
		return false, &typing.AlreadyExistError{Name: "Facility is booked at that time"}
	}

	return true, nil
}

// isAbleToApproveFacilityRequest is function to check if a facility is able to be approved according to user psermission
func isAbleToApproveFacilityRequest(fs *FacilityServer, in *facility.ApproveFacilityRequestRequest) (bool, typing.CustomError) {
	facilityRequest, err := fs.dbs.GetFacilityRequest(in.RequestId)
	if err != nil {
		return false, err
	}

	havingPermissionChannel := make(chan bool)
	overlapTimeChannel := make(chan bool)
	errorChannel := make(chan typing.CustomError, 2)

	go func() {
		facility, err := fs.dbs.GetFacilityInfo(facilityRequest.FacilityId)
		if err != nil {
			errorChannel <- err
			havingPermissionChannel <- false
			return
		}

		result, err := hasPermission(fs.account, in.UserId, facility.OrganizationId, common.Permission_UPDATE_FACILITY)
		if err != nil {
			errorChannel <- err
			havingPermissionChannel <- false
			return
		}
		havingPermissionChannel <- result
	}()

	go func() {
		isTimeOverlap, err := fs.dbs.IsOverlapTime(facilityRequest.FacilityId, facilityRequest.Start, facilityRequest.Finish, false)
		if err != nil {
			errorChannel <- err
			overlapTimeChannel <- true
			return
		}

		overlapTimeChannel <- isTimeOverlap
	}()

	isPermission := <-havingPermissionChannel
	isTimeOverlap := <-overlapTimeChannel

	close(errorChannel)
	for err := range errorChannel {
		return false, err
	}
	close(overlapTimeChannel)
	close(havingPermissionChannel)

	if !isPermission {
		return false, &typing.PermissionError{Type: common.Permission_UPDATE_FACILITY}
	}

	if isTimeOverlap {
		return false, &typing.AlreadyExistError{Name: "Facility is booked at that time"}
	}

	return true, nil
}

// isAbleToRejectFacilityRequest is function to check if a facility is able to be rejected according to user psermission
func isAbleToRejectFacilityRequest(fs *FacilityServer, in *facility.RejectFacilityRequestRequest) (bool, typing.CustomError) {
	facilityRequest, err := fs.dbs.GetFacilityRequest(in.RequestId)
	if err != nil {
		return false, err
	}

	facility, err := fs.dbs.GetFacilityInfo(facilityRequest.FacilityId)
	if err != nil {
		return false, err
	}

	isPermission, err := hasPermission(fs.account, in.UserId, facility.OrganizationId, common.Permission_UPDATE_FACILITY)
	if err != nil {
		return false, err
	}

	if !isPermission {
		return false, &typing.PermissionError{Type: common.Permission_UPDATE_FACILITY}
	}

	return true, nil
}

func handlePermissionChannel(permissionEventChannel <-chan bool, permissionFacilityChannel <-chan bool) (bool, common.Permission, typing.CustomError) {
	var isPermissionEvent bool
	for i := 0; i < 2; i++ {
		select {
		case isPermissionEvent := <-permissionEventChannel:
			if isPermissionEvent {
				return true, 0, nil
			}

		case isPermissionFacility := <-permissionFacilityChannel:
			if isPermissionFacility {
				return true, 0, nil
			}
		}
	}

	if !isPermissionEvent {
		return false, common.Permission_UPDATE_EVENT, nil
	}
	return false, common.Permission_UPDATE_FACILITY, nil
}

// isAbleToViewFacilityRequest a function to check whether user can view the targed facility request
func isAbleToViewFacilityRequest(fs *FacilityServer, userID int64, facilityRequest *common.FacilityRequest) (bool, common.Permission, typing.CustomError) {
	facility, err := fs.dbs.GetFacilityInfo(facilityRequest.FacilityId)
	if err != nil {
		return false, 0, err
	}

	permissionEventChannel := make(chan bool)
	permissionFacilityChannel := make(chan bool)
	errorChannel := make(chan typing.CustomError)

	go func() {
		event, err := getEvent(fs.participant, facilityRequest.EventId)
		if err != nil {
			errorChannel <- err
			permissionEventChannel <- false
			return
		}
		result, err := hasPermission(fs.account, userID, event.OrganizationId, common.Permission_UPDATE_EVENT)
		if err != nil {
			errorChannel <- err
			permissionEventChannel <- false
			return
		}
		permissionEventChannel <- result
	}()
	go func() {
		result, err := hasPermission(fs.account, userID, facility.OrganizationId, common.Permission_UPDATE_FACILITY)
		if err != nil {
			errorChannel <- err
			permissionFacilityChannel <- false
			return
		}
		permissionFacilityChannel <- result
	}()

	result, permission, err := handlePermissionChannel(permissionEventChannel, permissionFacilityChannel)

	close(errorChannel)
	for err := range errorChannel {
		return false, 0, err
	}

	return result, permission, err
}

// isAbleToViewFacilityRequestFull a function to check whether user can view the targed facility request
func isAbleToViewFacilityRequestFull(fs *FacilityServer, userID int64, facilityRequestFull *facility.FacilityRequestWithFacilityInfo) (bool, common.Permission, typing.CustomError) {
	event, err := getEvent(fs.participant, facilityRequestFull.EventId)
	for err != nil {
		return false, 0, err
	}
	permissionEventChannel := make(chan bool)
	permissionFacilityChannel := make(chan bool)
	errorChannel := make(chan typing.CustomError)

	go func() {
		result, err := hasPermission(fs.account, userID, event.OrganizationId, common.Permission_UPDATE_EVENT)
		if err != nil {
			errorChannel <- err
			permissionEventChannel <- false
			return
		}
		permissionEventChannel <- result
	}()
	go func() {
		result, err := hasPermission(fs.account, userID, facilityRequestFull.OrganizationId, common.Permission_UPDATE_FACILITY)
		if err != nil {
			errorChannel <- err
			permissionFacilityChannel <- false
			return
		}
		permissionFacilityChannel <- result
	}()

	result, permission, err := handlePermissionChannel(permissionEventChannel, permissionFacilityChannel)

	close(errorChannel)
	for err := range errorChannel {
		return false, 0, err
	}

	return result, permission, err
}

// isAbleToGetAvailableTimeOfFacility a function to check whether user can check facility availability
func isAbleToGetAvailableTimeOfFacility(startTime time.Time, finishTime time.Time) typing.CustomError {
	if helper.DayDifference(startTime, finishTime)+1 <= 0 {
		return &typing.InputError{Name: "Start must be earlier than Finish"}
	}

	now := time.Now()
	if helper.DayDifference(now, finishTime) >= 30 {
		return &typing.InputError{Name: "Booking date can only be within 30 days period from today"}
	}

	dayDifference := helper.DayDifference(now, startTime)
	if dayDifference < 0 {
		return &typing.InputError{Name: "Booking time must not be in the past"}
	}

	return nil
}

// createResultEmptyArray is function to create 2D empy boolean array according to input
func createResultEmptyArray(startTime time.Time, finishTime time.Time, operatingHours map[int32]*common.OperatingHour) []*facility.GetAvailableTimeOfFacilityResponse_Day {
	dayDifference := helper.DayDifference(startTime, finishTime) + 1
	result := make([]*facility.GetAvailableTimeOfFacilityResponse_Day, dayDifference)
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

	return result
}

// generateFacilityAvailabilityResult is a function to genereate facility request from empty 2D boolean array
func generateFacilityAvailabilityResult(resultArray []*facility.GetAvailableTimeOfFacilityResponse_Day, startTime time.Time, operatingHours map[int32]*common.OperatingHour, facilityRequests []*common.FacilityRequest) *facility.GetAvailableTimeOfFacilityResponse {
	for _, request := range facilityRequests {
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
		for i, item := range resultArray[index].Items {
			currentHour := int(startHour) + i
			if item && currentHour <= requestStartHour || currentHour >= requestFinishHour {
				resultArray[index].Items[i] = false
			}
		}

	}

	return &facility.GetAvailableTimeOfFacilityResponse{Day: resultArray}
}

// getFacilityInfoWithRequests is function to preapare facility info for GetAvailableTimeOfFacility API
func getFacilityInfoWithRequests(fs *FacilityServer, facilityID int64, start *timestamp.Timestamp, end *timestamp.Timestamp) (*FacilityInfoWithRequest, typing.CustomError) {
	errorChannel := make(chan typing.CustomError, 2)
	faicilityInfoChannel := make(chan *common.Facility)
	faiclityRequestsChannel := make(chan []*common.FacilityRequest)

	go func() {
		facilityInfo, err := fs.dbs.GetFacilityInfo(facilityID)
		if err != nil {
			errorChannel <- err
		}
		faicilityInfoChannel <- facilityInfo
	}()
	go func() {
		facilityRequests, err := fs.dbs.GetApprovedFacilityRequestList(facilityID, start, end)
		if err != nil {
			errorChannel <- err
		}
		faiclityRequestsChannel <- facilityRequests
	}()

	facilityInfo := <-faicilityInfoChannel
	facilityRequests := <-faiclityRequestsChannel

	close(errorChannel)
	for err := range errorChannel {
		return nil, err
	}
	close(faicilityInfoChannel)
	close(faiclityRequestsChannel)

	return &FacilityInfoWithRequest{Info: facilityInfo, Requests: facilityRequests}, nil
}
