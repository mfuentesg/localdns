## localDNS Server

`localDNS` is a simple DNS server built to handle local records and thought to be easy, fast and small. It is able to handle local and external DNS queries,
it means, if the record is not available in the `localDNS` database it will forward the query to some known public DNS server like `8.8.8.8` or `1.1.1.1`.

## Motivation

I just wanted to have a custom DNS on my home to handle custom domains, I tried `Pi-Hole`, but it was too much for what I really wanted, just `dns`.

But, the most powerful reason is **Just for fun**.

## Development

`localdns` supports embedded databases like [sqlite](https://www.sqlite.org/index.html) for persistent storing.

If you want to run your own instance of `localDNS`, go to the [releases page](https://github.com/mfuentesg/localdns/releases) and download the most convenient binary for you, or, clone this repo and execute `go run main.go` and the magic will start.

### gRPC

`localDNS` exposes a gRPC server in order, it brings you the flexibility to develop something awesome, like a custom UI to add/remove/edit records.

## Systemd

This service can be installed as systemd service. Use the following snippets as base for your configuration.
These files assume that you have installed `localDNS` at `/opt/localdns/` directory and them will check for `/opt/localdns/localdns` binary file.

**/etc/systemd/system/localdns.service**

```ini
[Unit]
Description=localDNS service
ConditionPathExists=/opt/localdns/localdns
After=network.target

[Service]
Type=simple
ExecStart=/opt/localdns/localdns
WorkingDirectory=/opt/localdns/
Restart=on-failure
User=root
Group=root
RestartSec=10
StandardOutput=syslog
StandardError=syslog
SyslogIdentifier=localdns

[Install]
WantedBy=multi-user.target
```

**/etc/systemd/system/localdns-watcher.service**

```ini
[Unit]
Description=localDNS watcher to restart localdns service on config changes
After=network.target

[Service]
Type=oneshot
ExecStart=/usr/bin/systemctl restart localdns.service

[Install]
WantedBy=multi-user.target
```

**/etc/systemd/system/localdns-watcher.path**

```ini
[Path]
PathModified=/opt/localdns/.localdns.yaml

[Install]
WantedBy=multi-user.target
```

## ToDo

- [x] Add persistent layer
- [x] Add support for A records
- [ ] Add support for AAAA records
- [x] Support for IPv4
- [ ] Support for IPv6
- [ ] Prometheus metrics
- [x] go-releaser binary creation
- [x] GitHub actions pipeline
- [x] Add gRPC API
- [ ] Logging strategy
- [ ] Configuration layer

## Roadmap

- [ ] Provide a UI for easy management
- [ ] User access and permissions for UI
