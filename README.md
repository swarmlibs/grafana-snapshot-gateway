# About
A snapshot gateway for self-hosted Grafana server

## Usage

```sh
usage: grafana-snapshot-gateway --grafana-url=GRAFANA-URL [<flags>]

Flags:
  --[no-]help                Show context-sensitive help (also try --help-long and --help-man).
  --listen-addr=":3003"      The address to listen on for HTTP requests.
  --advertise-addr=ADVERTISE-ADDR  
                             The address to advertise in URLs.
  --grafana-url=GRAFANA-URL  Grafana URL
  --grafana-basic-auth=GRAFANA-BASIC-AUTH  
                             Grafana credentials
  --[no-]check-snapshot-before-delete  
                             Check if snapshot exists before delete
  --[no-]version             Prints current version.
  --[no-]short-version       Print just the version number.
```
