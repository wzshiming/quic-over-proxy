package quic_over_proxy

import (
	"context"
	"errors"
	"net"

	"github.com/quic-go/quic-go"
)

func DialerAsListen(dialer Dialer) ListenConfig {
	return &reverseListenConfig{dialer}
}

type reverseListenConfig struct {
	dialer Dialer
}

func (l *reverseListenConfig) Listen(ctx context.Context, network, address string) (net.Listener, error) {
	return &listener{
		ctx:     ctx,
		network: network,
		address: address,
		dialer:  l.dialer,
	}, nil
}

type listener struct {
	ctx     context.Context
	network string
	address string
	dialer  Dialer
	done    chan struct{}
}

func (l *listener) Accept() (net.Conn, error) {
	if l.done == nil {
		l.done = make(chan struct{})
	} else {
		<-l.done
	}

	for {
		conn, err := l.dialer.DialContext(l.ctx, l.network, l.address)
		if err != nil {
			if l.ctx.Err() != nil {
				return nil, l.ctx.Err()
			}

			if errors.Is(err, &quic.IdleTimeoutError{}) {
				continue
			}

			return nil, err
		} else {
			return &netConnWrapper{
				Conn: conn,
				done: l.done,
			}, nil
		}
	}

}

type netConnWrapper struct {
	done chan struct{}
	net.Conn
}

func (c *netConnWrapper) Close() error {
	c.done <- struct{}{}
	return c.Conn.Close()
}

func (l *listener) Close() error {
	return nil
}

func (l *listener) Addr() net.Addr {
	return addr{
		network: l.network,
		address: l.address,
	}
}

type addr struct {
	network string
	address string
}

func (a addr) Network() string {
	return a.network
}
func (a addr) String() string {
	return a.address
}
