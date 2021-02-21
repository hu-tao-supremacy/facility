package database

import (
	"encoding/json"

	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/jmoiron/sqlx/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/timestamppb"

	common "onepass.app/facility/hts/common"
	facility "onepass.app/facility/hts/facility"
	model "onepass.app/facility/internal/model"
	typing "onepass.app/facility/internal/typing"
)

func convertOperatingHoursModelToProto(OperatingHours types.JSONText) ([]*common.OperatingHour, typing.CustomError) {
	var message []*model.OperatingHour
	err := json.Unmarshal(OperatingHours, &message)

	if err != nil {
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

func convertFacilityModelToProto(data *model.Facility) (*common.Facility, typing.CustomError) {
	OperatingHours, err := convertOperatingHoursModelToProto(data.OperatingHours)
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

func convertFacilityRequestModelToProto(data *model.FacilityRequest) *common.FacilityRequest {
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

func convertFacilityRequestWithInfoModelToProto(data *model.FacilityRequestWithInfo) (*facility.FacilityRequestWithFacilityInfo, typing.CustomError) {
	var rejectReason *wrappers.StringValue
	if data.RejectReason.Valid {
		rejectReason = &wrappers.StringValue{Value: data.RejectReason.String}
	}

	OperatingHours, err := convertOperatingHoursModelToProto(data.OperatingHours)
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

func checkOperatingHours() {

}
