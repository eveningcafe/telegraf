package vcs_ceph

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/internal"
	"github.com/influxdata/telegraf/plugins/inputs"
	"io/ioutil"
	"log"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	plugin     = "vcs-ceph"
	typeMon    = "monitor"
	typeOsd    = "osd"
	osdPrefix  = "ceph-osd"
	monPrefix  = "ceph-mon"
	sockSuffix = "asok"
)

type Ceph struct {
	Cluster                string
	CephBinary             string
	PgGoodState            []string
	EnableMgrMetric        bool
	OsdPrefix              string
	MonPrefix              string
	SocketDir              string
	SocketSuffix           string
	CephUser               string
	CephConfig             string
	GatherAdminSocketStats bool
	GatherClusterStats     bool
	TimeoutExec            internal.Duration `toml:"timeout_exec"`
}

func (c *Ceph) Description() string {
	return "Collects performance metrics from the MON and OSD nodes in a Ceph storage cluster."
}

var sampleConfig = `
  ## This is the recommended interval to poll.  Too frequent and you will lose
  ## Recommend interval
  interval = '1m'
  ## All configuration values are optional, defaults are shown below
  ## ceph cluster name
  cluster = "my_ceph"

  ## location of ceph binary
  ceph_binary = "/usr/bin/ceph"
  ## timeout excute for each command like ceph -s, ceph osd pool stat, ceph df and ceph osd df, defaut 30s
  timeout_exec = "30s"

  ## state of pg that plugin will consider in good_state , wrap into metric: ceph_pgmap,cluster=my_ceph num_good_pgs=***
  pg_good_state = ["active+clean","active+clean+scrubbing", "active+clean+scrubbing+deep"]

  ## directory in which to look for socket files
  socket_dir = "/var/run/ceph"

  ## prefix of MON and OSD socket files, used to determine socket type
  mon_prefix = "ceph-mon"
  osd_prefix = "ceph-osd"

  ## suffix used to identify socket files
  socket_suffix = "asok"

  ## Ceph user to authenticate as
  ceph_user = "client.admin"

  ## Ceph configuration to use to locate the cluster
  ceph_config = "/etc/ceph/ceph.conf"

  ## Whether to gather statistics via the admin socket
  gather_admin_socket_stats = true

  ## Whether to gather statistics via ceph commands
  gather_cluster_stats = false
`

func (c *Ceph) SampleConfig() string {
	return sampleConfig
}

func (c *Ceph) Gather(acc telegraf.Accumulator) error {
	if c.GatherAdminSocketStats {
		if err := c.gatherAdminSocketStats(acc); err != nil {
			return err
		}
	}

	if c.GatherClusterStats {
		if err := c.gatherClusterStats(acc); err != nil {
			return err
		}
	}

	return nil
}

func (c *Ceph) gatherAdminSocketStats(acc telegraf.Accumulator) error {
	sockets, err := findSockets(c)
	if err != nil {
		return fmt.Errorf("failed to find sockets at path '%s': %v", c.SocketDir, err)
	}

	for _, s := range sockets {
		dump, err := perfDump(c.CephBinary, s)
		if err != nil {
			acc.AddError(fmt.Errorf("E! error reading from socket '%s': %v", s.socket, err))
			continue
		}
		data, err := parseDump(dump)
		if err != nil {
			acc.AddError(fmt.Errorf("E! error parsing dump from socket '%s': %v", s.socket, err))
			continue
		}
		for tag, metrics := range data {
			acc.AddFields("ceph",
				map[string]interface{}(metrics),
				map[string]string{"type": s.sockType, "id": s.sockId, "collection": tag, "cluster": c.Cluster})
		}
	}
	return nil
}



