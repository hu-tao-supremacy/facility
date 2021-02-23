package main

import common "onepass.app/facility/hts/common"

// FacilityInfoWithRequest is a struct to combine facility info and request
type FacilityInfoWithRequest struct {
	Info     *common.Facility
	Requests []*common.FacilityRequest
}
