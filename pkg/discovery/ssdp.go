package discovery

import (
	"github.com/bcurren/go-ssdp"
	"net"
	"time"
)

var DENON_SCHEMA = "urn:schemas-denon-com:device:ACT-Denon:1"

func DiscoverDevices() ([]net.IP, error) {
	devices, err := ssdp.Search(DENON_SCHEMA, 3*time.Second)
	if err != nil {
		return nil, err
	}
	var ips []net.IP
	for _, device := range devices {
		ips = append(ips, device.ResponseAddr.IP)
	}
	return ips, nil
}