func (c *Ceph) gatherClusterStats(acc telegraf.Accumulator) error {
	jobs := []struct {
		command string
		parser  func(telegraf.Accumulator, string, string) error
	}{
		{"status", c.decodeStatus},
		{"df", c.decodeDf},
		{"osd df", c.decodeOsdDfStats},
		{"osd pool stats", c.decodeOsdPoolStats},
	}
	// For each job, execute against the cluster, parse and accumulate the data points
	for _, job := range jobs {
		output, err := c.execWithTimeout(job.command)
		if err != nil {
			if (strings.Contains(err.Error(), "status 1")) {
				return fmt.Errorf("fail to authenticate ceph cluster , please check %s, %s", c.CephConfig, c.CephUser)
			}else{
				return fmt.Errorf("error executing command: %v", err)
			}
		}
		err = job.parser(acc, output,c.Cluster)
		if err != nil {
			return fmt.Errorf("error parsing output: %v", err)
		}
	}

	return nil
}

func init() {

	inputs.Add(plugin, func() telegraf.Input {
		return &Ceph{
			Cluster:                "my_ceph",
			CephBinary:             "/usr/bin/ceph",
			TimeoutExec:            internal.Duration{Duration: time.Second * 30},
			PgGoodState:            []string{"active+clean","active+clean+scrubbing", "active+clean+scrubbing+deep"},
			OsdPrefix:              osdPrefix,
			MonPrefix:              monPrefix,
			SocketDir:              "/var/run/ceph",
			SocketSuffix:           sockSuffix,
			CephUser:               "client.admin",
			CephConfig:             "/etc/ceph/ceph.conf",
			GatherAdminSocketStats: false,
			GatherClusterStats:     false,
		}
	})
}

var perfDump = func(binary string, socket *socket) (string, error) {
	cmdArgs := []string{"--admin-daemon", socket.socket}
	if socket.sockType == typeOsd {
		cmdArgs = append(cmdArgs, "perf", "dump")
	} else if socket.sockType == typeMon {
		cmdArgs = append(cmdArgs, "perfcounters_dump")
	} else {
		return "", fmt.Errorf("ignoring unknown socket type: %s", socket.sockType)
	}

	cmd := exec.Command(binary, cmdArgs...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("error running ceph dump: %s", err)
	}

	return out.String(), nil
}

var findSockets = func(c *Ceph) ([]*socket, error) {
	listing, err := ioutil.ReadDir(c.SocketDir)
	if err != nil {
		return []*socket{}, fmt.Errorf("Failed to read socket directory '%s': %v", c.SocketDir, err)
	}
	sockets := make([]*socket, 0, len(listing))
	for _, info := range listing {
		f := info.Name()
		var sockType string
		var sockPrefix string
		if strings.HasPrefix(f, c.MonPrefix) {
			sockType = typeMon
			sockPrefix = monPrefix
		}
		if strings.HasPrefix(f, c.OsdPrefix) {
			sockType = typeOsd
			sockPrefix = osdPrefix

		}
		if sockType == typeOsd || sockType == typeMon {
			path := filepath.Join(c.SocketDir, f)
			sockets = append(sockets, &socket{parseSockId(f, sockPrefix, c.SocketSuffix), sockType, path})
		}
	}
	return sockets, nil
}

func parseSockId(fname, prefix, suffix string) string {
	s := fname
	s = strings.TrimPrefix(s, prefix)
	s = strings.TrimSuffix(s, suffix)
	s = strings.Trim(s, ".-_")
	return s
}

type socket struct {
	sockId   string
	sockType string
	socket   string
}

type metric struct {
	pathStack []string // lifo stack of name components
	value     float64
}

// Pops names of pathStack to build the flattened name for a metric
func (m *metric) name() string {
	buf := bytes.Buffer{}
	for i := len(m.pathStack) - 1; i >= 0; i-- {
		if buf.Len() > 0 {
			buf.WriteString(".")
		}
		buf.WriteString(m.pathStack[i])
	}
	return buf.String()
}

type metricMap map[string]interface{}

type taggedMetricMap map[string]metricMap

// Parses a raw JSON string into a taggedMetricMap
// Delegates the actual parsing to newTaggedMetricMap(..)
func parseDump(dump string) (taggedMetricMap, error) {
	data := make(map[string]interface{})
	err := json.Unmarshal([]byte(dump), &data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse json: '%s': %v", dump, err)
	}

	return newTaggedMetricMap(data), nil
}

