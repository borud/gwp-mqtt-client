package addr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	regularGWURLs = []string{
		"hb://dummy.gw",
		"hb://dummy.gw/fm:1",
		"hb://dummy.gw/fm:1/1234",
	}

	extraSlashGWURLs = []string{
		"hb://dummy.gw/",
		"hb://dummy.gw/fm:1/",
		"hb://dummy.gw/fm:1/1234/",
	}
)

func TestGW(t *testing.T) {
	for _, u := range regularGWURLs {
		gw, err := Parse(u)
		assert.NoError(t, err)
		assert.Equal(t, u, gw.String())
	}

	// Make sure normalization works predictably
	for _, u := range extraSlashGWURLs {
		gw, err := Parse(u)
		assert.NoError(t, err)
		assert.Equal(t, u, gw.String()+"/")
	}

	_, err := Parse("http://foo.bar")
	assert.ErrorIs(t, err, ErrWrongURLScheme)

	_, err = Parse("hb://.gw")
	assert.ErrorIs(t, err, ErrInvalidGatewayID)

	addr, err := Parse("hb://some.gw/foo:1/1234567")
	assert.NoError(t, err)

	gwAddr, ok := addr.(*GW)
	assert.True(t, ok)
	assert.Equal(t, uint64(1234567), gwAddr.NodeIDAsUint64())
}
