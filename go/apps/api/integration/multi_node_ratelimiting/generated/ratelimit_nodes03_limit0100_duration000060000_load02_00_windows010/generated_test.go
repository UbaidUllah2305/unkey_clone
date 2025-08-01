// Code generated by go generate; DO NOT EDIT.
package ratelimit_nodes03_limit0100_duration000060000_load02_00_windows010

import (
	"testing"

	"github.com/unkeyed/unkey/go/apps/api/integration"
	run "github.com/unkeyed/unkey/go/apps/api/integration/multi_node_ratelimiting"
	"github.com/unkeyed/unkey/go/pkg/testutil"
)

func TestIntegration_RateLimit_Nodes03_Limit0100_Duration000060000_Load02_00_Windows010(t *testing.T) {
	testutil.SkipUnlessIntegration(t)

	h := integration.New(t, integration.Config{
		NumNodes: 3,
	})

	run.RunRateLimitTest(
		t,
		h,
		100,   // limit
		60000, // duration
		10,    // window count
		2,     // load factor
		3,     // node count
	)
}