// Builds a TaggedMetricMap out of a generic string map.
// The top-level key is used as a tag and all sub-keys are flattened into metrics
func newTaggedMetricMap(data map[string]interface{}) taggedMetricMap {
	tmm := make(taggedMetricMap)
	for tag, datapoints := range data {
		mm := make(metricMap)
		for _, m := range flatten(datapoints) {
			mm[m.name()] = m.value
		}
		tmm[tag] = mm
	}
	return tmm
}

// Recursively flattens any k-v hierarchy present in data.
// Nested keys are flattened into ordered slices associated with a metric value.
// The key slices are treated as stacks, and are expected to be reversed and concatenated
// when passed as metrics to the accumulator. (see (*metric).name())
func flatten(data interface{}) []*metric {
	var metrics []*metric

	switch val := data.(type) {
	case float64:
		metrics = []*metric{{make([]string, 0, 1), val}}
	case map[string]interface{}:
		metrics = make([]*metric, 0, len(val))
		for k, v := range val {
			for _, m := range flatten(v) {
				m.pathStack = append(m.pathStack, k)
				metrics = append(metrics, m)
			}
		}
	default:
		log.Printf("I! Ignoring unexpected type '%T' for value %v", val, val)
	}

	return metrics
}

// exec executes the 'ceph' command with the supplied arguments, returning JSON formatted output
func (c *Ceph) exec(command string) (string, error) {
	cmdArgs := []string{"--conf", c.CephConfig, "--name", c.CephUser, "--format", "json"}
	cmdArgs = append(cmdArgs, strings.Split(command, " ")...)

	cmd := exec.Command(c.CephBinary, cmdArgs...)

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("error running ceph %v: %s", command, err)
	}

	output := out.String()

	// Ceph doesn't sanitize its output, and may return invalid JSON.  Patch this
	// up for them, as having some inaccurate data is better than none.
	output = strings.Replace(output, "-inf", "0", -1)
	output = strings.Replace(output, "inf", "0", -1)

	return output, nil
}

// exec executes the 'ceph' command with the supplied arguments, returning JSON formatted output
func (c *Ceph) execWithTimeout(command string) (string, error) {
	cmdArgs := []string{"--conf", c.CephConfig, "--name", c.CephUser, "--format", "json"}
	cmdArgs = append(cmdArgs, strings.Split(command, " ")...)


	// Create a new context and add a timeout to it
	ctx, cancel := context.WithTimeout(context.Background(), c.TimeoutExec.Duration)
	defer cancel() // The cancel should be deferred so resources are cleaned up

	// Create the command with our context
	cmd := exec.CommandContext(ctx, c.CephBinary, cmdArgs...)

	// This time we can simply use Output() to get the result.
	out, err := cmd.Output()

	// We want to check the context error to see if the timeout was executed.
	// The error returned by cmd.Output() will be OS specific based on what
	// happens when a process is killed.
	if ctx.Err() == context.DeadlineExceeded {
		return "",fmt.Errorf("Command ceph %v timed out in %s", command, c.TimeoutExec)
	}
	// If there's no context error, we know the command completed (or errored).
	if err != nil {
		return "",fmt.Errorf("Non-zero exit code: %v", err)
	}
	output := string(out)
	output = strings.Replace(output, "-inf", "0", -1)
	output = strings.Replace(output, "inf", "0", -1)
	return output, nil

}

