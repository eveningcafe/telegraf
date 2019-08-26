package vcs_kafka

import (
	"github.com/influxdata/telegraf/testutil"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestClusterStats(t *testing.T) {

	//c := &Ceph{
	//	CephBinary:             "/usr/bin/ceph",
	//	Cluster:                "A",
	//	CephUser:               "client.vcs-monitor",
	//	CephConfig:             "/etc/ceph/ceph.conf",
	//	TimeoutExec:            "30s",
	//	GatherAdminSocketStats: false,
	//	GatherClusterStats:     true,
	//}
	k := &Kafka{
		Debug:   true,
		Detail:  true,
		Brokers: []string{"localhost:9092"},
		FilterTopics: "black",
		FilterConsummerGroups: "g2",
		Version: "v0.10.0.0",
	}
	acc := &testutil.Accumulator{}
	err := k.Gather(acc)
	require.NoError(t, err)
}
