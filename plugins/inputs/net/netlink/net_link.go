package main

import (
	"fmt"
	"github.com/vishvananda/netlink"
)

func main() {
	lo, _ := netlink.LinkByName("eth1")
	state := (*lo.Attrs()).OperState
	operState := state
	fmt.Printf(operState.String())
}