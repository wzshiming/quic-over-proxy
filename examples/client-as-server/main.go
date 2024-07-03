package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	quic_over_proxy "github.com/wzshiming/quic-over-proxy"
	"github.com/wzshiming/quic-over-proxy/examples/utils"
	"github.com/wzshiming/socks5"
)

func main() {
	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			fmt.Fprint(w, "hello, world!")
		case http.MethodPost:
			result, err := io.ReadAll(r.Body)
			if err != nil {
				fmt.Fprint(w, "error:", err.Error())
				return
			}
			fmt.Fprintf(w, "echo:'%s'", string(result))
		default:
			fmt.Fprint(w, "method is not supported")
		}
	})

	proxyDialer, err := socks5.NewDialer("socks5://" + os.Getenv("PROXY"))
	if err != nil {
		panic(err)
	}

	dialer := quic_over_proxy.NewDialer(utils.GenerateTLSConfig(), nil, proxyDialer, false)

	listenConfig := quic_over_proxy.DialerAsListen(dialer)

	listener, err := listenConfig.Listen(context.Background(), "udp", os.Getenv("TARGET"))
	if err != nil {
		panic(err)
	}

	logger.Println("Client as Server is starting...")
	server := http.Server{
		Handler: mux,
	}
	logger.Println(server.Serve(listener))
}
