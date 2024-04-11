package proxy_test

import (
	"fmt"
	"github.com/editorpost/spider/collect/proxy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewPool(t *testing.T) {

	p := proxy.NewPool("https://octopart.com/irfb3077pbf-infineon-65873800")
	p.SetCheckContent("Price and Stock")
	require.NoError(t, p.Start())

	u, err := p.GetProxyURL(nil)
	assert.NoError(t, err)
	assert.NotNil(t, u)
	fmt.Println(u.String())
}
