package godns

import (
	"net/netip"
	"strings"
)

// Get the parent domain
// ie. www.google.com. -> google.com.
func getParentDomain(domain string) string {
	if domain == "." {
		return domain
	}
	parts := strings.Split(domain, ".")
	parents := parts[1:]
	parentDomain := strings.Join(parents, ".")
	if parentDomain == "" {
		parentDomain += "."
	}
	return parentDomain
}

// Return the domain, its parent, its grandparent until
// the root domain.
// ie. www.google.com -> www.google.com., google.com., .com., .
func getAllDomains(domain string) []string {
	domains := []string{domain}
	for domain != "." {
		parent := getParentDomain(domain)
		domains = append(domains, parent)
		domain = parent
	}
	return domains
}

// Filter records by type.
// ie. Get all A records.
func filterRecordsByType(records []DnsRecord, qtype Type) []DnsRecord {
	filtered := []DnsRecord{}
	for _, record := range records {
		if record.Type == qtype {
			filtered = append(filtered, record)
		}
	}
	return filtered
}

// Get IPv4 addresses from A records
func getIpv4sFromRecords(records []DnsRecord) []netip.Addr {
	aRecords := filterRecordsByType(records, TypeA)
	ipv4s := []netip.Addr{}
	for _, record := range aRecords {
		ip, err := netip.ParseAddr(record.RData)
		if err != nil {
			continue
		}
		if ip.Is4() {
			ipv4s = append(ipv4s, ip)
		}
	}
	return ipv4s
}
