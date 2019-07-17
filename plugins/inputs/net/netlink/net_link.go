package main

import (
	"fmt"
	"github.com/vishvananda/netlink"
	"net"
)

func main() {
	lo, _ := netlink.LinkByName("eth1")
	var err error

	links, err := netlink.LinkList()
	for _, io := range links {
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

		link, _ = netlink.LinkByName("io.Name")

		fields := map[string]interface{}{
			"bytes_sent":   io. BytesSent,
			"bytes_recv":   io.BytesRecv,
			"packets_sent": io.PacketsSent,
			"packets_recv": io.PacketsRecv,
			"err_in":       io.Errin,
			"err_out":      io.Errout,
			"drop_in":      io.Dropin,
			"drop_out":     io.Dropout,
			"operation":    (*link.Attrs()).OperState.String(),
		}
	}
	state := (*lo.Attrs()).OperState
	operState := state
	fmt.Printf(err)
}