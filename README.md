# quic-over-proxy

Let the QUIC protocol using proxy server and reverse penetration.

Refer to and modified from [WankkoRee/go-http3-proxy](https://github.com/WankkoRee/go-http3-proxy)

### Test

```bash 
## Terminal_1
$ ADDRESS=:1080 go run ./examples/proxy
> Proxy is listening...

## Terminal_2
$ ADDRESS=:8080 go run ./examples/server
> Server is listening...

## Terminal_3
## This terminal depends on Terminal_1 & Terminal_2
$ TARGET=127.0.0.1:8080 PROXY=127.0.0.1:1080 go run ./examples/client
> 200 echo:'hello, server!'
> 200 echo:'hello, server!'
> 200 echo:'hello, server!'
...
```

### Reverse

```bash 
## Terminal_1
$ ADDRESS=:1080 go run ./examples/proxy
> Proxy is listening...

## Terminal_2
$ TARGET=127.0.0.1:8080 PROXY=127.0.0.1:1080 go run ./examples/client-as-server
> Client as Server is starting...

## Terminal_3
## This terminal depends on Terminal_1 & Terminal_2
$ ADDRESS=:8080 go run ./examples/server-as-client
> 200 echo:'hello, client as server!'
> 200 echo:'hello, client as server!'
> 200 echo:'hello, client as server!'
...
```

## Docker

This example can also be run using docker.

### Running

```bash 
$ make -C examples down build up
proxy-1             | Proxy is listening...
server-1            | Server is listening...
client-as-server-1  | Client as Server is starting...
client-1            | 200 echo:'hello, server!'
server-as-client-1  | 200 echo:'hello, client as server!'
client-1            | 200 echo:'hello, server!'
server-as-client-1  | 200 echo:'hello, client as server!'
client-1            | 200 echo:'hello, server!'
server-as-client-1  | 200 echo:'hello, client as server!'
...
```
