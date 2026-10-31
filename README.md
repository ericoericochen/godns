# godns

The Domain Name System (DNS) lets computers lookup ip addresses for a domain or more generally any data associated with a domain called resource records. DNS resolvers perform this lookup process by iteratively sending requests to nameservers, following referrals until they reach an authoritative nameserver that returns the final answers.

This is an implementation of a stub and recursive DNS resolver in Go using nothing more than the standard libraries.

I heavily leaned on these 2 RFCs for reference.

- [RFC 1034: Domain names - concepts and facilities](https://www.rfc-editor.org/info/rfc1034/)
- [RFC 1035: Domain names - implementation and specification](https://www.rfc-editor.org/info/rfc1035/)

---

I started this to gain a much deeper understanding of how DNS works by building out its core components from scratch. Previously, from college courses and reading some blogs, I had surface level knowledge but as with most things, there are nunances that you don't see until you get your hands dirty. So that I did. Here's a rough timeline of how the project unfolded.

I started with the DNS protocol to understand how a [message](https://www.rfc-editor.org/info/rfc1035/#section-4) was laid out in bytes to implement an encoder and decoder. Next I built a [stub resolver](https://www.rfc-editor.org/info/rfc1034/#section-5.3.1) to send queries to a public DNS server like Cloudflare's 1.1.1.1. Here I used the encoder to convert a query to a byte sequence to send as an [UDP datagram](https://en.wikipedia.org/wiki/User_Datagram_Protocol) and the decoder to parse the server response. Then I felt pretty confident to take a pass at a [recursive resolver](https://www.rfc-editor.org/info/rfc1034/#section-5.3.3), starting with a root nameserver and inspecting the authority section to query the next nameservers until I got to an authoritative nameserver. Adding a DNS cache was tricky. I learned about LFU eviction policy and how it's implemented with a [doubly linked list](https://arpitbhayani.me/blogs/lfu/). Finally, I wrapped the resolvers in an UDP server running on port 53 and tested it by sending queries via `dig`.

I walked away having a much richer understanding of how DNS works.

## Setup

Make sure to have Go and Docker installed.

To get started with running the DNS server, run the following docker commands. You can choose which resolver to run by setting the env variable `RESOLVER_TYPE` to `recursive` or `stub`.

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

To query the DNS server, you can use `dig`.

```bash
dig @127.0.0.1 www.google.com A +noedns
```

## Project Layout

```graphql
.
├── cmd
│   └── server
│       └── main.go             # DNS server entry point
├── fixtures                    # Sample DNS messages and output
├── cache.go                    # DNS cache with LFU and TTL eviction
├── cache_test.go               # DNS cache tests
├── decoder.go                  # Decode DNS messages from bytes
├── Dockerfile                  # Build and run the DNS server
├── encoder.go                  # Encode DNS messages to bytes
├── enums.go                    # DNS opcodes, response codes, and record types
├── go.mod
├── message.go                  # DNS message structures
├── message_builder.go          # DNS message builder
├── README.md
├── resolver.go                 # Shared resolver interface
├── resolver_recursive.go       # Recursive DNS resolver
├── resolver_stub.go            # Stub DNS resolver
├── root_nameservers.go         # Root nameserver addresses
├── server.config.json          # DNS server configuration
├── server.go                   # UDP DNS server
└── transport.go                # Send DNS queries over UDP
```

## Encoding / Decoding

Here's the command to capture the DNS message from `dig` to test the encoder and decoder.

```bash
sudo tshark -i en0 \
  -f "src host 1.1.1.1 and udp src port 53" \
  -c 1 \
  -T fields \
  -e udp.payload |
tr -d ':\r\n ' |
xxd -r -p > response.bin
```
