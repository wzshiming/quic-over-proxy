package quic_over_proxy

import (
	"context"
	"crypto/tls"
	"net"

	"github.com/quic-go/quic-go"
)

type listenConfig struct {
	tlsCfg *tls.Config
	cfg    *quic.Config
	active bool
}

func NewListenConfig(tlsCfg *tls.Config, cfg *quic.Config, active bool) ListenConfig {
	return &listenConfig{
		tlsCfg: tlsCfg,
		cfg:    cfg,
		active: active,
	}
}

func (q *listenConfig) Listen(ctx context.Context, network, address string) (net.Listener, error) {
	network = "udp"

	listener, err := quic.ListenAddr(address, q.tlsCfg, q.cfg)
	if err != nil {
		return nil, err
	}

	return &listenerWrapper{
		Listener: listener,
		active:   q.active,
	}, nil
}

type listenerWrapper struct {
	*quic.Listener
	ch     chan net.Conn
	active bool
}

func (l *listenerWrapper) Accept() (net.Conn, error) {
	if l.ch == nil {
		l.ch = make(chan net.Conn)
		l.startAccept(context.Background())
	}
	return <-l.ch, nil
}

func (l *listenerWrapper) startAccept(ctx context.Context) {
	go func() {
		for {
			conn, err := l.Listener.Accept(ctx)
			if err != nil {
				return
			}
			go func() {
				err := l.startAcceptStream(ctx, conn)
				if err != nil {
					return
				}
			}()
		}
	}()
}
func (l *listenerWrapper) startAcceptStream(ctx context.Context, conn quic.Connection) error {
	var (
		stm quic.Stream
		err error
	)

	for {
		if l.active {
			stm, err = conn.OpenStreamSync(ctx)
		} else {
			stm, err = conn.AcceptStream(ctx)
		}
		if err != nil {
			return err
		}

		c := &streamWrapper{
			Stream:     stm,
			localAddr:  conn.LocalAddr(),
			remoteAddr: conn.RemoteAddr(),
		}
		l.ch <- c
	}
}
