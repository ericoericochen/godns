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
- [ ] Add DNS Cache to Resolvers
- [ ] Implement DNS Server
