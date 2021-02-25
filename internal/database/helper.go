package database

import (
	"encoding/json"
	"time"

	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/jmoiron/sqlx/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/timestamppb"

	common "onepass.app/facility/hts/common"
	facility "onepass.app/facility/hts/facility"
	"onepass.app/facility/internal/helper"
	model "onepass.app/facility/internal/model"
	typing "onepass.app/facility/internal/typing"
)

// ConvertOperatingHoursModelToProto is fuction to convert operationHours to proto
func ConvertOperatingHoursModelToProto(operatingHours types.JSONText) ([]*common.OperatingHour, typing.CustomError) {
	var message []*model.OperatingHour

	if err := json.Unmarshal(operatingHours, &message); err != nil {
		return nil, &typing.DatabaseError{StatusCode: codes.DataLoss, Err: err}
	}

	result := make([]*common.OperatingHour, len(message))
	for i, OperatingHour := range message {
		result[i] = &common.OperatingHour{
			Day:        common.DayOfWeek(common.DayOfWeek_value[OperatingHour.Day]),
			StartHour:  OperatingHour.StartHour,
			FinishHour: OperatingHour.FinishHour,
		}
	}
	return result, nil
}

// OperatingHoursModelToProto type of function to inject to helper
type OperatingHoursModelToProto func(operatingHours types.JSONText) ([]*common.OperatingHour, typing.CustomError)

// Helper is struct to inject function that can m=be mock
type Helper struct {
	Convert       OperatingHoursModelToProto
	DayDifference helper.DayDifferenceFunc
}

func (dbHelper *Helper) convertFacilityModelToProto(data *model.Facility) (*common.Facility, typing.CustomError) {
	OperatingHours, err := dbHelper.Convert(data.OperatingHours)
	if err != nil {
		return nil, err
	}

	return &common.Facility{
		Id:             data.ID,
		OrganizationId: data.OrganizationID,
		Name:           data.Name,
		Latitude:       data.Latitude,
		Longitude:      data.Longitude,
		OperatingHours: OperatingHours,
		Description:    data.Description,
	}, nil
}

func (dbHelper *Helper) convertFacilityRequestModelToProto(data *model.FacilityRequest) *common.FacilityRequest {
	var rejectReason *wrappers.StringValue
	if data.RejectReason.Valid {
		rejectReason = &wrappers.StringValue{Value: data.RejectReason.String}
	}
	return &common.FacilityRequest{
		Id:           data.ID,
		EventId:      data.EventID,
		FacilityId:   data.FacilityID,
		Status:       common.Status(common.Status_value[data.Status]),
		RejectReason: rejectReason,
		Start:        timestamppb.New(data.Start),
		Finish:       timestamppb.New(data.Finish),
	}
}

func (dbHelper *Helper) convertFacilityRequestWithInfoModelToProto(data *model.FacilityRequestWithInfo) (*facility.FacilityRequestWithFacilityInfo, typing.CustomError) {
	var rejectReason *wrappers.StringValue
	if data.RejectReason.Valid {
		rejectReason = &wrappers.StringValue{Value: data.RejectReason.String}
	}

	OperatingHours, err := dbHelper.Convert(data.OperatingHours)
	if err != nil {
		return nil, err
	}

	return &facility.FacilityRequestWithFacilityInfo{
		Id:             data.ID,
		EventId:        data.EventID,
		FacilityId:     data.FacilityID,
		Status:         common.Status(common.Status_value[data.Status]),
		RejectReason:   rejectReason,
		Start:          timestamppb.New(data.Start),
		Finish:         timestamppb.New(data.Finish),
		OrganizationId: data.OrganizationID,
		FacilityName:   data.FacilityName,
		Latitude:       data.Latitude,
		Longitude:      data.Longitude,
		OperatingHours: OperatingHours,
		Description:    data.Description,
	}, nil
}

func (dbHelper *Helper) checkDateInput(start time.Time, finish time.Time, operatingHours []*common.OperatingHour) typing.CustomError {
	if dbHelper.DayDifference(start, finish) != 0 {
		return &typing.InputError{Name: "Start and Finish must be the same day"}
	}

	now := time.Now()
	dayDifferenceFromNow := dbHelper.DayDifference(now, start)
	if dayDifferenceFromNow >= 30 {
		return &typing.InputError{Name: "Booking date can only be within 30 days period from today"}
	}

	HourStart, MinuteStart, secondStart := start.Clock()
	HourFinish, MinuteFinish, secondFinish := finish.Clock()

	if dayDifferenceFromNow < 0 || (dayDifferenceFromNow == 0 && HourStart < now.Hour()) {
		return &typing.InputError{Name: "Booking time must not be in the past"}
	}
	if (MinuteStart + secondStart + MinuteFinish + secondFinish) != 0 {
		return &typing.InputError{Name: "Minutes and seconds must be 0"}
	}

	if HourStart > HourFinish {
		return &typing.InputError{Name: "Start must be earlier than Finish"}
	}

	weekDayStart := start.Weekday()
	var operatingHour *common.OperatingHour
	for _, value := range operatingHours {
		day := int(value.Day.Number())
		if day == int(weekDayStart) {
			operatingHour = value
		}
	}
	if operatingHour == nil {
		return &typing.InputError{Name: "Not in operatingHours"}
	}

	isStartAfterOpening := int(operatingHour.StartHour) <= HourStart
	isFinishBeforeClose := HourFinish <= int(operatingHour.FinishHour)
	if !isStartAfterOpening || !isFinishBeforeClose {
		return &typing.InputError{Name: "Not in operatingHours"}
	}

	return nil
}
