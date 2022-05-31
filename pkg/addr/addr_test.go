package addr

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

var urls = []string{
	"hb://dummy.gw",
	"hb://dummy.gw/fm:1",
	"hb://dummy.gw/fm:1/1234",
	"hb://dummy.up",
	"hb://dummy.up/t1",
}

func TestParse(t *testing.T) {
	for _, u := range urls {
		addr, err := Parse(u)
		assert.NoError(t, err)

		switch t := addr.(type) {
		case *GW:
			log.Printf("      GW: %s", t.RoutingPath())
		case *Upstream:
			log.Printf("UPSTREAM: %s", t.RoutingPath())
		}
	}
}

func BenchmarkParseGW(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Parse("hb://dummy.gw/fm:1/1234")
	}
}

func BenchmarkParseUpstream(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Parse("hb://dummy.up/some/longer/path/that/we/might/see")
	}
}