// CephStatus is used to unmarshal "ceph -s" output
type CephStatus struct {
	Fsid   string `json:"fsid"`
	Health struct {
		Checks struct {
			MONDOWN struct {
				Severity string `json:"severity"`
				Summary  struct {
					Message string `json:"message"`
				} `json:"summary"`
			} `json:"MON_DOWN"`
		} `json:"checks"`
		Status  string `json:"status"`
		Summary []struct {
			Severity string `json:"severity"`
			Summary  string `json:"summary"`
		} `json:"summary"`
		OverallStatus string `json:"overall_status"`
	} `json:"health"`
	ElectionEpoch int      `json:"election_epoch"`
	Quorum        []int    `json:"quorum"`
	QuorumNames   []string `json:"quorum_names"`
	Monmap        struct {
		Epoch    int    `json:"epoch"`
		Fsid     string `json:"fsid"`
		Modified string `json:"modified"`
		Created  string `json:"created"`
		Features struct {
			Persistent []string      `json:"persistent"`
			Optional   []interface{} `json:"optional"`
		} `json:"features"`
		Mons []struct {
			Rank       int    `json:"rank"`
			Name       string `json:"name"`
			Addr       string `json:"addr"`
			PublicAddr string `json:"public_addr"`
		} `json:"mons"`
	} `json:"monmap"`
	OSDMap struct {
		OSDMap struct {
			Epoch          float64 `json:"epoch"`
			NumOSDs        float64 `json:"num_osds"`
			NumUpOSDs      float64 `json:"num_up_osds"`
			NumInOSDs      float64 `json:"num_in_osds"`
			Full           bool    `json:"full"`
			NearFull       bool    `json:"nearfull"`
			NumRemappedPGs float64 `json:"num_remapped_pgs"`
		} `json:"osdmap"`
	} `json:"osdmap"`
	PGMap struct {
		PGsByState []struct {
			StateName string  `json:"state_name"`
			Count     float64 `json:"count"`
		} `json:"pgs_by_state"`
		Version          float64  `json:"version"`
		NumPGs           float64  `json:"num_pgs"`
		DataBytes        float64  `json:"data_bytes"`
		BytesUsed        float64  `json:"bytes_used"`
		BytesAvail       float64  `json:"bytes_avail"`
		BytesTotal       float64  `json:"bytes_total"`
		ReadBytesSec     float64  `json:"read_bytes_sec"`
		WriteBytesSec    float64  `json:"write_bytes_sec"`
		OpPerSec         *float64 `json:"op_per_sec"` // This field is no longer reported in ceph 10 and later
		ReadOpPerSec     float64  `json:"read_op_per_sec"`
		WriteOpPerSec    float64  `json:"write_op_per_sec"`
	} `json:"pgmap"`
	Fsmap struct {
		Epoch  int           `json:"epoch"`
		ByRank []interface{} `json:"by_rank"`
	} `json:"fsmap"`
	Mgrmap struct {
		Epoch      int    `json:"epoch"`
		ActiveGid  int    `json:"active_gid"`
		ActiveName string `json:"active_name"`
		ActiveAddr string `json:"active_addr"`
		Available  bool   `json:"available"`
		Standbys   []struct {
			Gid              int      `json:"gid"`
			Name             string   `json:"name"`
			AvailableModules []string `json:"available_modules"`
		} `json:"standbys"`
		Modules          []string      `json:"modules"`
		AvailableModules []interface{} `json:"available_modules"`
		Services         struct {
		} `json:"services"`
	} `json:"mgrmap"`
	Servicemap struct {
		Epoch    int    `json:"epoch"`
		Modified string `json:"modified"`
		Services struct {
		} `json:"services"`
	} `json:"servicemap"`
}

// decodeStatus decodes the output of 'ceph -s'
func (c *Ceph) decodeStatus(acc telegraf.Accumulator, input string, cluster string) error {
	data := &CephStatus{}
	if err := json.Unmarshal([]byte(input), data); err != nil {
		return fmt.Errorf("failed to parse json: '%s': %v", input, err)
	}

	decoders := []func(telegraf.Accumulator, *CephStatus, string) error{
		c.decodeStatusHealth,
		c.decodeStatusMonmap,
		c.decodeStatusOsdmap,
		c.decodeStatusPgmap,
		c.decodeStatusPgmapState,
	}

	for _, decoder := range decoders {
		if err := decoder(acc, data, cluster); err != nil {
			return err
		}
	}

	return nil
}

