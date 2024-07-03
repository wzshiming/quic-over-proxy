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

	listenConfig := quic_over_proxy.NewListenConfig(utils.GenerateTLSConfig(), nil, false)
	listener, err := listenConfig.Listen(context.Background(), "udp", os.Getenv("ADDRESS"))
	if err != nil {
		panic(err)
	}

	logger.Println("Server is listening...")
	server := http.Server{
		Handler: mux,
	}
	logger.Println(server.Serve(listener))
}
