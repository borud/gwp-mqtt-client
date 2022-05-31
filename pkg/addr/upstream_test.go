package addr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	regularUpstreamURLs = []string{
		"hb://dummy.up",
		"hb://dummy.up/t1",
	}

	extraSlashUpostreamURLs = []string{
		"hb://dummy.up/",
		"hb://dummy.up/t1/",
	}
)

func TestUpstream(t *testing.T) {
	for _, u := range regularUpstreamURLs {
		up, err := Parse(u)
		assert.NoError(t, err)
		assert.Equal(t, u, up.String())
	}

	for _, u := range extraSlashUpostreamURLs {
		up, err := Parse(u)
		assert.NoError(t, err)
		assert.Equal(t, u, up.String()+"/")
	}

	_, err := Parse("http://foo.bar")
	assert.ErrorIs(t, err, ErrWrongURLScheme)

	_, err = Parse("hb://.up")
	assert.ErrorIs(t, err, ErrInvalidUpstreamID)
}

func TestUpstreamLongPath(t *testing.T) {
	addr, err := Parse("hb://foo.up/long/path/here")
	assert.NoError(t, err)
	assert.Equal(t, "long/path/here", addr.(*Upstream).Path)
}