// decodeStatusHealth decodes the health portion of the output of 'ceph status'
func (c *Ceph)decodeStatusHealth(acc telegraf.Accumulator, data *CephStatus, cluster string) error {
	fields := map[string]interface{}{
		"status":         data.Health.Status,
		"overall_status": data.Health.OverallStatus,
	}
	acc.AddFields("ceph_health", fields, map[string]string{
		"cluster": cluster,
	})
	return nil
}

// decodeStatusOsdmap decodes the OSD map portion of the output of 'ceph -s'
func (c *Ceph)decodeStatusMonmap(acc telegraf.Accumulator, data *CephStatus, cluster string) error {
	var mon_list []string
	var mon_in_quorum []string
	var mon_out_quorum []string
	for _, v := range data.Monmap.Mons {
		mon_list = append(mon_list, v.Name)
	}

	for _, v := range data.QuorumNames {
		mon_in_quorum = append(mon_in_quorum, v)
	}

	for _, ca := range mon_list {
		belong := 0
		for _, in := range mon_in_quorum {
			if ca == in {
				belong = 1
			}
		}
		// if not belong to mon in quorum then belong to out quorum
		if belong == 0 {
			mon_out_quorum = append(mon_out_quorum, ca)
		}
	}

	fields := map[string]interface{}{
		"num_mon":            len(mon_list),
		"num_in_quorum":      len(mon_in_quorum),
		"num_out_of_quorum":  len(mon_out_quorum),
		"mon_list":           strings.Join(mon_list, ", "),
		"mons_in_quorum":     strings.Join(mon_in_quorum, ", "),
		"mons_out_of_quorum": strings.Join(mon_out_quorum, ", "),
	}
	acc.AddFields("ceph_mon", fields, map[string]string{
		"cluster": cluster,
	})
	return nil
}

// decodeStatusOsdmap decodes the OSD map portion of the output of 'ceph -s'
func (c *Ceph)decodeStatusOsdmap(acc telegraf.Accumulator, data *CephStatus, cluster string) error {
	fields := map[string]interface{}{
		"epoch":            data.OSDMap.OSDMap.Epoch,
		"num_osds":         data.OSDMap.OSDMap.NumOSDs,
		"num_up_osds":      data.OSDMap.OSDMap.NumUpOSDs,
		"num_in_osds":      data.OSDMap.OSDMap.NumInOSDs,
		"full":             data.OSDMap.OSDMap.Full,
		"nearfull":         data.OSDMap.OSDMap.NearFull,
		"num_remapped_pgs": data.OSDMap.OSDMap.NumRemappedPGs,
	}
	acc.AddFields("ceph_osdmap", fields, map[string]string{
		"cluster": cluster,
	})
	return nil
}

// decodeStatusPgmap decodes the PG map portion of the output of 'ceph -s'
func (c *Ceph)decodeStatusPgmap(acc telegraf.Accumulator, data *CephStatus, cluster string) error {
	fields := map[string]interface{}{
		"version":            data.PGMap.Version,
		"num_pgs":            data.PGMap.NumPGs,
		"data_bytes":         data.PGMap.DataBytes,
		"bytes_used":         data.PGMap.BytesUsed,
		"bytes_avail":        data.PGMap.BytesAvail,
		"bytes_total":        data.PGMap.BytesTotal,
		"read_bytes_sec":     data.PGMap.ReadBytesSec,
		"write_bytes_sec":    data.PGMap.WriteBytesSec,
		"op_per_sec":         data.PGMap.OpPerSec, // This field is no longer reported in ceph 10 and later
		"read_op_per_sec":    data.PGMap.ReadOpPerSec,
		"write_op_per_sec":   data.PGMap.WriteOpPerSec,
	}

	// find ceph_pgmap,cluster=my_ceph num_pgs_good=***
	num_pgs_good := float64(0)
	for _, pgState := range data.PGMap.PGsByState {
		for _, p := range c.PgGoodState {
			if p == pgState.StateName {
				num_pgs_good += pgState.Count
				break
			}
		}
	}
	fields["num_pgs_good"]=num_pgs_good

	acc.AddFields("ceph_pgmap", fields, map[string]string{
		"cluster": cluster,
	})

	return nil
}

