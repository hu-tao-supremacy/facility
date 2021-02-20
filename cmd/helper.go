package main

import (
	_ "github.com/lib/pq"
	common "onepass.app/facility/hts/common"
)

// hasPermission is mock function for account.hasPermission
func hasPermission(UserID int64, OrganizationID int64, PermissionName common.Permission) bool {
	// time.Sleep(1 * time.Second)
	return true
}

// hasPermission is mock function for account.hasPermission
func hasEvent(UserID int64, OrganizationID int64, PermissionName int64) bool {
	// time.Sleep(1 * time.Second)
	return true
}
