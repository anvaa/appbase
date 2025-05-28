package global

import (
	
	"fmt"
	"net"
	"strconv"
	"strings"
)

func StringToInt(str string) int {

	if str == "" {
		return 0
	}

	i, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return i
}

func StringToInt64(str string) int64 {

	if str == "" {
		return 0
	}

	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0
	}
	return i
}

func StringToBits(s string) byte {
	b, err := strconv.ParseUint(s, 10, 8)
	if err != nil {
		return 0
	}
	return uint8(b)
}

func IntToString(val int) string {
	return strconv.Itoa(val)
}

func ActToString(t int) string {

	if t == 0 {
		t = 3600 // 1 hours
	}

	// sec := t % 60
	min := t / 60
	hour := min / 60
	min = min % 60
	day := hour / 24
	hour = hour % 24

	var timeString string
	switch {
	case day > 0:
		timeString = fmt.Sprintf("%d days", day)
		if hour > 0 {
			timeString = fmt.Sprintf("%d days, %d hours", day, hour)
		}
		if min > 0 {
			timeString = fmt.Sprintf("%d days, %d hours, %d minutes", day, hour, min)
		}
	case hour > 0:
		timeString = fmt.Sprintf("%d hours", hour)
		if min > 0 {
			timeString = fmt.Sprintf("%d hours, %d minutes", hour, min)
		}
	case min > 0:
		timeString = fmt.Sprintf("%d minutes", min)
	default:
		timeString = fmt.Sprintf("%d minutes", min)
	}

	return timeString
}

func CalculateAccessTime(t string) int64 {

	min := StringToInt(t)
	if min == 0 {
		min = 1
	}

	return int64(min * 60)
}

func ShortenText(txt string, nrc int) string {
	if len(txt) > nrc {
		nrc = nrc - 3
		txt = txt[:nrc] + " .."
	}
	return txt
}

func GetIPv4Addresses() ([]string, error) {
	var ips []string
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, iface := range interfaces {
		// Check if the interface is up and skip if it's not.
		if iface.Flags&net.FlagUp == 0 {
			continue // Interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // Loopback interface
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

			// To get the IP address in IPv4 format, check if it is IPv4 first.
			if ip.To4() != nil {
				// if ip is docker ip, skip
				_ip := strings.Split(ip.String(), ".")
				if _ip[0] == "172" {
					continue
				}
				ips = append(ips, ip.String())
			}
		}
	}

	return ips, nil
}