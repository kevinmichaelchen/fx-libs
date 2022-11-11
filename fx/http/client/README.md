An [FX module](https://github.com/uber-go/fx) that provides an HTTP client.

## Installing
```shell
go get -v github.com/kevinmichaelchen/fx-libs/fx/http/client
```

## Environment variables
| Env Var                                          | Purpose                                                                                                                                                   | Default |
|--------------------------------------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------|---------|
| `HTTP_CLIENT_TRANSPORT_DIAL_TIMEOUT`             | max amount of time a dial will wait for a connect to complete                                                                                             | 10s     |
| `HTTP_CLIENT_TRANSPORT_DIAL_KEEP_ALIVE`          | the interval between keep-alive probes for an active network connection.                                                                                  | 30s     |
| `HTTP_CLIENT_TRANSPORT_IDLE_CONN_TIMEOUT`        | max amount of time an idle (keep-alive) connection will remain idle before closing itself                                                                 | 90s     |
| `HTTP_CLIENT_TRANSPORT_MAX_IDLE_CONNS`           | max number of idle (keep-alive) connections across all hosts                                                                                              | 100     |
| `HTTP_CLIENT_TRANSPORT_MAX_IDLE_CONNS_PER_HOST`  | max idle (keep-alive) connections to keep per-host                                                                                                        | 10      |
| `HTTP_CLIENT_TRANSPORT_EXPECT_CONTINUE_TIMEOUT`  | amount of time to wait for a server's first response headers after fully writing the request headers if the request has an "Expect: 100-continue" header  | 1s      |
| `HTTP_CLIENT_TRANSPORT_RESPONSE_HEADER_TIMEOUT`  | max amount of time to wait for a server's response headers after fully writing the request (including its body, if any)                                   | 15s     |
| `HTTP_CLIENT_TRANSPORT_TLS_HANDSHAKE_TIMEOUT`    | max amount of time to wait for TLS handshake                                                                                                              | 5s      |
