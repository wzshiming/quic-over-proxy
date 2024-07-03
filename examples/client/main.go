package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	quic_over_proxy "github.com/wzshiming/quic-over-proxy"
	"github.com/wzshiming/quic-over-proxy/examples/utils"
	"github.com/wzshiming/socks5"
)

func main() {
	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	var proxyDialer quic_over_proxy.Dialer
	if p := os.Getenv("PROXY"); p != "" {
		dialer, err := socks5.NewDialer("socks5://" + p)
		if err != nil {
			panic(err)
		}
		proxyDialer = dialer
	}

	dialer := quic_over_proxy.NewDialer(utils.GenerateTLSConfig(), nil, proxyDialer, true)
	client := http.Client{
		Transport: &http.Transport{
			DialContext: dialer.DialContext,
		},
		Timeout: 10 * time.Second,
	}

	for {
		time.Sleep(time.Second)
		req, err := http.NewRequest(
			http.MethodPost,
			"http://"+os.Getenv("TARGET"),
			bytes.NewReader([]byte(`hello, server!`)),
		)
		if err != nil {
			panic(err)
		}

		resp, err := client.Do(req)
		if err != nil {
			logger.Println(err)
			continue
		}

		result, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		logger.Println(resp.StatusCode, string(result))
	}
}
