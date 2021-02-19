package database

import (
	"google.golang.org/protobuf/types/known/wrapperspb"
	facility "onepass.app/facility/hts/facility"
	model "onepass.app/facility/internal/model"
)

func convertFacilityModelToProto(data *model.Facility) *facility.Facility {
	var operatingHours, description *wrapperspb.StringValue
	if data.OperatingHours.Valid {
		operatingHours = &wrapperspb.StringValue{
			Value: data.OperatingHours.String,
		}
	}
	if data.Description.Valid {
		description = &wrapperspb.StringValue{
			Value: data.Description.String,
		}
	}
	return &facility.Facility{
		Id:             data.ID,
		OrganizationId: data.OrganizationID,
		Name:           data.Name,
		Latitude:       data.Latitude,
		Longitude:      data.Longitude,
		OperatingHours: operatingHours,
		Description:    description,
	}
}
