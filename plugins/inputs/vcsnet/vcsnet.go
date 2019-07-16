package vcsnet

import (
	"fmt"
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/filter"
	"github.com/influxdata/telegraf/plugins/inputs"
	"github.com/influxdata/telegraf/plugins/inputs/system"
	"github.com/prometheus/procfs/sysfs"
	"net"
	"strings"
)

type VCSNetIOStats struct {
	filter filter.Filter
	ps     system.PS
	fs                  sysfs.FS

	skipChecks          bool
	IgnoreProtocolStats bool
	Interfaces          []string
}


func (_ *VCSNetIOStats) Description() string {
	return "Read metrics about network interface usage"
}

var netSampleConfig = `
  ## By default, telegraf gathers stats from any up interface (excluding loopback)
  ## Setting interfaces will tell it to gather these explicit interfaces,
  ## regardless of status.
  ##
  # interfaces = ["eth0"]
  ##
  ## On linux systems telegraf also collects protocol stats.
  ## Setting ignore_protocol_stats to true will skip reporting of protocol metrics.
  ##
  # ignore_protocol_stats = false
  ##
`

func (_ *VCSNetIOStats) SampleConfig() string {
	return netSampleConfig
}

func (s *VCSNetIOStats) Gather(acc telegraf.Accumulator) error {

	netio, err := s.ps.NetIO()
	if err != nil {
		return fmt.Errorf("error getting net io info: %s", err)
	}

	if s.filter == nil {
		if s.filter, err = filter.Compile(s.Interfaces); err != nil {
			return fmt.Errorf("error compiling filter: %s", err)
		}
	}

	interfaces, err := net.Interfaces()
	if err != nil {
		return fmt.Errorf("error getting list of interfaces: %s", err)
	}
	interfacesByName := map[string]net.Interface{}
	for _, iface := range interfaces {
		interfacesByName[iface.Name] = iface
	}

	for _, io := range netio {
		if len(s.Interfaces) != 0 {
			var found bool

			if s.filter.Match(io.Name) {
				found = true
			}

			if !found {
				continue
			}
		} else if !s.skipChecks {
			iface, ok := interfacesByName[io.Name]
			if !ok {
				continue
			}

			if iface.Flags&net.FlagLoopback == net.FlagLoopback {
				continue
			}

			if iface.Flags&net.FlagUp == 0 {
				continue
			}
		}

		tags := map[string]string{
			"interface": io.Name,
		}

		if err != nil {
			return fmt.Errorf("error getting list of interfaces: %s", err)
		}
		fields := map[string]interface{}{
			"bytes_sent":   io.BytesSent,
			"bytes_recv":   io.BytesRecv,
			"packets_sent": io.PacketsSent,
			"packets_recv": io.PacketsRecv,
			"err_in":       io.Errin,
			"err_out":      io.Errout,
			"drop_in":      io.Dropin,
			"drop_out":     io.Dropout,
		}
		speedBytes, upValue, err := s.getNetStatus(io.Name)
		if err != nil {
			return fmt.Errorf("error getting status of interfaces: %s", err)
		}
		fields["speedBytes"] = speedBytes
		fields["upValue"] = upValue

		acc.AddCounter("vcsnet", fields, tags)
	}

	// Get system wide stats for different network protocols
	// (ignore these stats if the call fails)
	if !s.IgnoreProtocolStats {
		netprotos, _ := s.ps.NetProto()
		fields := make(map[string]interface{})
		for _, proto := range netprotos {
			for stat, value := range proto.Stats {
				name := fmt.Sprintf("%s_%s", strings.ToLower(proto.Protocol),
					strings.ToLower(stat))
				fields[name] = value
			}
		}
		tags := map[string]string{
			"interface": "all",
		}
		acc.AddFields("vcsnet", fields, tags)
	}

	return nil

}

func init() {
	inputs.Add("vcsnet", func() telegraf.Input {
		return &VCSNetIOStats{ps: system.NewSystemPS()}
	})
}


//get list of NIC
func (s *VCSNetIOStats) getNetClassInfo() (sysfs.NetClass, error) {

	//if not test, run real file system
    var err error
	if !s.skipChecks {
		s.fs, err = sysfs.NewFS("/sys")
		if err != nil {
			return nil, fmt.Errorf("failed to open sysfs: %v", err)
		}
	}

	netClass, err := s.fs.NewNetClass()
	if err != nil {
		return netClass, fmt.Errorf("error obtaining net class info: %s", err)
	}

	return netClass, nil
}


//get speed, operstate of NIC
func (s *VCSNetIOStats) getNetStatus(interfaceName string) (uint64, uint64, error) {

	netClass, err := s.getNetClassInfo()
	if err != nil {
		return 0, 0, fmt.Errorf("could not get net class info: %s", err)
	}
	speedBytes := uint64(0)
	upValue := uint64(0)

	for _, ifaceInfo := range netClass {
		if ifaceInfo.Name != interfaceName{
			continue
		} else {

			if ifaceInfo.Speed != nil {
				speedBytes = uint64(*ifaceInfo.Speed / 8 * 1000 * 1000)
			}

			if ifaceInfo.OperState == "up" {
				upValue = 1
			} else  if ifaceInfo.OperState == "down" {
				upValue = 0
			}

		}

	}
	return speedBytes, upValue, err
}

