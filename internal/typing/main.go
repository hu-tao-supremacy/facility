package typing

import (
	"google.golang.org/grpc/codes"
	"onepass.app/facility/hts/common"
)

// CustomError is
type CustomError interface {
	Code() codes.Code
	Error() string
}

// DatabaseError is for database error
type DatabaseError struct {
	StatusCode codes.Code
	Err        error
}

func (e *DatabaseError) Error() string { return e.Err.Error() }

// Code is for getting code
func (e *DatabaseError) Code() codes.Code { return e.StatusCode }

// NotFoundError is a error for not found
type NotFoundError struct {
	Name string
}

func (e *NotFoundError) Error() string { return e.Name + ": not found" }

// Code is for getting code
func (e *NotFoundError) Code() codes.Code { return codes.NotFound }

// PermissionError is a denied permission
type PermissionError struct {
	Type common.Permission
}

func (e *PermissionError) Error() string { return e.Type.String() + " is denied" }

// Code is for getting code
func (e *PermissionError) Code() codes.Code { return codes.PermissionDenied }

// AlreadyExistError is error for existed entry
type AlreadyExistError struct {
	Name string
}

func (e *AlreadyExistError) Error() string { return e.Name + ": already exist" }

// Code is for getting code
func (e *AlreadyExistError) Code() codes.Code { return codes.AlreadyExists }

// InputError is error for input mismatch entry
type InputError struct {
	Name string
}

func (e *InputError) Error() string { return "input error: " + e.Name }

// Code is for getting code
func (e *InputError) Code() codes.Code { return codes.InvalidArgument }

// GRPCError is error for grpc client error
type GRPCError struct {
	Name string
}

func (e *GRPCError) Error() string { return "service error: " + e.Name }

// Code is for getting code
func (e *GRPCError) Code() codes.Code { return codes.Unavailable }
