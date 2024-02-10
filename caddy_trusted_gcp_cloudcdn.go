package caddy_trusted_gcp_cloudcdn

import (
	"net"
	"net/http"
	"net/netip"
	"strings"
	"sync"
	"time"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

func init() {
	caddy.RegisterModule(CaddyTrustedGCPCloudCDN{})
}

// The module auto trusted_proxies `GCP CloudCDN EDGE servers` from `_cloud-eoips.googleusercontent.com` TXT record
// Doc: https://cloud.google.com/cdn/docs/set-up-external-backend-internet-neg
// Range from: _cloud-eoips.googleusercontent.com
type CaddyTrustedGCPCloudCDN struct {
	// Interval to update the trusted proxies list. default: 1d
	Interval caddy.Duration `json:"interval,omitempty"`
	ranges   []netip.Prefix
	ctx      caddy.Context
	lock     *sync.RWMutex
}

func (CaddyTrustedGCPCloudCDN) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.ip_sources.gcp_cloudcdn",
		New: func() caddy.Module { return new(CaddyTrustedGCPCloudCDN) },
	}
}

func (s *CaddyTrustedGCPCloudCDN) Provision(ctx caddy.Context) error {
	s.ctx = ctx
	s.lock = new(sync.RWMutex)

	// update cron
	go func() {
		if s.Interval == 0 {
			s.Interval = caddy.Duration(24 * time.Hour) // default to 24 hours
		}
		ticker := time.NewTicker(time.Duration(s.Interval))
		s.lock.Lock()
		s.ranges, _ = s.fetchPrefixes()
		s.lock.Unlock()
		for {
			select {
			case <-ticker.C:
				prefixes, err := s.fetchPrefixes()
				if err != nil {
					break
				}
				s.lock.Lock()
				s.ranges = prefixes
				s.lock.Unlock()
			case <-s.ctx.Done():
				ticker.Stop()
				return
			}
		}
	}()
	return nil
}

func (s *CaddyTrustedGCPCloudCDN) fetchPrefixes() ([]netip.Prefix, error) {
	txtRecords, err := net.LookupTXT("_cloud-eoips.googleusercontent.com")
	if err != nil {
		return nil, err
	}
	var prefixes []netip.Prefix
	// txt like "v=spf1 ip4:34.96.0.0/20 ip4:34.127.192.0/18 ~all"
	for _, txt := range txtRecords {
		ss := strings.Fields(txt)
		for _, s := range ss {
			if strings.HasPrefix(s, "ip") {
				expr := s[strings.Index(s, ":")+1:]
				prefix, err := caddyhttp.CIDRExpressionToPrefix(expr)
				if err != nil {
					return nil, err
				}
				prefixes = append(prefixes, prefix)
			}
		}
	}
	return prefixes, nil
}

func (s *CaddyTrustedGCPCloudCDN) GetIPRanges(_ *http.Request) []netip.Prefix {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.ranges
}

// UnmarshalCaddyfile implements caddyfile.Unmarshaler. Syntax:
//
//	gcp_cloudcdn {
//	   interval <duration>
//	}
func (m *CaddyTrustedGCPCloudCDN) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	d.Next() // consume directive name

	if d.NextArg() {
		return d.ArgErr()
	}

	for nesting := d.Nesting(); d.NextBlock(nesting); {
		switch d.Val() {
		case "interval":
			if !d.NextArg() {
				return d.ArgErr()
			}
			val, err := caddy.ParseDuration(d.Val())
			if err != nil {
				return err
			}
			m.Interval = caddy.Duration(val)
		default:
			return d.ArgErr()
		}
	}

	return nil
}

// Interface guards
var (
	_ caddy.Module            = (*CaddyTrustedGCPCloudCDN)(nil)
	_ caddy.Provisioner       = (*CaddyTrustedGCPCloudCDN)(nil)
	_ caddyfile.Unmarshaler   = (*CaddyTrustedGCPCloudCDN)(nil)
	_ caddyhttp.IPRangeSource = (*CaddyTrustedGCPCloudCDN)(nil)
)
