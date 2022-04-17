package discover_ip

import (
	"net"
	"sync"
	"tasmota_backup/modules/device_info"
)

// Get array with IP addresses inside specified network.
func getHosts(cidr string) ([]string, int, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, 0, err
	}

	var ips []string
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
	}

	// remove network address and broadcast address
	lenIPs := len(ips)
	switch {
	case lenIPs < 2:
		return ips, lenIPs, nil

	default:
		return ips[1 : len(ips)-1], lenIPs - 2, nil
	}
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func Discover(cidr string) (foundDevices []device_info.TasmoDevice) {
	var wg sync.WaitGroup

	ips, _, _ := getHosts(cidr)

	for _, ip := range ips {
		wg.Add(1)
		go func(ip string) {
			defer wg.Done()
			dev, er := device_info.GetDeviceInfo(ip)
			if er == nil {
				foundDevices = append(foundDevices, dev)
			}
		}(ip)
	}
	wg.Wait()
	return
}