// decodeStatusPgmapState decodes the PG map state portion of the output of 'ceph -s'
func (c *Ceph)decodeStatusPgmapState(acc telegraf.Accumulator, data *CephStatus, cluster string) error {

	for _, pgState := range data.PGMap.PGsByState {
		tags := map[string]string{
			"cluster": cluster,
			"state":   pgState.StateName,
		}
		fields := map[string]interface{}{
			"count": pgState.Count,
		}
		acc.AddFields("ceph_pgmap_state", fields, tags)


	}

	return nil
}

// CephDF is used to unmarshal 'ceph df' output
type CephDf struct {
	Stats struct {
		TotalSpace      *float64 `json:"total_space"` // pre ceph 0.84
		TotalUsed       *float64 `json:"total_used"`  // pre ceph 0.84
		TotalAvail      *float64 `json:"total_avail"` // pre ceph 0.84
		TotalBytes      *float64 `json:"total_bytes"`
		TotalUsedBytes  *float64 `json:"total_used_bytes"`
		TotalAvailBytes *float64 `json:"total_avail_bytes"`
	} `json:"stats"`
	Pools []struct {
		Name  string `json:"name"`
		Stats struct {
			KBUsed      float64  `json:"kb_used"`
			BytesUsed   float64  `json:"bytes_used"`
			Objects     float64  `json:"objects"`
			PercentUsed *float64 `json:"percent_used"`
			MaxAvail    *float64 `json:"max_avail"`
		} `json:"stats"`
	} `json:"pools"`
}

// decodeDf decodes the output of 'ceph df'
func (c *Ceph)decodeDf(acc telegraf.Accumulator, input string, cluster string) error {
	data := &CephDf{}
	if err := json.Unmarshal([]byte(input), data); err != nil {
		return fmt.Errorf("failed to parse json: '%s': %v", input, err)
	}

	// ceph.usage: records global utilization and number of objects
	fields := map[string]interface{}{
		"total_space":       data.Stats.TotalSpace,
		"total_used":        data.Stats.TotalUsed,
		"total_avail":       data.Stats.TotalAvail,
		"total_bytes":       data.Stats.TotalBytes,
		"total_used_bytes":  data.Stats.TotalUsedBytes,
		"total_avail_bytes": data.Stats.TotalAvailBytes,
	}
	acc.AddFields("ceph_usage", fields, map[string]string{
		"cluster": cluster,
	})

	// ceph.pool.usage: records per pool utilization and number of objects
	for _, pool := range data.Pools {
		tags := map[string]string{
			"name":    pool.Name,
			"cluster": cluster,
		}
		fields := map[string]interface{}{
			"kb_used":      pool.Stats.KBUsed,
			"bytes_used":   pool.Stats.BytesUsed,
			"objects":      pool.Stats.Objects,
			"percent_used": pool.Stats.PercentUsed,
			"max_avail":    pool.Stats.MaxAvail,
		}
		acc.AddFields("ceph_pool_usage", fields, tags)
	}

	return nil
}

// CephOSDPoolStats is used to unmarshal 'ceph osd pool stats' output
type CephOSDPoolStats []struct {
	PoolName     string `json:"pool_name"`
	ClientIORate struct {
		ReadBytesSec  float64  `json:"read_bytes_sec"`
		WriteBytesSec float64  `json:"write_bytes_sec"`
		OpPerSec      *float64 `json:"op_per_sec"` // This field is no longer reported in ceph 10 and later
		ReadOpPerSec  float64  `json:"read_op_per_sec"`
		WriteOpPerSec float64  `json:"write_op_per_sec"`
	} `json:"client_io_rate"`
	RecoveryRate struct {
		RecoveringObjectsPerSec float64 `json:"recovering_objects_per_sec"`
		RecoveringBytesPerSec   float64 `json:"recovering_bytes_per_sec"`
		RecoveringKeysPerSec    float64 `json:"recovering_keys_per_sec"`
	} `json:"recovery_rate"`
}

