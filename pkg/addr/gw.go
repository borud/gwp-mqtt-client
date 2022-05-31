package addr

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// GW address. A gateway address has three components: ID, Downstream and NodeID.
// A downstream address can be a full address, containing all three, or it can contain
// prefixes to identify just the gateway or just the downstream on a given gateway.
//
//   hb://dummy.gw
//   hb://dummy.gw/fm:1
//   hb://dummy.gw/fm:1/1234
//
type GW struct {
	// The ID of the gateway. We usually use the hardware ID of the primary ethernet interface
	// as an ID.
	ID string `json:"id" db:"id"`

	// The downstream identifies what downstream driver and instance we are addressing.
	Downstream string `json:"downstream" db:"downstream"`

	// NodeID identifies the node we are addressing.
	NodeID string `json:"nodeID" db:"nodeID"`
}

// GW address related errors
var (
	ErrInvalidGatewayID = errors.New("invalid gateway ID")
)

func parseGWAddrInternal(parsedURL *url.URL) (Addr, error) {
	id := strings.Split(parsedURL.Host, ".")
	if len(id) != 2 || id[0] == "" {
		return nil, fmt.Errorf("%w: %s", ErrInvalidGatewayID, parsedURL.Host)
	}

	if id[1] != gwTLD {
		return nil, ErrUnknownTLD
	}

	gw := &GW{
		ID: id[0],
	}

	parts := strings.Split(parsedURL.Path, "/")

	// No need to check if downstream is empty since the null value is an empty string
	if len(parts) > 1 {
		gw.Downstream = parts[1]
	}

	// No need to check if nodeID is empty since the null value is an empty string
	if len(parts) > 2 {
		gw.NodeID = parts[2]
	}

	return gw, nil
}

// NodeIDAsUint64 returns the node ID as an uint64. If the value is not a valid uint64 representation
// we return 0.
func (g *GW) NodeIDAsUint64() uint64 {
	n, err := strconv.ParseUint(g.NodeID, 10, 64)
	if err != nil {
		return 0
	}
	return n
}

// GetID returns the ID element of this address.
func (g GW) GetID() string {
	return g.ID
}

// TLD returns the TLD portion of the address.
func (g GW) TLD() string {
	return gwTLD
}

func (g GW) String() string {
	sb := strings.Builder{}
	sb.WriteString(addrScheme)
	sb.WriteString("://")
	sb.WriteString(g.ID)
	sb.WriteString(".")
	sb.WriteString(gwTLD)

	if g.Downstream == "" {
		return sb.String()
	}

	sb.WriteString("/")
	sb.WriteString(g.Downstream)

	if g.NodeID == "" {
		return sb.String()
	}

	sb.WriteString("/")
	sb.WriteString(g.NodeID)

	return sb.String()
}

// RoutingPath is the representation of this address used by the router. It is responsible for ordering fields
// in a way that provides us with a clear prefix hierarchy.
func (g GW) RoutingPath() string {
	sb := strings.Builder{}
	sb.WriteString(gwTLD)

	if g.ID == "" {
		return sb.String()
	}

	sb.WriteString("/")
	sb.WriteString(g.ID)

	if g.Downstream == "" {
		return sb.String()
	}

	sb.WriteString("/")
	sb.WriteString(g.Downstream)

	if g.NodeID == "" {
		return sb.String()
	}

	sb.WriteString("/")
	sb.WriteString(g.NodeID)

	return sb.String()
}
