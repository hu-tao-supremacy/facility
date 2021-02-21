package main

import (
	_ "github.com/lib/pq"
	common "onepass.app/facility/hts/common"
	facility "onepass.app/facility/hts/facility"
	typing "onepass.app/facility/internal/typing"
)

// hasPermission is mock function for account.hasPermission
func hasPermission(UserID int64, OrganizationID int64, PermissionName common.Permission) bool {
	// time.Sleep(1 * time.Second)
	return true
}

// hasEvent is mock function for organization.hasEvent
func hasEvent(UserID int64, PermissionName int64, EventID int64) bool {
	// time.Sleep(1 * time.Second)
	return true
}

// getEvent is mock function for Participant.getEvent
func getEvent(EventID int64) common.Event {
	// time.Sleep(1 * time.Second)
	return common.Event{}
}

// isAbleToCreateFacilityRequest is function to check if a facility is able to book according to user psermission
func isAbleToCreateFacilityRequest(fs *FacilityServer, in *facility.CreateFacilityRequestRequest, permission common.Permission) (bool, typing.CustomError) {
	havingPermissionChannel := make(chan bool)
	eventOwnerChannel := make(chan bool)
	overlapTimeChannel := make(chan bool)
	errorChannel := make(chan typing.CustomError)

	go func() {
		isTimeOverlap, err := fs.dbs.IsOverlapTime(in.FacilityId, in.Start, in.End)
		errorChannel <- err
		overlapTimeChannel <- isTimeOverlap
	}()

	event := getEvent(in.EventId)
	go func() {
		havingPermissionChannel <- hasPermission(in.UserId, event.OrganizationId, permission)
	}()
	go func() {
		eventOwnerChannel <- hasEvent(in.UserId, event.OrganizationId, in.EventId)
	}()

	isPermission := <-havingPermissionChannel
	isEventOwner := <-eventOwnerChannel
	overlapError := <-errorChannel
	isTimeOverlap := <-overlapTimeChannel

	close(havingPermissionChannel)
	close(eventOwnerChannel)
	close(overlapTimeChannel)
	close(errorChannel)

	if !(isPermission && isEventOwner) {
		return false, &typing.PermissionError{Type: permission}
	}

	if overlapError != nil {
		return false, overlapError
	}

	if isTimeOverlap {
		return false, &typing.AlreadyExistError{Name: "Facility is booked at that time"}
	}

	return true, nil
}

// isAbleToApproveFacilityRequest is function to check if a facility is able to be approved according to user psermission
func isAbleToApproveFacilityRequest(fs *FacilityServer, in *facility.ApproveFacilityRequestRequest, permission common.Permission) (bool, typing.CustomError) {
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
			return
		}

		havingPermissionChannel <- hasPermission(in.UserId, facility.OrganizationId, permission)
	}()

	go func() {
		isTimeOverlap, err := fs.dbs.IsOverlapTime(facilityRequest.FacilityId, facilityRequest.Start, facilityRequest.Finish)
		if err != nil {
			errorChannel <- err
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
		return false, &typing.PermissionError{Type: permission}
	}

	if isTimeOverlap {
		return false, &typing.AlreadyExistError{Name: "Facility is booked at that time"}
	}

	return true, nil
}

// isAbleToRejectFacilityRequest is function to check if a facility is able to be rejected according to user psermission
func isAbleToRejectFacilityRequest(fs *FacilityServer, in *facility.RejectFacilityRequestRequest, permission common.Permission) (bool, typing.CustomError) {
	facilityRequest, err := fs.dbs.GetFacilityRequest(in.RequestId)
	if err != nil {
		return false, err
	}

	facility, err := fs.dbs.GetFacilityInfo(facilityRequest.FacilityId)
	if err != nil {
		return false, err
	}

	isPermission := hasPermission(in.UserId, facility.OrganizationId, permission)
	if err != nil {
		return false, err
	}

	if !isPermission {
		return false, &typing.PermissionError{Type: permission}
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
			isPermissionEvent = false

		case isPermissionFacility := <-permissionFacilityChannel:
			if isPermissionFacility {
				return true, 0, nil
			}
			isPermissionFacility = false

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

	go func() {
		event := getEvent(facilityRequest.EventId)
		permissionEventChannel <- hasPermission(userID, event.OrganizationId, common.Permission_UPDATE_EVENT)
	}()
	go func() {
		permissionFacilityChannel <- hasPermission(userID, facility.OrganizationId, common.Permission_UPDATE_FACILITY)
	}()

	return handlePermissionChannel(permissionEventChannel, permissionFacilityChannel)
}

// isAbleToViewFacilityRequestFull a function to check whether user can view the targed facility request
func isAbleToViewFacilityRequestFull(fs *FacilityServer, userID int64, facilityRequestFull *facility.FacilityRequestWithFacilityInfo) (bool, common.Permission, typing.CustomError) {
	event := getEvent(facilityRequestFull.EventId)
	permissionEventChannel := make(chan bool)
	permissionFacilityChannel := make(chan bool)

	go func() {
		permissionEventChannel <- hasPermission(userID, event.OrganizationId, common.Permission_UPDATE_EVENT)
	}()
	go func() {
		permissionFacilityChannel <- hasPermission(userID, facilityRequestFull.OrganizationId, common.Permission_UPDATE_FACILITY)
	}()

	return handlePermissionChannel(permissionEventChannel, permissionFacilityChannel)
}
