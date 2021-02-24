package helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
	common "onepass.app/facility/hts/common"
)

func TestSomething5(t *testing.T) {

	// t.Error("fizzbuzz of 1 should be '1' but have", v) // --> (6)
	// v := 1                                             // --> (4)
	assert.True(t, true, "True is true!")

	// if "1" != v {    // --> 5
	// }
	a := &common.Facility{}
	assert.Empty(t, a, "A is empty")
	// log.Println(a)
}
