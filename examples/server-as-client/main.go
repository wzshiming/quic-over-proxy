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
)

func main() {
	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	listenConfig := quic_over_proxy.NewListenConfig(utils.GenerateTLSConfig(), nil, true)

	dialer := quic_over_proxy.ListenAsDialer(listenConfig)
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
			"http://"+os.Getenv("ADDRESS"),
			bytes.NewReader([]byte(`hello, client as server!`)),
		)
		if err != nil {
			logger.Println(err)
			continue
		}

		resp, err := client.Do(req)
		if err != nil {
			logger.Println(err)
			continue
		}

		result, err := io.ReadAll(resp.Body)
		if err != nil {
			logger.Println(err)
			continue
		}

		logger.Println(resp.StatusCode, string(result))
	}
}
