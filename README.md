## DNS Server

Simple DNS server built to handle local records and thought to be easy, fast and small.

## Motivation

I just wanted to have a custom DNS on my home to handle custom domains, I tried
Pi-Hole, but it was too much for what I really wanted, just local DNS.

## Development

`localdns` uses [pogreb](https://github.com/akrylysov/pogreb) as embedded database.

## ToDo

- [x] Add persistent layer
- [x] Add support for A records
- [ ] Add support for AAAA records
- [ ] Support for IPv4 and IPv6 (?)
- [ ] Prometheus metrics
- [ ] Add go-releaser for binary
- [ ] Add GitHub actions pipeline
- [ ] Add gRPC API

## Roadmap

- [ ] Provide a UI for easy management
- [ ] User access and permissions
