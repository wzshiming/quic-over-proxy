package quic_over_proxy

import (
	"context"
	"crypto/tls"
	"net"

	"github.com/quic-go/quic-go"
)

type dialer struct {
	proxyDialer Dialer
	tlsCfg      *tls.Config
	cfg         *quic.Config
	active      bool
}

func NewDialer(tlsCfg *tls.Config, cfg *quic.Config, proxyDialer Dialer, active bool) Dialer {
	return &dialer{
		proxyDialer: proxyDialer,
		tlsCfg:      tlsCfg,
		cfg:         cfg,
		active:      active,
	}
}

func (q *dialer) DialContext(ctx context.Context, network, address string) (net.Conn, error) {
	network = "udp"

	var (
		proxyConn net.Conn
		err       error
	)
	if q.proxyDialer != nil {
		proxyConn, err = q.proxyDialer.DialContext(ctx, network, address)
	} else {
		proxyConn, err = net.Dial(network, address)
	}
	if err != nil {
		return nil, err
	}

	remoteAddr, err := net.ResolveUDPAddr(network, address)
	if err != nil {
		return nil, err
	}

	conn, err := quic.Dial(ctx, &packetConnWrapper{proxyConn.(net.PacketConn)}, remoteAddr, q.tlsCfg, q.cfg)
	if err != nil {
		return nil, err
	}

	var stm quic.Stream
	if q.active {
		stm, err = conn.OpenStreamSync(ctx)
	} else {
		stm, err = conn.AcceptStream(ctx)
	}
	if err != nil {
		return nil, err
	}

	return &streamWrapper{
		Stream:     stm,
		remoteAddr: conn.RemoteAddr(),
		localAddr:  conn.LocalAddr(),
	}, nil
}
