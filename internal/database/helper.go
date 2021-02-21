package database

import (
	"github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/protobuf/types/known/timestamppb"

	common "onepass.app/facility/hts/common"
	facility "onepass.app/facility/hts/facility"
	model "onepass.app/facility/internal/model"
)

func convertFacilityModelToProto(data *model.Facility) *common.Facility {
	return &common.Facility{
		Id:             data.ID,
		OrganizationId: data.OrganizationID,
		Name:           data.Name,
		Latitude:       data.Latitude,
		Longitude:      data.Longitude,
		OperatingHours: data.OperatingHours,
		Description:    data.Description,
	}
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

func convertFacilityRequestWithInfoModelToProto(data *model.FacilityRequestWithInfo) *facility.FacilityRequestWithFacilityInfo {
	var rejectReason *wrappers.StringValue
	if data.RejectReason.Valid {
		rejectReason = &wrappers.StringValue{Value: data.RejectReason.String}
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
		OperatingHours: data.OperatingHours,
		Description:    data.Description,
	}
}
