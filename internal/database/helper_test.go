package database

import (
	"log"
	"testing"

	"github.com/jmoiron/sqlx/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/proto"
	common "onepass.app/facility/hts/common"
	model "onepass.app/facility/internal/model"
	typing "onepass.app/facility/internal/typing"
)

func TestConvertOperatingHoursModelToProto(t *testing.T) { // -
	assert := assert.New(t)

	data := []byte(`[{"day": "MON", "start_hour": 10, "finish_hour": 19}]`)
	operatingHour1 := (*types.JSONText)(&data)
	expected := []*common.OperatingHour{{Day: common.DayOfWeek_MON, StartHour: 10, FinishHour: 19}}
	operatingHoursProto, err := ConvertOperatingHoursModelToProto(*operatingHour1)
	assert.Nil(err)
	assert.Equal(1, len(operatingHoursProto))
	assert.Equal(expected[0], operatingHoursProto[0])

	data = []byte(`[{"day":error "MON", "start_hour": 10, "finish_hour": 19}]`)
	operatingHour2 := (*types.JSONText)(&data)
	_, err = ConvertOperatingHoursModelToProto(*operatingHour2)
	assert.NotNil(err)

	data = []byte(`[{"day": "MON", "start_hour": 10, "finish_hour": 19}, {"day": "SAT", "start_hour": 2, "finish_hour": 9}]`)
	operatingHour3 := (*types.JSONText)(&data)
	expected3 := []*common.OperatingHour{{Day: common.DayOfWeek_MON, StartHour: 10, FinishHour: 19}, {Day: common.DayOfWeek_SAT, StartHour: 2, FinishHour: 9}}
	operatingHoursProto, err = ConvertOperatingHoursModelToProto(*operatingHour3)
	assert.Nil(err)
	assert.Equal(2, len(operatingHoursProto))
	assert.Equal(expected3[0], operatingHoursProto[0])
	assert.Equal(expected3[1], operatingHoursProto[1])

	data = []byte(`[]`)
	operatingHour4 := (*types.JSONText)(&data)
	operatingHoursProto, err = ConvertOperatingHoursModelToProto(*operatingHour4)
	assert.Nil(err)
	assert.Equal(0, len(operatingHoursProto))
}

type MockHelper struct {
	mock.Mock
}

func (hp *MockHelper) convertOperatingHoursModelToProto(operatingHours types.JSONText) ([]*common.OperatingHour, typing.CustomError) {

	hp.Called(operatingHours)
	return nil, nil
}

func TestConvertFacilityModelToProto(t *testing.T) {
	assert := assert.New(t)
	testObj := new(MockHelper)
	helper := Helper{Convert: testObj.convertOperatingHoursModelToProto}

	testObj.On("convertOperatingHoursModelToProto", types.JSONText(nil)).Return(nil, nil)

	modelFacility := model.Facility{}
	expected := common.Facility{}
	protoFacility, err := helper.convertFacilityModelToProto(&modelFacility)
	assert.True(proto.Equal(&expected, protoFacility))
	assert.Nil(err)
	assert.Equal(&expected, protoFacility)

	data := []byte(`[{"day":error "MON", "start_hour": 10, "finish_hour": 19}]`)
	Source := (types.JSONText)(data)
	testObj.On("convertOperatingHoursModelToProto", types.JSONText(Source)).Return(nil, nil)

	modelFacility = model.Facility{
		ID:             12,
		OrganizationID: 6,
		Name:           "ISE",
		Latitude:       12.2,
		Longitude:      43.3,
		OperatingHours: Source,
		Description:    "description",
	}
	protoFacility, err = helper.convertFacilityModelToProto(&modelFacility)
	expected = common.Facility{
		Id:             12,
		OrganizationId: 6,
		Name:           "ISE",
		Latitude:       12.2,
		Longitude:      43.3,
		Description:    "description"}
	assert.True(proto.Equal(&expected, protoFacility))
	assert.Nil(err)
	assert.Equal(&expected, protoFacility)
}

func TestConvertFacilityModelToProtoEmpty(t *testing.T) {
	assert := assert.New(t)
	testObj := new(MockHelper)
	helper := Helper{Convert: testObj.convertOperatingHoursModelToProto}

	data := []byte(`[{"day":error "MON", "start_hour": 10, "finish_hour": 19}]`)
	Source := (types.JSONText)(data)
	mockOperationHours := []common.OperatingHour{{
		Day:        common.DayOfWeek_WED,
		StartHour:  2,
		FinishHour: 10,
	}}
	testObj.On("convertOperatingHoursModelToProto", types.JSONText(Source)).Return(mockOperationHours, nil)

	modelFacility := model.Facility{
		ID:             12,
		OrganizationID: 6,
		Name:           "ISE",
		Latitude:       12.2,
		Longitude:      43.3,
		OperatingHours: Source,
		Description:    "description",
	}
	protoFacility, err := helper.convertFacilityModelToProto(&modelFacility)
	expected := common.Facility{
		Id:             12,
		OrganizationId: 6,
		Name:           "ISE",
		Latitude:       12.2,
		Longitude:      43.3,
		Description:    "description"}
	log.Println(protoFacility, "da")
	assert.True(proto.Equal(&expected, protoFacility))
	assert.Nil(err)
	assert.Equal(&expected, protoFacility)
}
