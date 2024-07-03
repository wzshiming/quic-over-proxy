package main

import (
	"fmt"
	"log"
	"os"

	"github.com/wzshiming/socks5"
)

func main() {
	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
	logger.Println("Proxy is listening...")
	fmt.Println(socks5.NewServer().ListenAndServe("tcp", os.Getenv("ADDRESS")))
}
