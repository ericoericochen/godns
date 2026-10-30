package godns

import "net/netip"

const (
	A_ROOT_NAMESERVER      = "a.root-servers.net."
	A_ROOT_NAMESERVER_IPV4 = "198.41.0.4"

	B_ROOT_NAMESERVER      = "b.root-servers.net."
	B_ROOT_NAMESERVER_IPV4 = "170.247.170.2"

	C_ROOT_NAMESERVER      = "c.root-servers.net."
	C_ROOT_NAMESERVER_IPV4 = "192.33.4.12"

	D_ROOT_NAMESERVER      = "d.root-servers.net."
	D_ROOT_NAMESERVER_IPV4 = "199.7.91.13"

	E_ROOT_NAMESERVER      = "e.root-servers.net."
	E_ROOT_NAMESERVER_IPV4 = "192.203.230.10"

	F_ROOT_NAMESERVER      = "f.root-servers.net."
	F_ROOT_NAMESERVER_IPV4 = "192.5.5.241"

	G_ROOT_NAMESERVER      = "g.root-servers.net."
	G_ROOT_NAMESERVER_IPV4 = "192.112.36.4"

	H_ROOT_NAMESERVER      = "h.root-servers.net."
	H_ROOT_NAMESERVER_IPV4 = "198.97.190.53"

	I_ROOT_NAMESERVER      = "i.root-servers.net."
	I_ROOT_NAMESERVER_IPV4 = "192.36.148.17"

	J_ROOT_NAMESERVER      = "j.root-servers.net."
	J_ROOT_NAMESERVER_IPV4 = "192.58.128.30"

	K_ROOT_NAMESERVER      = "k.root-servers.net."
	K_ROOT_NAMESERVER_IPV4 = "193.0.14.129"

	L_ROOT_NAMESERVER      = "l.root-servers.net."
	L_ROOT_NAMESERVER_IPV4 = "199.7.83.42"

	M_ROOT_NAMESERVER      = "m.root-servers.net."
	M_ROOT_NAMESERVER_IPV4 = "202.12.27.33"
)

var rootNameserverIpMap = map[string]netip.Addr{
	A_ROOT_NAMESERVER: netip.MustParseAddr(A_ROOT_NAMESERVER_IPV4),
	B_ROOT_NAMESERVER: netip.MustParseAddr(B_ROOT_NAMESERVER_IPV4),
	C_ROOT_NAMESERVER: netip.MustParseAddr(C_ROOT_NAMESERVER_IPV4),
	D_ROOT_NAMESERVER: netip.MustParseAddr(D_ROOT_NAMESERVER_IPV4),
	E_ROOT_NAMESERVER: netip.MustParseAddr(E_ROOT_NAMESERVER_IPV4),
	F_ROOT_NAMESERVER: netip.MustParseAddr(F_ROOT_NAMESERVER_IPV4),
	G_ROOT_NAMESERVER: netip.MustParseAddr(G_ROOT_NAMESERVER_IPV4),
	H_ROOT_NAMESERVER: netip.MustParseAddr(H_ROOT_NAMESERVER_IPV4),
	I_ROOT_NAMESERVER: netip.MustParseAddr(I_ROOT_NAMESERVER_IPV4),
	J_ROOT_NAMESERVER: netip.MustParseAddr(J_ROOT_NAMESERVER_IPV4),
	K_ROOT_NAMESERVER: netip.MustParseAddr(K_ROOT_NAMESERVER_IPV4),
	L_ROOT_NAMESERVER: netip.MustParseAddr(L_ROOT_NAMESERVER_IPV4),
	M_ROOT_NAMESERVER: netip.MustParseAddr(M_ROOT_NAMESERVER_IPV4),
}
