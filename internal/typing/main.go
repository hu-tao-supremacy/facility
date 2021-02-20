package typing

import (
	"google.golang.org/grpc/codes"
	"onepass.app/facility/hts/common"
)

// DatabaseError is for database error
type DatabaseError struct {
	StatusCode codes.Code
	Err        error
}

func (e *DatabaseError) Error() string { return e.Err.Error() }

// NotFoundError is a error for not found
type NotFoundError struct {
	Name string
}

func (e *NotFoundError) Error() string { return e.Name + ": not found" }

// PermissionError is a denied permission
type PermissionError struct {
	Type common.Permission
}

func (e *PermissionError) Error() string { return e.Type.String() + " is denied" }
