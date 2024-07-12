// Package constants contains network plugin constants.
package constants

const (
	// PluginCode is the code of the network plugin.
	PluginCode string = "Network"

	// ValidationTypeDNS is the DNS validation type.
	ValidationTypeDNS string = "network-dns"

	// ValidationTypeICMP is the ICMP validation type.
	ValidationTypeICMP string = "network-icmp"

	// ValidationTypeIPRange is the IP range validation type.
	ValidationTypeIPRange string = "network-ip-range"

	// ValidationTypeMTU is the MTU validation type.
	ValidationTypeMTU string = "network-mtu"

	// ValidationTypeTCPConn is the TCP connection validation type.
	ValidationTypeTCPConn string = "network-tcp-conn"

	// ValidationTypeHTTPFile is the HTTP file validation type.
	ValidationTypeHTTPFile string = "http-file"
)
