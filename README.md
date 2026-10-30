# godns

The Domain Name System (DNS) lets clients lookup ip addresses for a hostname or more generally any data associated with a hostname called resource records. This is an implementation of a stub and recursive DNS resolver in Go using only the standard libraries. Each resolver supports a DNS cache with a LFU eviction policy.

We support the following resource records

- A
- NS
- CNAME
- AAAA

## Setup

Make sure to have Go and Docker installed.

To get started with running the DNS resolver server. You can choose which resolver to run by setting the env variable `RESOLVER_TYPE` to `recursive` or `stub`.

Build the docker image.

```bash
docker build -t godns .
```

Run the server on UDP port 53.

```bash
docker run \
  -p 53:53/udp \
  -e RESOLVER_TYPE=recursive \
  godns
```

## Encoding / Decoding

Capture DNS response from dig.

```bash
sudo tshark -i en0 \
  -f "src host 1.1.1.1 and udp src port 53" \
  -c 1 \
  -T fields \
  -e udp.payload |
tr -d ':\r\n ' |
xxd -r -p > response.bin
```

```bash
dig @1.1.1.1 www.google.com A
```

Plan

- [x] Parse Binary DNS Message
- [x] Implement a Stub Resolver
- [x] Implement a Recursive Resolver
- [x] Add DNS Cache to Stub Resolver
- [x] Add DNS Cache to Recursive Resolver
- [x] Implement DNS UDP Server
- [ ] Dockerize DNS UDP Server

- [ ] Encode AAAA RDATA

Run DNS Server on port 53

```
docker run -p 53:53/udp godns
```
