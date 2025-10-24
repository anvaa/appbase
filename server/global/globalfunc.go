package global

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

// StringToInt converts a string to int, returns 0 on error or empty string.
func StringToInt(str string) int {
	i, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return i
}

// StringToInt64 converts a string to int64, returns 0 on error or empty string.
func StringToInt64(str string) int64 {
	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0
	}
	return i
}

// StringToBits converts a string to uint8, returns 0 on error.
func StringToBits(s string) uint8 {
	b, err := strconv.ParseUint(s, 10, 8)
	if err != nil {
		return 0
	}
	return uint8(b)
}

// IntToString converts an int to string.
func IntToString(val int) string {
	return strconv.Itoa(val)
}

// UuidToString converts a uint32 to string.
func UuidToString(uuid uint32) string {
	return strconv.FormatUint(uint64(uuid), 10)
}

// ActToString returns a human-readable duration string from seconds.
func ActToString(t int) string {
	if t == 0 {
		t = 3600 // 1 hour
	}
	min := t / 60
	hour := min / 60
	min = min % 60
	day := hour / 24
	hour = hour % 24

	switch {
	case day > 0 && hour > 0 && min > 0:
		return fmt.Sprintf("%d days, %d hours, %d minutes", day, hour, min)
	case day > 0 && hour > 0:
		return fmt.Sprintf("%d days, %d hours", day, hour)
	case day > 0:
		return fmt.Sprintf("%d days", day)
	case hour > 0 && min > 0:
		return fmt.Sprintf("%d hours, %d minutes", hour, min)
	case hour > 0:
		return fmt.Sprintf("%d hours", hour)
	default:
		return fmt.Sprintf("%d minutes", min)
	}
}

// CalculateAccessTime converts minutes (as string) to seconds (int64).
func CalculateAccessTime(t string) int64 {
	min := StringToInt(t)
	if min == 0 {
		min = 1
	}
	return int64(min * 60)
}

// ShortenText truncates a string to nrc characters, appending " .." if truncated.
func ShortenText(txt string, nrc int) string {
	if len(txt) > nrc && nrc > 3 {
		return txt[:nrc-3] + " .."
	}
	return txt
}

// GetIPv4Addresses returns a slice of non-loopback, non-docker IPv4 addresses.
func GetIPv4Addresses() ([]string, error) {
	var ips []string
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip4 := ip.To4(); ip4 != nil {
				octets := strings.Split(ip4.String(), ".")
				if len(octets) > 0 && octets[0] == "172" { // Skip Docker IPs
					continue
				}
				ips = append(ips, ip4.String())
			}
		}
	}
	return ips, nil
}

func GetHostName() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return hostname
}
