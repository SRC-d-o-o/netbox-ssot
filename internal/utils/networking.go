package utils

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/bl4ko/netbox-ssot/internal/constants"
)

func ReverseLookup(ipAddress string) string {
	// Create a context with the specified timeout
	TIMEOUT := 2 * time.Second //nolint:mnd
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	// Check if ipAddress contains a mask, and remove it
	ipAddress = strings.Split(ipAddress, "/")[0]

	// Use a custom resolver with the context
	resolver := &net.Resolver{}
	names, err := resolver.LookupAddr(ctx, ipAddress)
	if err != nil || len(names) == 0 {
		return ""
	}

	// Return the first domain name, stripping the trailing dot if present
	domain := strings.TrimSuffix(names[0], ".")
	return domain
}

// Function that receives hostname and performs a forward lookup
// to get the IP address. If the forward lookup fails, it returns an empty string.
func Lookup(hostname string) string {
	ips, err := net.LookupIP(hostname)
	if err != nil || len(ips) == 0 {
		return ""
	}
	return ips[0].String()
}

// Function that converts string representation of ipv4 mask (e.g. 255.255.255.128) to
// bit representation (e.g. 25).
func MaskToBits(mask string) (int, error) {
	ipMask := net.IPMask(net.ParseIP(mask).To4())
	if ipMask == nil {
		return 0, fmt.Errorf("invalid mask: %s", mask)
	}
	ones, _ := ipMask.Size()
	return ones, nil
}

// GetIPVersion returns the version of the IP address.
// It returns 4 for IPv4, 6 for IPv6, and 0 if the IP address is invalid.
func GetIPVersion(ipAddress string) int {
	ip := net.ParseIP(ipAddress)
	if ip == nil {
		return 0
	}
	if ip.To4() != nil {
		return constants.IPv4
	}
	return constants.IPv6
}

// RemoveZoneIndexFromIPAddress removes zone index from the IPv6 address:
// e.g. 2001:db8::1%eth0 -> 2001:db8::1.
// e.g. 2001:db8::1%2/64 -> 2001:db8::1/64.
func RemoveZoneIndexFromIPAddress(ipAddress string) string {
	if strings.Contains(ipAddress, "%") {
		base := strings.Split(ipAddress, "%")[0]
		maskArr := strings.Split(ipAddress, "/")
		mask := ""
		if len(maskArr) > 1 {
			mask = "/" + maskArr[1]
		}
		ipAddress = base + mask
	}
	return ipAddress
}

// SubnetContainsIPAddress checks if given IP address is part of the
// given subnet (e.g. ipAddress "172.31.4.129" and
// subnet "172.31.4.145/25").
func SubnetContainsIPAddress(ipAddress string, subnet string) bool {
	ipAddress = RemoveZoneIndexFromIPAddress(ipAddress)
	address := strings.Split(ipAddress, "/")[0]
	ip := net.ParseIP(address)
	if ip == nil {
		return false
	}
	_, ipnet, err := net.ParseCIDR(subnet)
	if err != nil {
		return false
	}
	return ipnet.Contains(ip)
}

// VerifySubnet checks if a given subnet is valid.
func VerifySubnet(subnet string) bool {
	_, _, err := net.ParseCIDR(subnet)
	return err == nil
}

// subnetsContainIPAddress checks if array of subnets contain,
// the ip address.
func subnetsContainIPAddress(ipAddress string, subnets []string) bool {
	for _, subnet := range subnets {
		if SubnetContainsIPAddress(ipAddress, subnet) {
			return true
		}
	}
	return false
}

func IsPermittedIPAddress(ipAddress string, permittedSubnets []string, ignoredSubnets []string) bool {
	if subnetsContainIPAddress(ipAddress, ignoredSubnets) {
		return false
	}
	if len(permittedSubnets) == 0 {
		return true
	}
	return subnetsContainIPAddress(ipAddress, permittedSubnets)
}

// GetmaskAndPrefixFromIPAddress extracts mask and prefix
// from a given ipAddress of format ip/mask.
// 192.168.1.1/24 --> (192.168.1.0/24, 24).
func GetPrefixAndMaskFromIPAddress(ipAddress string) (string, int, error) {
	_, ipNet, err := net.ParseCIDR(ipAddress)
	if err != nil {
		return "", 0, err
	}
	maskBits, _ := ipNet.Mask.Size()
	return ipNet.String(), maskBits, err
}
