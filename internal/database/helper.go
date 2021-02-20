package database

import (
	common "onepass.app/facility/hts/common"
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