// CephOSDPoolStats is used to unmarshal 'ceph osd df' output
type CephOSDDfStats struct {
	Nodes []struct {
		ID          int     `json:"id"`
		DeviceClass string  `json:"device_class"`
		Name        string  `json:"name"`
		Type        string  `json:"type"`
		TypeID      int     `json:"type_id"`
		CrushWeight float64 `json:"crush_weight"`
		Depth       int     `json:"depth"`
		PoolWeights struct {
		} `json:"pool_weights"`
		Reweight    float64 `json:"reweight"`
		Kb          float64 `json:"kb"`
		KbUsed      float64 `json:"kb_used"`
		KbAvail     float64 `json:"kb_avail"`
		Utilization float64 `json:"utilization"`
		Var         float64 `json:"var"`
		Pgs         int     `json:"pgs"`
	} `json:"nodes"`
	Stray   []interface{} `json:"stray"`
	Summary struct {
		TotalKb            int64   `json:"total_kb"`
		TotalKbUsed        int     `json:"total_kb_used"`
		TotalKbAvail       int64   `json:"total_kb_avail"`
		AverageUtilization float64 `json:"average_utilization"`
		MinVar             float64 `json:"min_var"`
		MaxVar             float64 `json:"max_var"`
		Dev                float64 `json:"dev"`
	} `json:"summary"`
}

// decodeOsdPoolStats decodes the output of 'ceph osd pool stats'
func (c *Ceph)decodeOsdPoolStats(acc telegraf.Accumulator, input string, cluster string) error {
	data := CephOSDPoolStats{}
	if err := json.Unmarshal([]byte(input), &data); err != nil {
		return fmt.Errorf("failed to parse json: '%s': %v", input, err)
	}

	// ceph.pool.stats: records pre pool IO and recovery throughput
	for _, pool := range data {
		tags := map[string]string{
			"name":    pool.PoolName,
			"cluster": cluster,
		}
		fields := map[string]interface{}{
			"read_bytes_sec":             pool.ClientIORate.ReadBytesSec,
			"write_bytes_sec":            pool.ClientIORate.WriteBytesSec,
			"op_per_sec":                 pool.ClientIORate.OpPerSec, // This field is no longer reported in ceph 10 and later
			"read_op_per_sec":            pool.ClientIORate.ReadOpPerSec,
			"write_op_per_sec":           pool.ClientIORate.WriteOpPerSec,
			"recovering_objects_per_sec": pool.RecoveryRate.RecoveringObjectsPerSec,
			"recovering_bytes_per_sec":   pool.RecoveryRate.RecoveringBytesPerSec,
			"recovering_keys_per_sec":    pool.RecoveryRate.RecoveringKeysPerSec,
		}
		acc.AddFields("ceph_pool_stats", fields, tags)
	}

	return nil
}

// decodeOsdPoolStats decodes the output of 'ceph osd df'
func (c *Ceph) decodeOsdDfStats(acc telegraf.Accumulator, input string, cluster string) error {
	data := CephOSDDfStats{}
	if err := json.Unmarshal([]byte(input), &data); err != nil {
		return fmt.Errorf("failed to parse json: '%s': %v", input, err)
	}

	// ceph.pool.stats: records pre pool IO and recovery throughput
	for _, osd := range data.Nodes {
		tags := map[string]string{
			"name":    osd.Name,
			"class":   osd.DeviceClass,
			"cluster": cluster,
		}
		var availability float64
		if osd.Kb == 0 {
			availability = 0
		} else {
			availability = osd.KbAvail / osd.Kb * 100
		}
		fields := map[string]interface{}{
			"crush_weight": osd.CrushWeight,
			"reweight":     osd.Reweight,
			"size_gb":      osd.Kb / (1024 * 1024),
			"use_gb":       osd.KbUsed / (1024 * 1024),
			"avail_gb":     osd.KbAvail / (1024 * 1024),
			"pgs":          osd.Pgs,
			"utilization":  osd.Utilization,
			"availability": availability,
			"var":          osd.Var,
		}
		acc.AddFields("ceph_osd_stats", fields, tags)
	}

	return nil
}
