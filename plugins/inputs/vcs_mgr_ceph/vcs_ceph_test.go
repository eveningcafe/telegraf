package vcs_mgr_ceph

import (
	"testing"

	"github.com/influxdata/telegraf/testutil"
)

func TestTrig(t *testing.T) {
	c := &VCSCeph{
	}

	for i := 0.0; i < 10.0; i++ {

		var acc testutil.Accumulator


		c.Gather(&acc)

		fields := make(map[string]interface{})


		acc.AssertContainsFields(t, "v", fields)
	}
}
