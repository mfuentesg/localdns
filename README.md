## localDNS Server

`localDNS` is a simple DNS server built to handle local records and thought to be easy, fast and small. It is able to handle local and external DNS queries,
it means, if the record is not available in the `localDNS` database it will forward the query to some known public DNS server like `8.8.8.8` or `1.1.1.1`.

## Motivation

I just wanted to have a custom DNS on my home to handle custom domains, I tried `Pi-Hole`, but it was too much for what I really wanted, just `dns`.

But, the most powerful reason is **Just for fun**.

## Development

`localdns` uses an [embedded database](https://github.com/akrylysov/pogreb) to store each record in a persistent way.

If you want to run your own instance of `localDNS`, go the releases and download the most convenient binary for you or, clone this repo and execute `go run main.go` and the magic will start.

### gRPC

`localDNS` exposes a gRPC server in order, it brings you the flexibility to develop something awesome, like a custom UI to add/remove/edit records.

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
