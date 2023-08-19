package geoip

import "net"

func IPRange(cidr string) (start int, end int, err error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return 0, 0, err
	}

	startIP := ip.Mask(ipnet.Mask)
	endIP := make(net.IP, len(startIP))
	copy(endIP, startIP)

	for i := len(endIP) - 1; i >= len(endIP)-4; i-- {
		endIP[i] |= ^ipnet.Mask[i]
	}

	return int(IP2Int(startIP)), int(IP2Int(endIP)), nil
}

func IP2Int(ip net.IP) uint32 {
	ip = ip.To4()
	return (uint32(ip[0]) << 24) | (uint32(ip[1]) << 16) | (uint32(ip[2]) << 8) | uint32(ip[3])
}
