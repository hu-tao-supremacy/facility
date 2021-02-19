package database

import (
	"regexp"
	"strings"

	"google.golang.org/protobuf/types/known/wrapperspb"
	facility "onepass.app/facility/hts/facility"
	model "onepass.app/facility/internal/model"
)

var camelRegex = regexp.MustCompile("[A-Z]?[a-z0-9]+")

func camelToSnakeCase(str string) string {
	matches := camelRegex.FindAllString(str, -1)
	lowers := make([]string, len(matches))

	for i, match := range matches {
		lowers[i] = strings.ToLower(match)
	}

	return strings.Join(lowers, "_")
}

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
