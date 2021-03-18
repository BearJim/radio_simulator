package simulator_util

import (
	"fmt"
	"net"
	"strconv"

	"git.cs.nctu.edu.tw/calee/sctp"
)

func IPsToSCTPAddr(ipAddrs []string, port int) (*sctp.SCTPAddr, error) {
	ips := []net.IPAddr{}
	for _, ipAddr := range ipAddrs {
		if ip, err := net.ResolveIPAddr("ip", ipAddr); err != nil {
			return nil, fmt.Errorf("Error resolving address '%s': %v", ipAddr, err)
		} else {
			ips = append(ips, *ip)
		}
	}
	sctpAddr := &sctp.SCTPAddr{
		IPAddrs: ips,
		Port:    port,
	}
	return sctpAddr, nil
}

func TACConfigToHexString(intString string) (hexString string) {
	tmp, _ := strconv.ParseUint(intString, 10, 32)
	hexString = fmt.Sprintf("%06x", tmp)
	return
}
