package addr

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

// Addr is an address used for addressing various endpoints in
// the healthy buildings system.
type Addr interface {
	// RoutingPath is the representation of this address used by the router. It is responsible for ordering fields
	// in a way that provides us with a clear prefix hierarchy.
	RoutingPath() string

	// TLD returns the tld of the pseudo-hostname in HB addresses.  This can be used to distinguish between
	// the different specific address types we support, but in most cases you should use a type switch
	// structure as it is more reliable and possibly more performant.
	TLD() string

	// Return the ID of the address.
	GetID() string

	// Return the URL string for this address.
	String() string
}

// Constants for addresses
const (
	addrScheme  = "hb"
	gwTLD       = "gw"
	upstreamTLD = "up"
)

// Errors common to all address types
var (
	ErrWrongURLScheme     = errors.New("wrong URL scheme for address")
	ErrInvalidHost        = errors.New("the host part of the url is invalid")
	ErrUnknownTLD         = errors.New("unknown TLD, only know " + gwTLD + " and " + upstreamTLD)
	ErrNotGWAddress       = errors.New("not a GW address")
	ErrNotUpstreamAddress = errors.New("not a Upstream address")
)

// ParseGWAddr is a convenience function for parsing GW addresses.
func ParseGWAddr(addr string) (*GW, error) {
	a, err := Parse(addr)
	if err != nil {
		return nil, err
	}

	gwAddr, ok := a.(*GW)
	if !ok {
		return nil, ErrNotGWAddress
	}
	return gwAddr, nil
}

// ParseUpstreamAddr is a convenience function for parsing Upstream addresses.
func ParseUpstreamAddr(addr string) (*Upstream, error) {
	a, err := Parse(addr)
	if err != nil {
		return nil, err
	}

	upstreamAddr, ok := a.(*Upstream)
	if !ok {
		return nil, ErrNotUpstreamAddress
	}
	return upstreamAddr, nil
}

// Parse an address.
func Parse(addr string) (Addr, error) {
	parsedURL, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}

	if parsedURL.Scheme != addrScheme {
		return nil, fmt.Errorf("%w: %s", ErrWrongURLScheme, parsedURL.Scheme)
	}

	hostParts := strings.Split(parsedURL.Host, ".")
	if len(hostParts) != 2 {
		return nil, ErrInvalidHost
	}

	switch hostParts[1] {
	case gwTLD:
		return parseGWAddrInternal(parsedURL)
	case upstreamTLD:
		return parseUpstreamAddrInternal(parsedURL)
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnknownTLD, hostParts[1])
	}
}
