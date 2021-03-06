## localdns Server

`localdns` is a simple DNS server built to handle local records and thought to be easy, fast and small. It is able to
handle local and external DNS queries,
it means, if the record is not available in the `localdns` database it will forward the query to some known public DNS
server like `8.8.8.8` or `1.1.1.1`.

## Motivation

I just wanted to have a custom DNS on my home to handle custom domains, I tried `Pi-Hole`, but it was too much for what
I really wanted, just `dns`.

But, the most powerful reason is **Just for fun**.

## Development

`localdns` uses [sqlite](https://www.sqlite.org/index.html) as storage layer.

If you want to run your own instance of `localdns`, go to
the [releases page](https://github.com/mfuentesg/localdns/releases) and download the most convenient binary for you, or,
clone this repo and execute `go run main.go` and the magic will start.

### Config file

`localdns` is fully configurable, just edit the autogenerated `.localdns.yaml` file and change it as your convenience.

## Migrations

`localdns` uses [golang-migrate](https://github.com/golang-migrate/migrate) for the migration system, all sql files will
be included as part of the final binary.
Migrations can be found at [`./migrations`](https://github.com/mfuentesg/localdns/tree/master/migrations) folder, and
them are applied once you start the `localdns` server.

### gRPC

`localdns` exposes a gRPC server in order, it brings you the flexibility to develop something awesome, like a custom UI
to add/remove/edit records or a CLI.

## Systemd

This service can be installed as systemd service. Use the following snippets as base for your own configuration.
These files assume that you have installed `localdns` at `/opt/localdns/` directory and them will check
for `/opt/localdns/localdns` binary file.

**/etc/systemd/system/localdns.service**

```ini
[Unit]
Description = localdns service
ConditionPathExists = /opt/localdns/localdns
After = network.target

[Service]
Type = simple
ExecStart = /opt/localdns/localdns
WorkingDirectory = /opt/localdns/
Restart = on-failure
User = root
Group = root
RestartSec = 10
StandardOutput = syslog
StandardError = syslog
SyslogIdentifier = localdns

[Install]
WantedBy = multi-user.target
```

**/etc/systemd/system/localdns-watcher.service**

```ini
[Unit]
Description = localdns watcher to restart localdns service on config changes
After = network.target

[Service]
Type = oneshot
ExecStart = /usr/bin/systemctl restart localdns.service

[Install]
WantedBy = multi-user.target
```

**/etc/systemd/system/localdns-watcher.path**

```ini
[Path]
PathModified = /opt/localdns/.localdns.yaml

[Install]
WantedBy = multi-user.target
```

## ToDo

- [x] Add persistent layer
- [x] Add support for A records
- [x] Support for IPv4
- [x] go-releaser binary creation
- [x] GitHub actions pipeline
- [x] Add gRPC API
- [x] Logging strategy
- [x] Configuration layer
- [ ] Prometheus metrics
- [ ] DNS over https support
- [ ] DNS over tls support
