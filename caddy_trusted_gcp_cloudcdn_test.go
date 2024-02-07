package caddy_trusted_gcp_cloudcdn

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/stretchr/testify/assert"
)

func TestFetch(t *testing.T) {
	c := &CaddyTrustedGCPCloudCDN{
		ctx: caddy.Context{
			Context: context.TODO(),
		},
		lock: new(sync.RWMutex),
	}
	prefixes, err := c.fetchPrefixes()
	assert.Nil(t, err)
	assert.True(t, len(prefixes) > 0, "prefixes is empty")
}

func TestSyntax(t *testing.T) {
	err := testSyntax(`gcp_cloudcdn`)
	assert.Nil(t, err, err)
	err = testSyntax(`gcp_cloudcdn {
		interval 12h
	}`)
	assert.Nil(t, err, err)
	err = testSyntax(`gcp_cloudcdn {
		interval 0.8h
		invalid_name 100
	}`)
	assert.NotNil(t, err, "invalid_name should be invalid")
}

func testSyntax(config string) error {
	d := caddyfile.NewTestDispenser(config)
	c := &CaddyTrustedGCPCloudCDN{}
	err := c.UnmarshalCaddyfile(d)
	if err != nil {
		return fmt.Errorf("unmarshal error for %q: %v", config, err)
	}

	ctx, cancel := caddy.NewContext(caddy.Context{Context: context.TODO()})
	defer cancel()

	err = c.Provision(ctx)
	if err != nil {
		return fmt.Errorf("provision error for %q: %v", config, err)
	}
	return nil
}
