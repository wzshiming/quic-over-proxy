package quic_over_proxy

import (
	"context"
	"net"
)

func ListenAsDialer(l ListenConfig) Dialer {
	d := &reverseDialer{
		listenConfig: l,
	}
	return d
}

type reverseDialer struct {
	ch           chan net.Conn
	listenConfig ListenConfig
	err          error
}

func (d *reverseDialer) DialContext(ctx context.Context, network, address string) (net.Conn, error) {
	network = "udp"
	if d.err != nil {
		err := d.err
		close(d.ch)
		d.err = nil
		d.ch = nil
		return nil, err
	}
	if d.ch == nil {
		err := d.start(ctx, network, address)
		if err != nil {
			return nil, err
		}
	}

	return <-d.ch, nil
}

func (d *reverseDialer) start(ctx context.Context, network, address string) error {

	listener, err := d.listenConfig.Listen(ctx, network, address)
	if err != nil {
		return err
	}

	d.ch = make(chan net.Conn)

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				d.err = err
				listener.Close()
				return
			}
			d.ch <- conn
		}
	}()
	return nil
}
