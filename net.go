package quic_over_proxy

import (
	"context"
	"net"
)

type Dialer interface {
	DialContext(ctx context.Context, network, address string) (net.Conn, error)
}

type ListenConfig interface {
	Listen(ctx context.Context, network, address string) (net.Listener, error)
}
