package quic_over_proxy

import (
	"net"

	"github.com/quic-go/quic-go"
	"github.com/wzshiming/socks5"
)

type packetConnWrapper struct {
	net.PacketConn
}

func (c *packetConnWrapper) SetReadBuffer(bytes int) error {
	// TODO: remove depend wzshiming/socks5
	socks5udpConn, ok := c.PacketConn.(*socks5.UDPConn)
	if ok {
		udpConn := socks5udpConn.PacketConn.(*net.UDPConn)
		return udpConn.SetReadBuffer(bytes)
	}
	return nil
}

func (c *packetConnWrapper) SetWriteBuffer(bytes int) error {
	// TODO: remove depend wzshiming/socks5
	socks5udpConn, ok := c.PacketConn.(*socks5.UDPConn)
	if ok {
		udpConn := socks5udpConn.PacketConn.(*net.UDPConn)
		return udpConn.SetWriteBuffer(bytes)
	}
	return nil
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
