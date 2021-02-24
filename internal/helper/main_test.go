package helper

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	common "onepass.app/facility/hts/common"
)

func TestSomething(t *testing.T) {

	// t.Error("fizzbuzz of 1 should be '1' but have", v) // --> (6)
	// v := 1                                             // --> (4)
	assert.True(t, true, "True is true!")

	// if "1" != v {    // --> 5
	// }
	a := &common.Facility{}
	log.Println(a)
}
