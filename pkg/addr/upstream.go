package addr

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

// Upstream is an upstream address. It consists of two parts: an ID and a path. The ID is the
// ID of the upstream, or more typically, the backend instance.  The path component of the
// address is used when generating responses to commands to tell both the machinery that maps
// addresses to MQTT topics and the receiving end what it needs to know in order to correlate
// responses.
type Upstream struct {
	// ID is the identity of the upstream
	ID string `json:"id" db:"id"`

	// Path is the (optional) path of the upstream.
	Path string `json:"path" db:"path"`
}

// Upstream related errors
var (
	ErrInvalidUpstreamID = errors.New("invalid upstream ID")
)

func parseUpstreamAddrInternal(parsedURL *url.URL) (Addr, error) {
	// the host portion should have exactly two components and the first component
	// cannot be empty
	id := strings.Split(parsedURL.Host, ".")
	if len(id) != 2 || id[0] == "" {
		return nil, fmt.Errorf("%w: %s", ErrInvalidUpstreamID, parsedURL.Host)
	}

	if id[1] != upstreamTLD {
		return nil, ErrUnknownTLD
	}

	if parsedURL.Path != "" {
		// slightly naughty to modify this in place
		parsedURL.Path = strings.TrimRight(strings.TrimLeft(parsedURL.Path, "/"), "/")
	}

	return &Upstream{
		ID:   strings.ToLower(id[0]),
		Path: strings.ToLower(parsedURL.Path),
	}, nil
}

// GetID returns the ID element of this address.
func (u Upstream) GetID() string {
	return u.ID
}

// TLD returns the TLD portion of the address.
func (u Upstream) TLD() string {
	return upstreamTLD
}

func (u Upstream) String() string {
	sb := strings.Builder{}
	sb.WriteString(addrScheme)
	sb.WriteString("://")
	sb.WriteString(u.ID)
	sb.WriteString(".")
	sb.WriteString(upstreamTLD)

	if u.Path != "" {
		sb.WriteString("/")
		sb.WriteString(u.Path)
	}

	return sb.String()
}

// RoutingPath is the representation of this address used by the router. It is responsible for ordering fields
// in a way that provides us with a clear prefix hierarchy.
func (u Upstream) RoutingPath() string {
	sb := strings.Builder{}
	sb.WriteString(upstreamTLD)

	if u.ID == "" {
		return sb.String()
	}

	sb.WriteString("/")
	sb.WriteString(u.ID)

	if u.Path == "" {
		return sb.String()
	}

	sb.WriteString("/")
	sb.WriteString(u.Path)
	return sb.String()
}
