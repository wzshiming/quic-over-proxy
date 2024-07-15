package quic_over_proxy

import (
	"net"

	"github.com/quic-go/quic-go"
)

type packetConnWithSetBuffer interface {
	net.PacketConn
	SetReadBuffer(bytes int) error
	SetWriteBuffer(bytes int) error
}

type streamWrapper struct {
	quic.Stream
	localAddr  net.Addr
	remoteAddr net.Addr
}

func (c *streamWrapper) LocalAddr() net.Addr {
	return c.localAddr
}

func (c *streamWrapper) RemoteAddr() net.Addr {
	return c.remoteAddr
}
