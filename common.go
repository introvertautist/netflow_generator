package main

import (
    "encoding/binary"
    "net"
    "fmt"
    "math/rand"
)

// Colored - Wrapper to color message
func Colored(colorString string) func(...interface{}) string {
    sprint := func(args ...interface{}) string {
        return fmt.Sprintf(colorString,
            fmt.Sprint(args...))
    }

    return sprint
}

// Colors
var (
    Info    = Colored("\033[1;32m%s\033[0m")
    Warn    = Colored("\033[1;33m%s\033[0m")
    Error   = Colored("\033[1;31m%s\033[0m")
    Debug   = Colored("\033[1;37m%s\033[0m")
    Debug2  = Colored("\033[1;35m%s\033[0m")
)

// Increment given Ip
// http://play.golang.org/p/m8TNTtygK0
func nextIP(ip net.IP) {
    for j := len(ip) - 1; j >= 0; j-- {
        ip[j]++
        if ip[j] > 0 {
            break
        }
    }
}

// GetIPRange - Get list of IP in given gange
func GetIPRange(addr string) ([]string) {
    ip, ipnet, _ := net.ParseCIDR(addr)

    var ips []string
    for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); nextIP(ip) {
        ips = append(ips, ip.String())
    }

    return ips
}

// GetRandonIP - Get random IP from list
func GetRandonIP(ipRange []string) (string) {
    return ipRange[rand.Intn(len(ipRange))]
}

// IPtoUint32 - Convert IP to uint32
func IPtoUint32(s string) uint32 {
    ip := net.ParseIP(s)
    return binary.BigEndian.Uint32(ip.To4())
}

// RandomUint32 - Generate random uint32
func RandomUint32(max int) uint32 {
    return uint32(rand.Intn(max))
}

// RandomUint16 - Generate random uint16
func RandomUint16(max int) uint16 {
    return uint16(rand.Intn(max))
}

// RandomInt - Generate random int
func RandomInt(min, max int) int {
    return rand.Intn(max-min) + min
}
