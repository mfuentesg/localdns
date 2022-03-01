## localDNS Server

Simple DNS server built to handle local records and thought to be easy, fast and small.

local**DNS** is not a recursive DNS server, it means that another DNS server is required to resolve the unknown domains for `localDNS`.

## Motivation

I just wanted to have a custom DNS on my home to handle custom domains, I tried Pi-Hole, but it was too much for what I really wanted, just localDNS.

**Just for fun**

## Development

`localdns` uses an [embedded database](https://github.com/akrylysov/pogreb) to store each record in a persistent way.

If you want to contribute, just execute `go run main.go` and you will get the server running :).

## ToDo

- [x] Add persistent layer
- [x] Add support for A records
- [ ] Add support for AAAA records
- [ ] Support for IPv4 and IPv6 (?)
- [ ] Prometheus metrics
- [ ] Add go-releaser for binary
- [x] Add GitHub actions pipeline
- [ ] Add gRPC API
- [ ] Improve logging strategy

## Roadmap

- [ ] Provide a UI for easy management
- [ ] User access and permissions for UI
