[![Build Status](https://github.com/xcaddyplugins/caddy-trusted-gcp-cloudcdn/workflows/update/badge.svg)](https://github.com/xcaddyplugins/caddy-trusted-gcp-cloudcdn)
[![Licenses](https://img.shields.io/github/license/xcaddyplugins/caddy-trusted-gcp-cloudcdn)](LICENSE)
[![donate](https://img.shields.io/badge/Donate-PayPal-green.svg)](https://www.buymeacoffee.com/illi)

# trusted_proxies GCP Cloud CDN module for `Caddy`

The module auto trusted_proxies `GCP CloudCDN EDGE servers` from `_cloud-eoips.googleusercontent.com` TXT record

Doc: https://cloud.google.com/cdn/docs/set-up-external-backend-internet-neg


## Install

The simplest, cross-platform way to get started is to download Caddy from [GitHub Releases](https://github.com/xcaddyplugins/caddy-trusted-gcp-cloudcdn/releases) and place the executable file in your PATH.


## Build from source

Requirements:

- [Go installed](https://golang.org/doc/install)
- [xcaddy](https://github.com/caddyserver/xcaddy)

Build:

```bash
$ xcaddy build --with github.com/xcaddyplugins/caddy-trusted-gcp-cloudcdn
```

## `Caddyfile` Syntax

```Caddyfile
trusted_proxies gcp_cloudcdn {
	interval <duration>
}
```

- `interval` How often to fetch the latest IP list. format is [caddy.Duration](https://caddyserver.com/docs/conventions#durations). For example `12h` represents **12 hours**, and "1d" represents **one day**. default value `1d`.

## `Caddyfile` Example

```Caddyfile
trusted_proxies gcp_cloudcdn {
	interval 1d
}
```

### `Caddyfile` Use Default Settings Example

```Caddyfile
trusted_proxies gcp_cloudcdn
```

## `Caddyfile` Global Trusted Example

Insert the following configuration of `Caddyfile` to apply it globally.

```Caddyfile
{
	servers {
		trusted_proxies gcp_cloudcdn
	}
}
```
