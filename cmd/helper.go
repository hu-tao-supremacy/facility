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
func hasEvent(UserID int64, OrganizationID int64, PermissionName int64) bool {
	// time.Sleep(1 * time.Second)
	return true
}

// isAbleToCreateFacilityRequest is function to check if a facility is able to book according to user psermission
func isAbleToCreateFacilityRequest(fs *FacilityServer, in *facility.CreateFacilityRequestRequest, permssion common.Permission) (bool, typing.CustomError) {
	havingPermissionChannel := make(chan bool)
	eventOwnerChannel := make(chan bool)
	overlapTimeChannel := make(chan bool)
	errorChannel := make(chan typing.CustomError)
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
		return false, &typing.PermissionError{Type: permssion}
	}

	if overalpErr != nil {
		return false, overalpErr
	}
	if isTimeOverlap {
		return false, &typing.AlreadyExistError{Name: "Facility is booked at that time"}
	}

	return true, nil
}
